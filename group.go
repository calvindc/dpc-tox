package dpc_tox

//#include <tox/tox.h>
//#include <stdlib.h>
import "C"
import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
)

//tox_add_groupchat

// ConferenceNew creates and connects to a new text conference.
func (t *Tox) ConferenceNew() (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var toxErrConferenceNew C.Tox_Err_Conference_New
	conferenceNumber := C.tox_conference_new(t.tox, &toxErrConferenceNew)

	switch ToxErrConferenceNew(toxErrConferenceNew) {
	case TOX_ERR_CONFERENCE_NEW_OK:
		return uint32(conferenceNumber), nil
	case TOX_ERR_CONFERENCE_NEW_INIT:
		return uint32(conferenceNumber), ErrConferenceNewFailedInitialize
	default:
		return uint32(conferenceNumber), ErrFuncFail
	}

	return uint32(conferenceNumber), ErrUnknown
}

// ConferenceDelete this function deletes a conference.
func (t *Tox) ConferenceDelete(conferenceNumber uint32) (bool, error) {
	if t.tox == nil {
		return false, ErrToxInit
	}

	var toxErrConferenceDelete C.Tox_Err_Conference_Delete
	ret := C.tox_conference_delete(t.tox, (C.uint32_t)(conferenceNumber), &toxErrConferenceDelete)
	if !bool(ret) {
		return bool(ret), ErrConferenceDeleteFailed
	}
	switch ToxErrConferenceDelete(toxErrConferenceDelete) {
	case TOX_ERR_CONFERENCE_DELETE_OK:
		return true, nil
	case TOX_ERR_CONFERENCE_DELETE_CONFERENCE_NOT_FOUND:
		return false, ErrConferenceDeleteConferenceNotFound
	default:
		return false, ErrFuncFail
	}

	return false, ErrUnknown
}

// ConferencePeerGetName
func (t *Tox) ConferencePeerGetName(conferenceNumber, peerNumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrToxInit
	}
	length, err := t.ConferencePeerGetNameSize(conferenceNumber, peerNumber)
	if err != nil {
		return "", ErrFuncFail
	}
	name := make([]byte, length)
	if length > 0 {
		var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query = C.TOX_ERR_CONFERENCE_PEER_QUERY_OK
		ret := C.tox_conference_peer_get_name(t.tox, (C.uint32_t)(conferenceNumber), (C.uint32_t)(peerNumber), (*C.uint8_t)(&name[0]), &toxErrConferencePeerQuery)
		if ret != true || ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(name), nil
}

func (t *Tox) ConferencePeerGetNameSize(conferenceNumber, peerNumber uint32) (int64, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}
	var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query = C.TOX_ERR_CONFERENCE_PEER_QUERY_OK
	ret := C.tox_conference_peer_get_name_size(t.tox, (C.uint32_t)(conferenceNumber), (C.uint32_t)(peerNumber), &toxErrConferencePeerQuery)
	if ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
		return 0, ErrFuncFail
	}
	return int64(ret), nil
}

func (t *Tox) ConferencePeerGetPublicKey(conferenceNumber uint32, peerNumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrToxInit
	}
	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)
	var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query
	r := C.tox_conference_peer_get_public_key(t.tox, (C.uint32_t)(conferenceNumber), (C.uint32_t)(peerNumber), (*C.uint8_t)(&publickey[0]), &toxErrConferencePeerQuery)
	if bool(r) != true || ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
		return "", ErrFuncFail
	}

	pubkey := strings.ToUpper(hex.EncodeToString(publickey[:]))
	return pubkey, nil
}

func (t *Tox) ConferenceInvite(friendNumber uint32, conferenceNumber uint32) (int, error) {
	if t.tox == nil {
		return -2, ErrToxInit
	}
	// if give a friendNumber which not exists,the tox_invite_friend has a strange behaive: cause other tox_* call failed
	// and the call will return true, but only strange thing accurs so just precheck the friendNumber and then go
	friendExist, err := t.FriendExists(friendNumber)
	if err != nil || friendExist == false {
		return -1, errors.New(fmt.Sprintf("friend not exists: %d", friendNumber))
	}

	var toxErrConferenceInvite C.Tox_Err_Conference_Invite
	r := C.tox_conference_invite(t.tox, (C.uint32_t)(friendNumber), (C.uint32_t)(conferenceNumber), &toxErrConferenceInvite)
	if r == false {
		return 0, errors.New(fmt.Sprintf("conference invite failed: %d", toxErrConferenceInvite))
	}
	switch ToxErrConferenceInvite(toxErrConferenceInvite) {
	case TOX_ERR_CONFERENCE_INVITE_OK:
		return 1, nil
	case TOX_ERR_CONFERENCE_INVITE_CONFERENCE_NOT_FOUND:
		return 0, ErrConferenceInviteConferenceNotFound
	case TOX_ERR_CONFERENCE_INVITE_FAIL_SEND:
		return 0, ErrConferenceInviteFailSend
	case TOX_ERR_CONFERENCE_INVITE_NO_CONNECTION:
		return 0, ErrConferenceInviteNoConnection
	default:
		return 0, ErrFuncFail
	}

	return 0, ErrUnknown
}

