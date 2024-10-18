package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/calvindc/dpc-tox/librarywrapper/libtox"
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
	var options *libtox.Options

	flag.StringVar(&filepath, "save", "./testdata/test_savedata", "path to save file")
	flag.Parse()

	fmt.Printf("[INFO] Using Tox Library version %d.%d.%d\n", libtox.VersionMajor(), libtox.VersionMinor(), libtox.VersionPatch())

	if !libtox.VersionIsCompatible(0, 2, 19) {
		fmt.Println("[ERROR] The compiled library (toxcore) is not compatible with this example.")
		fmt.Println("[ERROR] Please update your Tox library. If this error persists, please report it to the dpc_tox developers.")
		fmt.Println("[ERROR] Thanks!")
		return
	}

	savedata, err := loadData(filepath)
	if err == nil {
		fmt.Println("[INFO] Loading Tox profile from savedata...")
		options = &libtox.Options{
			IPv6Enabled:  true,
			UDPEnabled:   true,
			ProxyType:    libtox.TOX_PROXY_TYPE_NONE,
			ProxyHost:    "127.0.0.1",
			ProxyPort:    5555,
			StartPort:    0,
			EndPort:      0,
			TcpPort:      0, // only enable TCP server if your client provides an option to disable it
			SaveDataType: libtox.TOX_SAVEDATA_TYPE_TOX_SAVE,
			SaveData:     savedata}
	} else {
		fmt.Println("[INFO] Creating new Tox profile...")
		options = nil // default options
		newToxInstance = true
	}

	tox, err := libtox.New(options)
	if err != nil {
		panic(err)
	}

	if newToxInstance {
		err = tox.SelfSetName("tox-lib-debuger")
		if err != nil {
			fmt.Println(err)
		}
		err = tox.SelfSetStatusMessage("God has forgotten me!")
		if err != nil {
			fmt.Println(err)
		}
	}

	addr, _ := tox.SelfGetAddress()
	fmt.Println("TOX ID:\t\t", strings.ToUpper(hex.EncodeToString(addr)))
	public, _ := tox.SelfGetPublicKey()
	fmt.Println("TOX Public Key:\t", strings.ToUpper(hex.EncodeToString(public)))
	secert, _ := tox.SelfGetSecretKey()
	fmt.Println("TOX Secret Key:\t", strings.ToUpper(hex.EncodeToString(secert)))

	err = tox.SelfSetStatus(libtox.TOX_USERSTATUS_NONE)
	//err = tox.SelfSetStatus(dpc_tox.TOX_USERSTATUS_BUSY)

	// Register our callbacks
	tox.CallbackFriendRequest(onFriendRequest)
	tox.CallbackFriendMessage(onFriendMessage)
	tox.CallbackFileRecvControl(onFileRecvControl)
	tox.CallbackFileChunkRequest(onFileChunkRequest)
	tox.CallbackFileRecv(onFileRecv)
	tox.CallbackFileRecvChunk(onFileRecvChunk)

	tox.CallbackConferenceInvite(onConferenceInvite)
	tox.CallbackConferenceConnected(onConferenceConnected)
	tox.CallbackConferenceMessage(onConferenceMessage)

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

func onFriendRequest(t *libtox.Tox, publicKey []byte, message []byte, length uint32) {
	fmt.Printf("New friend request from %s\n", hex.EncodeToString(publicKey))
	fmt.Printf("With message: %v\n", string(message))
	// Auto-accept friend request
	t.FriendAddNorequest(publicKey)
}

func onFriendMessage(t *libtox.Tox, friendNumber uint32, messagetype libtox.ToxMessageType, message []byte, length uint32) {
	if messagetype == libtox.TOX_MESSAGE_TYPE_NORMAL {
		friendName, _ := t.FriendGetName(friendNumber)
		pubKey, _ := t.FriendGetPublickey(friendNumber)
		fmt.Printf("New message from friend number[%d], name[%s], ID[%s], message[%s]\n", friendNumber, friendName, hex.EncodeToString(pubKey), message)
	} else {
		fmt.Printf("New action from %d : %s\n", friendNumber, message)
	}

	switch string(message) {
	case "00":
		sendMes := []byte("Type \n00(for help)\n11(send you a pic)\n22(i do create a new group and invite you to my all group)\n33(i get all group peers info) \n to use some internal function for test.")
		t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, sendMes)

	case "11":
		file, err := os.Open("./testdata/response.png")
		if err != nil {
			defer file.Close()
			t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("File not found. Please 'cd' into tox project"))
			//file.Close()
			return
		}

		// get the file size
		stat, err := file.Stat()
		if err != nil {
			t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("Could not read file stats."))
			file.Close()
			return
		}

		fmt.Println("File size is ", stat.Size())

		fileNumber, err := t.FileSend(friendNumber, libtox.TOX_FILE_KIND_DATA, uint64(stat.Size()), nil, "fileName.png")
		if err != nil {
			t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("t.FileSend() failed."))
			file.Close()
			return
		}

		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: uint64(stat.Size())}

	case "22":
		conferenceNumber, err := t.ConferenceNew()
		if err != nil {
			t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("create group err:"+err.Error()))
		}
		t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte(fmt.Sprintf("i will create a new group named 'group-nano' with you, conferenceNumber=%d", conferenceNumber)))
		title := fmt.Sprintf("group-%d", time.Now().Nanosecond()/1e+6)
		success, err := t.ConferenceSetTitle(conferenceNumber, title)
		if err != nil {
			fmt.Println(fmt.Sprintf("ConferenceSetTitle for conferenceNumber=[%d],err=%v", conferenceNumber, err))
		}
		if success {
			t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("the new group build success, title is:"+title))
		}
		//test for invite to my room
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
	case "33": //check all peers info
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
			t.ConferenceSendMessage(theGp, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("welcome group...^~^"))
		}

	default:
		t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("i see, your message is:"+string(message)+", right?"))
	}
}

