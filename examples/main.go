package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/calvindc/dpc-tox"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Server struct {
	Address   string
	Port      uint16
	PublicKey []byte
}

const MAX_AVATAR_SIZE = 65536 // see github.com/Tox/Tox-STS/blob/master/STS.md#avatars

type FileTransfer struct {
	fileHandle *os.File
	fileSize   uint64
}

// Map of active file transfers
var transfers = make(map[uint32]FileTransfer)

func main() {
	var newToxInstance bool = false
	var filepath string
	var options *dpc_tox.Options

	flag.StringVar(&filepath, "save", "./example_savedata", "path to save file")
	flag.Parse()

	fmt.Printf("[INFO] Using Tox version %d.%d.%d\n", dpc_tox.VersionMajor(), dpc_tox.VersionMinor(), dpc_tox.VersionPatch())

	if !dpc_tox.VersionIsCompatible(0, 2, 19) {
		fmt.Println("[ERROR] The compiled library (toxcore) is not compatible with this example.")
		fmt.Println("[ERROR] Please update your Tox library. If this error persists, please report it to the dpc_tox developers.")
		fmt.Println("[ERROR] Thanks!")
		return
	}

	savedata, err := loadData(filepath)
	if err == nil {
		fmt.Println("[INFO] Loading Tox profile from savedata...")
		options = &dpc_tox.Options{
			IPv6Enabled:  true,
			UDPEnabled:   true,
			ProxyType:    dpc_tox.TOX_PROXY_TYPE_NONE,
			ProxyHost:    "127.0.0.1",
			ProxyPort:    5555,
			StartPort:    0,
			EndPort:      0,
			TcpPort:      0, // only enable TCP server if your client provides an option to disable it
			SaveDataType: dpc_tox.TOX_SAVEDATA_TYPE_TOX_SAVE,
			SaveData:     savedata}
	} else {
		fmt.Println("[INFO] Creating new Tox profile...")
		options = nil // default options
		newToxInstance = true
	}

	tox, err := dpc_tox.New(options)
	if err != nil {
		panic(err)
	}

	if newToxInstance {
		tox.SelfSetName("tox-debuger")
		tox.SelfSetStatusMessage("I am debuging tox!")
	}

	addr, _ := tox.SelfGetAddress()
	fmt.Println("TOX ID: ", strings.ToUpper(hex.EncodeToString(addr)))
	secert, _ := tox.SelfGetSecretKey()
	fmt.Println("TOX Secret Key: ", strings.ToUpper(hex.EncodeToString(secert)))

	err = tox.SelfSetStatus(dpc_tox.TOX_USERSTATUS_NONE)
	//err = tox.SelfSetStatus(dpc_tox.TOX_USERSTATUS_BUSY)

	// Register our callbacks
	tox.CallbackFriendRequest(onFriendRequest)
	tox.CallbackFriendMessage(onFriendMessage)
	tox.CallbackFileRecvControl(onFileRecvControl)
	tox.CallbackFileChunkRequest(onFileChunkRequest)
	tox.CallbackFileRecv(onFileRecv)
	tox.CallbackFileRecvChunk(onFileRecvChunk)

	tox.CallbackConferenceInvite(onConferenceInvite)

	/* Connect to the network
	 * Use more than one node in a real world szenario. This example relies one
	 * the following node to be up.
	 */
	pubkey, _ := hex.DecodeString("E20ABCF38CDBFFD7D04B29C956B33F7B27A3BB7AF0618101617B036E4AEA402D")
	server := &Server{"3.0.24.15", 33445, pubkey}

	// tox boot
	err = tox.Bootstrap(server.Address, server.Port, server.PublicKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("[INFO] Tox bootstrap sucessfully.")
	isRunning := true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ticker := time.NewTicker(25 * time.Millisecond)

	for isRunning {
		select {
		case <-c:
			fmt.Printf("\nSaving data...\n")
			if err := saveData(tox, filepath); err != nil {
				fmt.Println("[ERROR]", err)
			}
			fmt.Println("tox killing")
			isRunning = false
			tox.Kill()
		case <-ticker.C:
			tox.Iterate()
		}
	}
}

func onFriendRequest(t *dpc_tox.Tox, publicKey []byte, message string) {
	fmt.Printf("New friend request from %s\n", hex.EncodeToString(publicKey))
	fmt.Printf("With message: %v\n", message)
	// Auto-accept friend request
	t.FriendAddNorequest(publicKey)
}

func onFriendMessage(t *dpc_tox.Tox, friendNumber uint32, messagetype dpc_tox.ToxMessageType, message string) {
	if messagetype == dpc_tox.TOX_MESSAGE_TYPE_NORMAL {
		friendName, _ := t.FriendGetName(friendNumber)
		pubKey, _ := t.FriendGetPublickey(friendNumber)
		fmt.Printf("New message from friend number[%d], name[%s], ID[%s], message[%s]\n", friendNumber, friendName, hex.EncodeToString(pubKey), message)
	} else {
		fmt.Printf("New action from %d : %s\n", friendNumber, message)
	}

	switch message {
	case "00":
		t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "Type '/help case' to use some internal function.")

	case "11":
		file, err := os.Open("/home/cy/godev/src/github.com/calvindc/dpc-tox/examples/response.png")
		if err != nil {
			defer file.Close()
			t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "File not found. Please 'cd' into tox project")
			//file.Close()
			return
		}

		// get the file size
		stat, err := file.Stat()
		if err != nil {
			t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "Could not read file stats.")
			file.Close()
			return
		}

		fmt.Println("File size is ", stat.Size())

		fileNumber, err := t.FileSend(friendNumber, dpc_tox.TOX_FILE_KIND_DATA, uint64(stat.Size()), nil, "fileName.png")
		if err != nil {
			t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "t.FileSend() failed.")
			file.Close()
			return
		}

		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: uint64(stat.Size())}

	case "22":
		conferenceNumber, err := t.ConferenceNew()
		if err != nil {
			t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "create group err:"+err.Error())
		}
		t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, fmt.Sprintf("i will create a new group named 'GROUP-X' with you, conferenceNumber=%d", conferenceNumber))
		title := fmt.Sprintf("GROUP%d", time.Now().UnixNano())
		success, err := t.ConferenceSetTitle(conferenceNumber, title)
		if err != nil {
			fmt.Println(fmt.Sprintf("ConferenceSetTitle for conferenceNumber=[%d],err=%v", conferenceNumber, err))
		}
		if success {
			t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "the new group build success, title is:"+title)
		}
		//邀请到我建立的所有群里
		myAllGroup, err := t.ConferenceGetChatlist()
		if err != nil {
			fmt.Println(fmt.Sprintf("ConferenceGetChatlist err=%v", err))
		}
		for _, theGp := range myAllGroup {
			ret, err := t.ConferenceInvite(friendNumber, theGp)
			if err != nil {
				fmt.Println(fmt.Sprintf("ConferenceInvite invite [%d] into [%d] failed,err=%t,ret=%t", friendNumber, theGp, err, ret))
			}
			fmt.Println(fmt.Sprintf("ConferenceInvite invite [%d] into [%d] success,ret=%v", friendNumber, theGp, ret))
		}
	case "33": //查看所有群的成员情况
		myAllGroup, err := t.ConferenceGetChatlist()
		if err != nil {
			fmt.Println(fmt.Sprintf("ConferenceGetChatlist err=%v", err))
		}
		for _, theGp := range myAllGroup {
			var groupPeersInfo = make(map[uint32]string)
			groupPeersInfo, err := t.ConferenceGetPeers(theGp)
			if err != nil {
				fmt.Println(fmt.Sprintf("ConferenceGetPeers failed, groupdNumber=%t, err=%v", theGp, err))
			}
			for peerNumber, pubKey := range groupPeersInfo {
				fmt.Println(fmt.Sprintf("ConferenceGetPeers groupNubmber=%d, peerNumber=%d,peerPubkey=%s", theGp, peerNumber, pubKey))
			}
			t.ConferenceSendMessage(theGp, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "欢迎热烈讨论...^~^")
		}

	default:
		t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "i see, your message is:"+message+", right?")
	}
}