/*func (t *Tox) FriendExists(friendNumber uint32) bool {
	r := C.tox_friend_exists(t.tox, (C.uint32_t)(friendNumber))
	return bool(r)
}*/

func (t *Tox) ConferenceJoin(friendNumber uint32, cookie string) (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	if cookie == "" || len(cookie) < 20 {
		return 0, errors.New("Invalid cookie:" + cookie)
	}

	//data, err := hex.DecodeString(cookie)
	data := []byte(cookie)
	fmt.Println(len(data))
	//fmt.Println(err)
	if data == nil { // || len(data) < 10
		return 0, errors.New("Invalid data: " + cookie)
	}

	var toxErrConferenceJoin C.Tox_Err_Conference_Join
	ret := C.tox_conference_join(t.tox, (C.uint32_t)(friendNumber), (*C.uint8_t)(&data[0]), (C.size_t)(len(data)), &toxErrConferenceJoin)
	if ret == C.UINT32_MAX {
		return uint32(ret), errors.New(fmt.Sprintf("join group chat failed: %d", toxErrConferenceJoin))
	}

	return uint32(ret), nil
}

func (t *Tox) ConferenceSendMessage(conferenceNumber uint32, messageType ToxMessageType, message string) (bool, error) {
	if t.tox == nil {
		return false, ErrToxInit
	}

	if len(message) == 0 {
		return false, ErrArgs
	}

	var cMessageType C.TOX_MESSAGE_TYPE
	if messageType == TOX_MESSAGE_TYPE_ACTION {
		cMessageType = C.TOX_MESSAGE_TYPE_ACTION
	} else {
		cMessageType = C.TOX_MESSAGE_TYPE_NORMAL
	}
	cMessage := (*C.uint8_t)(&[]byte(message)[0])

	var toxErrConferenceSendMessage C.Tox_Err_Conference_Send_Message
	ret := C.tox_conference_send_message(t.tox, (C.uint32_t)(conferenceNumber), cMessageType, cMessage, (C.size_t)(len(message)), &toxErrConferenceSendMessage)
	if ret == false {
		return false, errors.New(fmt.Sprintf("group send message failed: %d", toxErrConferenceSendMessage))
	}
	if ToxErrConferenceSendMessage(toxErrConferenceSendMessage) != TOX_ERR_CONFERENCE_SEND_MESSAGE_OK {
		return false, ErrFuncFail
	}

	return bool(ret), nil
}

func (t *Tox) ConferenceSetTitle(conferenceNumber uint32, title string) (bool, error) {
	if t.tox == nil {
		return false, ErrToxInit
	}
	var cTitle (*C.uint8_t)
	if len(title) == 0 {
		cTitle = nil
	} else {
		cTitle = (*C.uint8_t)(&[]byte(title)[0])
	}
	var toxErrConferenceTitle C.Tox_Err_Conference_Title
	success := C.tox_conference_set_title(t.tox, (C.uint32_t)(conferenceNumber), cTitle, (C.size_t)(len(title)), &toxErrConferenceTitle)
	if !bool(success) || ToxErrConferenceTitle(toxErrConferenceTitle) != TOX_ERR_CONFERENCE_TITLE_OK {
		return false, ErrFuncFail
	}
	return true, nil
}

func (t *Tox) ConferenceGetTitle(conferenceNumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrToxInit
	}
	length, err := t.ConferenceGetTitleSize(conferenceNumber)
	if err != nil {
		return "", ErrFuncFail
	}
	title := make([]byte, length)
	var toxErrConferenceTitle C.Tox_Err_Conference_Title
	success := C.tox_conference_get_title(t.tox, (C.uint32_t)(conferenceNumber), (*C.uint8_t)(&title[0]), &toxErrConferenceTitle)
	if !bool(success) || ToxErrConferenceTitle(toxErrConferenceTitle) != TOX_ERR_CONFERENCE_TITLE_OK {
		return "", ErrFuncFail
	}

	return string(title), nil
}

func (t *Tox) ConferenceGetTitleSize(conferenceNumber uint32) (int64, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var toxErrConferenceTitle C.Tox_Err_Conference_Title
	ret := C.tox_conference_get_title_size(t.tox, (C.uint32_t)(conferenceNumber), &toxErrConferenceTitle)
	if ToxErrConferenceTitle(toxErrConferenceTitle) != TOX_ERR_CONFERENCE_TITLE_OK {
		return 0, ErrFuncFail
	}
	return int64(ret), nil
}

