package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/calvindc/dpc-tox/librarywrapper/libtox"
	"log"
	"os"
	"time"
)

func onFriendRequest(t *libtox.Tox, publicKey []byte, message []byte, length uint32) {
	log.Printf("New friend request from %s\n", hex.EncodeToString(publicKey))

	storage.StoreFriendRequest(hex.EncodeToString(publicKey), string(message))
	broadcastToClients(createSimpleJSONEvent("friend_requests_update"))
}

func onFriendMessage(t *libtox.Tox, friendnumber uint32, messagetype libtox.ToxMessageType, message []byte, length uint32) {
	type jsonEvent struct {
		Type     string `json:"type"`
		Friend   uint32 `json:"friend"`
		Time     int64  `json:"time"`
		Message  string `json:"message"`
		IsAction bool   `json:"isAction"`
	}

	e, _ := json.Marshal(jsonEvent{
		Type:     "friend_message",
		Friend:   friendnumber,
		Time:     time.Now().Unix() * 1000,
		Message:  string(message),
		IsAction: messagetype == libtox.TOX_MESSAGE_TYPE_ACTION,
	})

	publicKey, _ := tox.FriendGetPublickey(friendnumber)
	storage.StoreMessage(hex.EncodeToString(publicKey), true, messagetype == libtox.TOX_MESSAGE_TYPE_ACTION, string(message))

	broadcastToClients(string(e))
}

func onFriendConnectionStatusChanges(t *libtox.Tox, friendnumber uint32, connectionStatus libtox.ToxConnection) {
	type jsonEvent struct {
		Type   string `json:"type"`
		Friend uint32 `json:"friend"`
		Online bool   `json:"online"`
	}

	e, _ := json.Marshal(jsonEvent{
		Type:   "connection_status",
		Friend: friendnumber,
		Online: connectionStatus != libtox.TOX_CONNECTION_NONE,
	})

	broadcastToClients(string(e))
}

func onFriendNameChanges(t *libtox.Tox, friendnumber uint32, newname []byte, length uint32) {
	type jsonEvent struct {
		Type   string `json:"type"`
		Friend uint32 `json:"friend"`
		Name   string `json:"name"`
	}

	e, _ := json.Marshal(jsonEvent{
		Type:   "name_changed",
		Friend: friendnumber,
		Name:   string(newname),
	})

	broadcastToClients(string(e))
}

func onFriendStatusMessageChanges(t *libtox.Tox, friendnumber uint32, status []byte, length uint32) {
	type jsonEvent struct {
		Type      string `json:"type"`
		Friend    uint32 `json:"friend"`
		StatusMsg string `json:"status_msg"`
	}

	e, _ := json.Marshal(jsonEvent{
		Type:      "status_message_changed",
		Friend:    friendnumber,
		StatusMsg: string(status),
	})

	broadcastToClients(string(e))
}

func onFriendStatusChanges(t *libtox.Tox, friendnumber uint32, userstatus libtox.ToxUserStatus) {
	type jsonEvent struct {
		Type   string `json:"type"`
		Friend uint32 `json:"friend"`
		Status string `json:"status"`
	}

	e, _ := json.Marshal(jsonEvent{
		Type:   "status_changed",
		Friend: friendnumber,
		Status: getUserStatusAsString(userstatus),
	})

	broadcastToClients(string(e))
}

func onFileRecv(t *libtox.Tox, friendnumber uint32, filenumber uint32, kind libtox.ToxFileKind, filesize uint64, filename string, length uint32) {
	if kind == libtox.TOX_FILE_KIND_AVATAR {
		publicKey, _ := tox.FriendGetPublickey(friendnumber)
		file, err := os.Create("../html/avatars/" + hex.EncodeToString(publicKey) + ".png")
		if err != nil {
			log.Println("[ERROR] Error creating file", "../html/avatars/"+hex.EncodeToString(publicKey)+".png")
		}

		// only accept avatars with a file size <= CFG_MAX_AVATAR_SIZE
		if filesize <= CFG_MAX_AVATAR_SIZE {
			// append the file to the map of active file transfers
			transfers[filenumber] = FileTransfer{fileHandle: file, fileSize: filesize, fileKind: kind}

			t.FileControl(friendnumber, filenumber, libtox.TOX_FILE_CONTROL_RESUME)
		} else {
			t.FileControl(friendnumber, filenumber, libtox.TOX_FILE_CONTROL_CANCEL)
		}

	} else if kind == libtox.TOX_FILE_KIND_DATA {
		file, err := os.Create("../html/download/" + string(filename))
		if err != nil {
			log.Println("[ERROR] Error creating file", "../html/download/"+string(filename))
		}

		// append the file to the map of active file transfers
		transfers[filenumber] = FileTransfer{fileHandle: file, fileSize: filesize, fileKind: kind}

		// TODO do not accept any file send request without asking the user
		t.FileControl(friendnumber, filenumber, libtox.TOX_FILE_CONTROL_RESUME)

	} else {
		log.Print("onFileRecv: unknown TOX_FILE_KIND: ", kind)
	}
}

func onFileRecvControl(t *libtox.Tox, friendnumber uint32, filenumber uint32, fileControl libtox.ToxFileControl) {
	transfer, ok := transfers[filenumber]
	if !ok {
		log.Println("Error: File handle does not exist")
		return
	}

	// TODO handle TOX_FILE_CONTROL_RESUME and TOX_FILE_CONTROL_PAUSE
	if fileControl == libtox.TOX_FILE_CONTROL_CANCEL {
		// delete file handle
		transfer.fileHandle.Close()
		delete(transfers, filenumber)
	}
}

func onFileRecvChunk(t *libtox.Tox, friendnumber uint32, filenumber uint32, position uint64, data []byte, length uint32) {
	transfer, ok := transfers[filenumber]
	if !ok {
		if len(data) == 0 {
			// ignore the zero-length chunk that indicates that the transfer is
			// complete (see below)
			return
		}

		log.Println("Error: File handle does not exist")
		return
	}

	// write data to the file handle
	transfer.fileHandle.WriteAt(data, (int64)(position))

	// file transfer completed
	if position+uint64(len(data)) >= transfer.fileSize {
		// Some clients will send us another zero-length chunk without data (only
		// required for stream, not necessary for files with a known size) and some
		// will not.
		// We will delete the file handle now (we aleady reveived the whole file)
		// and ignore the file handle error when the empty chunk arrives.

		fileKind := transfer.fileKind

		transfer.fileHandle.Sync()
		transfer.fileHandle.Close()
		delete(transfers, filenumber)
		log.Println("File transfer completed (receiving)", filenumber)

		if fileKind == libtox.TOX_FILE_KIND_AVATAR {
			// update friendlist
			broadcastToClients(createSimpleJSONEvent("avatar_update"))
		}
	}
}