func onFileRecv(t *dpc_tox.Tox, friendNumber uint32, fileNumber uint32, kind dpc_tox.ToxFileKind, filesize uint64, filename string) {
	fmt.Println("callback onFileRecv")
	if kind == dpc_tox.TOX_FILE_KIND_AVATAR {

		if filesize > MAX_AVATAR_SIZE {
			// reject file send request
			t.FileControl(friendNumber, fileNumber, dpc_tox.TOX_FILE_CONTROL_CANCEL)
			return
		}

		publicKey, _ := t.FriendGetPublickey(friendNumber)
		file, err := os.Create("example_" + hex.EncodeToString(publicKey) + ".png")
		if err != nil {
			fmt.Println("[ERROR] Error creating file", "example_"+hex.EncodeToString(publicKey)+".png")
		}

		// append the file to the map of active file transfers
		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: filesize}

		// accept the file send request
		t.FileControl(friendNumber, fileNumber, dpc_tox.TOX_FILE_CONTROL_RESUME)

	} else {
		// accept files of any length

		file, err := os.Create("example_" + filename)
		if err != nil {
			fmt.Println("[ERROR] Error creating file", "example_"+filename)
		}

		// append the file to the map of active file transfers
		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: filesize}

		// accept the file send request
		t.FileControl(friendNumber, fileNumber, dpc_tox.TOX_FILE_CONTROL_RESUME)
	}
}