func onFileRecv(t *libtox.Tox, friendNumber uint32, fileNumber uint32, kind libtox.ToxFileKind, filesize uint64, filename string, length uint32) {
	fmt.Println("callback onFileRecv")
	if kind == libtox.TOX_FILE_KIND_AVATAR {

		if filesize > MAX_AVATAR_SIZE {
			// reject file send request
			t.FileControl(friendNumber, fileNumber, libtox.TOX_FILE_CONTROL_CANCEL)
			return
		}

		publicKey, _ := t.FriendGetPublickey(friendNumber)
		file, err := os.Create("./testdata/file_recv_" + hex.EncodeToString(publicKey) + ".png")
		if err != nil {
			fmt.Println("[ERROR] Error creating file", "test_"+hex.EncodeToString(publicKey)+".png")
		}

		// append the file to the map of active file transfers
		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: filesize}

		// accept the file send request
		t.FileControl(friendNumber, fileNumber, libtox.TOX_FILE_CONTROL_RESUME)

	} else {
		// accept files of any length

		file, err := os.Create("./testdata/file_recv_" + filename)
		if err != nil {
			fmt.Println("[ERROR] Error creating file", "test_"+filename)
		}

		// append the file to the map of active file transfers
		transfers[fileNumber] = FileTransfer{fileHandle: file, fileSize: filesize}

		// accept the file send request
		t.FileControl(friendNumber, fileNumber, libtox.TOX_FILE_CONTROL_RESUME)
	}
}

func onFileRecvControl(t *libtox.Tox, friendNumber uint32, fileNumber uint32, fileControl libtox.ToxFileControl) {
	fmt.Println("callback onFileRecvControl")
	transfer, ok := transfers[fileNumber]
	if !ok {
		fmt.Println("Error: File handle does not exist")
		return
	}

	if fileControl == libtox.TOX_FILE_CONTROL_CANCEL {
		// delete file handle
		transfer.fileHandle.Sync()
		transfer.fileHandle.Close()
		delete(transfers, fileNumber)
	}
}

func onFileChunkRequest(t *libtox.Tox, friendNumber uint32, fileNumber uint32, position uint64, length uint64) {
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

func onFileRecvChunk(t *libtox.Tox, friendNumber uint32, fileNumber uint32, position uint64, data []byte, length uint32) {
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
		t.FriendSendMessage(friendNumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("Thanks!"))
	}
}

func onConferenceInvite(t *libtox.Tox, friendnumber uint32, conferencetype libtox.ToxConferenceType, cookie []byte) {
	fmt.Println("callback onConferenceInvite")
	fmt.Printf("New conference invite from [%d], conferenceType=%v, ", friendnumber, conferencetype)
	fmt.Printf("With cookie: %s\n", cookie)

	switch conferencetype {
	case libtox.TOX_CONFERENCE_TYPE_TEXT:
		ret, err := t.ConferenceJoin(friendnumber, cookie)
		if err != nil {
			fmt.Println(fmt.Sprintf("ConferenceJoin err=%v, ret=%v", err, ret))
		}

	case libtox.TOX_CONFERENCE_TYPE_AV:
		t.FriendSendMessage(friendnumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("can not join av group"))
	default:
		t.FriendSendMessage(friendnumber, libtox.TOX_MESSAGE_TYPE_NORMAL, []byte("unknow conferencetype"))
	}
}
func onConferenceMessage(t *libtox.Tox, conferencenumber uint32, peernumber uint32, messagetype libtox.ToxMessageType, message []byte, length uint32) {
	/*if peernumber == 0 {
		return //self
	}*/
	fmt.Println("callback onConferenceInvite")
	fmt.Println("*******************************")
	fmt.Println(fmt.Sprintf("onConferenceMessage:conferencenumber=\t[%d]", conferencenumber))
	fmt.Println(fmt.Sprintf("onConferenceMessage:peernumber=\t\t[%d]", peernumber))
	fmt.Println(fmt.Sprintf("onConferenceMessage:messagetype=\t[%v]", messagetype))
	fmt.Println(fmt.Sprintf("onConferenceMessage:message=\t\t[%s]", message))
	fmt.Println("*******************************")

	//t.ConferenceSendMessage(conferencenumber, dpc_tox.TOX_MESSAGE_TYPE_NORMAL, "This is an automatic reply ["+message+"].")
}
func onConferenceConnected(t *libtox.Tox, conferencenumber uint32) {
	fmt.Println("callback onConferenceConnected")
	fmt.Println(fmt.Sprintf("i have joined the group[%d] and in", conferencenumber))
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
func saveData(t *libtox.Tox, filepath string) error {
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