func (t *Tox) ConferencePeerNumberIsOurs(conferenceNumber, peerNumber uint32) (bool, error) {
	if t.tox == nil {
		return false, ErrToxInit
	}
	var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query
	ret := C.tox_conference_peer_number_is_ours(t.tox, (C.uint32_t)(conferenceNumber), (C.uint32_t)(conferenceNumber), &toxErrConferencePeerQuery)
	if !ret || ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
		return false, ErrFuncFail
	}
	return bool(ret), nil
}

func (t *Tox) ConferencePeerCount(conferenceNumber uint32) (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query
	ret := C.tox_conference_peer_count(t.tox, (C.uint32_t)(conferenceNumber), &toxErrConferencePeerQuery)
	if ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
		return 0, ErrFuncFail
	}
	return uint32(ret), nil
}

// extra combined api
func (t *Tox) ConferenceGetNames(conferenceNumber uint32) ([]string, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	peerCount, err := t.ConferencePeerCount(conferenceNumber)
	if err != nil {
		return nil, ErrFuncFail
	}
	peerNames := make([]string, peerCount)
	if peerCount == 0 {
		return peerNames, nil
	}

	for idx := uint32(0); idx < math.MaxUint32; idx++ {
		pname, err := t.ConferencePeerGetName(conferenceNumber, idx)
		if err != nil {
			break
		}
		peerNames[idx] = pname
		if uint32(len(peerNames)) >= peerCount {
			break
		}
	}

	return peerNames, nil
}

func (t *Tox) ConferenceGetPeerPubkeys(conferenceNumber uint32) ([]string, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	peerPubkeys := make([]string, 0)
	peerCount, err := t.ConferencePeerCount(conferenceNumber)
	if err != nil {
		return nil, ErrFuncFail
	}

	for peerNumber := uint32(0); peerNumber < math.MaxUint32; peerNumber++ {
		pubkey, err := t.ConferencePeerGetPublicKey(conferenceNumber, peerNumber)
		if err != nil {
			break
		} else {
			peerPubkeys = append(peerPubkeys, pubkey)
		}
		if uint32(len(peerPubkeys)) >= peerCount {
			break
		}
	}
	return peerPubkeys, nil
}

// return [peerNumber]pubKey
func (t *Tox) ConferenceGetPeers(conferenceNumber uint32) (map[uint32]string, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	peers := make(map[uint32]string, 0)
	peerCount, err := t.ConferencePeerCount(conferenceNumber)
	if err != nil {
		return nil, ErrFuncFail
	}

	for peerNumber := uint32(0); peerNumber < math.MaxUint32; peerNumber++ {
		pubkey, err := t.ConferencePeerGetPublicKey(conferenceNumber, peerNumber)
		if err != nil {
			break
		} else {
			peers[peerNumber] = pubkey
		}
		if uint32(len(peers)) >= peerCount {
			break
		}
	}

	return peers, nil
}

// ConferenceGetChatlistSize
func (t *Tox) ConferenceGetChatlistSize() (uint32, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}
	ret := C.tox_conference_get_chatlist_size(t.tox)
	return uint32(ret), nil
}

func (t *Tox) ConferenceGetChatlist() ([]uint32, error) {
	if t.tox == nil {
		return nil, ErrToxInit
	}

	size, err := t.ConferenceGetChatlistSize()
	if err != nil {
		return nil, err
	}

	chatList := make([]uint32, size)

	if size > 0 {
		C.tox_conference_get_chatlist(t.tox, (*C.uint32_t)(&chatList[0]))
	}

	return chatList, nil
}

func (t *Tox) ConferenceGetType(conferenceNumber uint32) (int, error) {
	if t.tox == nil {
		return 0, ErrToxInit
	}

	var toxErrConferenceGetType C.Tox_Err_Conference_Get_Type
	ret := C.tox_conference_get_type(t.tox, (C.uint32_t)(conferenceNumber), &toxErrConferenceGetType)
	if ToxErrConferenceGetType(toxErrConferenceGetType) != TOX_ERR_CONFERENCE_GET_TYPE_OK {
		return int(ret), ErrFuncFail
	}

	return int(ret), nil
}

func (t *Tox) ConferenceGetIdentifier(conferenceNumber uint32) (string, error) {
	if t.tox == nil {
		return "", ErrToxInit
	}

	idbuf := [1 + C.TOX_PUBLIC_KEY_SIZE]byte{}

	C.tox_conference_get_id(t.tox, (C.uint32_t)(conferenceNumber), (*C.uint8_t)(&idbuf[0]))

	identifier := strings.ToUpper(hex.EncodeToString(idbuf[:]))
	identifier = identifier[2:] // 1B(type)+32B(identifier)

	return identifier, nil
}