func onFileRecvControl(t *dpc_tox.Tox, friendNumber uint32, fileNumber uint32, fileControl dpc_tox.ToxFileControl) {
	fmt.Println("callback onFileRecvControl")
	transfer, ok := transfers[fileNumber]
	if !ok {
		fmt.Println("Error: File handle does not exist")
		return
	}

	if fileControl == dpc_tox.TOX_FILE_CONTROL_CANCEL {
		// delete file handle
		transfer.fileHandle.Sync()
		transfer.fileHandle.Close()
		delete(transfers, fileNumber)
	}
}

func onFileChunkRequest(t *dpc_tox.Tox, friendNumber uint32, fileNumber uint32, position uint64, length uint64) {
	fmt.Println(fmt.Sprintf("callback onFileChunkRequest friendNumber:%d,,position:%d,length:%d", friendNumber, position, length))

	transfer, ok := transfers[fileNumber]
	if !ok {
		fmt.Println("Error: File handle does not exist")
		return
	}

	// a zero-length chunk request confirms that the file was successfully transferred
	if length == 0 {
		transfer.fileHandle.Close()
		delete(transfers, fileNumber)
		fmt.Println("File transfer completed (sending)", fileNumber)
		return
	}

	// read the requested data to send
	data := make([]byte, length)
	_, err := transfers[fileNumber].fileHandle.ReadAt(data, int64(position))
	if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	// send the requested data
	t.FileSendChunk(friendNumber, fileNumber, position, data)
}

func onFileRecvChunk(t *dpc_tox.Tox, friendNumber uint32, fileNumber uint32, position uint64, data []byte) {
	fmt.Println("callback onFileRecvChunk")
	transfer, ok := transfers[fileNumber]
	if !ok {
		if len(data) == 0 {
			// ignore the zero-length chunk that indicates that the transfer is
			// complete (see below)
			return
		}

		fmt.Println("Error: File handle does not exist")
		return
	}

	// write the received data to the file handle
	transfer.fileHandle.WriteAt(data, (int64)(position))

	// file transfer completed
	if position+uint64(len(data)) >= transfer.fileSize {
		// Some clients will send us another zero-length chunk without data (only
		// required for streams, not necessary for files with a known size) and some
		// will not.
		// We will delete the file handle now (we aleady received the whole file)
		// and ignore the file handle error when the zero-length chunk arrives.

		transfer.fileHandle.Sync()
		transfer.fileHandle.Close()
		delete(transfers, fileNumber)
		fmt.Println("File transfer completed (receiving)", fileNumber)
		t.FriendSendMessage(friendNumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "Thanks!")
	}
}

func onConferenceInvite(t *dpc_tox.Tox, friendnumber uint32, conferencetype dpc_tox.ToxConferenceType, cookie string) {
	fmt.Println("callback onConferenceInvite")
	fmt.Printf("New conference invite from [%d], conferenceType=%v, ", friendnumber, conferencetype)
	fmt.Printf("With cookie: %s\n", cookie)

	switch conferencetype {
	case dpc_tox.TOX_CONFERENCE_TYPE_TEXT:
		ret, err := t.ConferenceJoin(friendnumber, cookie)
		if err != nil {
			fmt.Println(fmt.Sprintf("ConferenceJoin err=%v, ret=%v", err, ret))
		}

	case dpc_tox.TOX_CONFERENCE_TYPE_AV:
		t.FriendSendMessage(friendnumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "can not join av group")
	default:
		t.FriendSendMessage(friendnumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "unknow conferencetype")
	}
}

// loadData reads a file and returns its content as a byte array
func loadData(filepath string) ([]byte, error) {
	if len(filepath) == 0 {
		return nil, errors.New("Empty path")
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return data, err
}

// saveData writes the savedata from toxcore to a file
func saveData(t *dpc_tox.Tox, filepath string) error {
	if len(filepath) == 0 {
		return errors.New("Empty path")
	}

	data, err := t.GetSavedata()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath, data, 0644)
	return err
}
