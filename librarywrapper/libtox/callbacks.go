package libtox

/*
#include <tox/tox.h>
#include "hooks-macro.c"
*/
import "C"
import "unsafe"

// OnSelfConnectionStatusChanges This event is triggered whenever there is a change in the DHT connectionstate.
/*
 * When disconnected, a client may choose to call tox_bootstrap again, to reconnect to the DHT.
 * Note that this state may frequently change for short amounts of time. Clients should therefore not immediately bootstrap onreceiving a disconnect.
 */
type OnSelfConnectionStatusChanges func(tox *Tox, status ToxConnection)

// OnFriendNameChanges This event is triggered when a friend changes their name.
type OnFriendNameChanges func(tox *Tox, friendnumber uint32, name []byte, length uint32)

// OnFriendStatusMessageChanges This event is triggered when a friend changes their status message.
type OnFriendStatusMessageChanges func(tox *Tox, friendnumber uint32, message []byte, length uint32)

// OnFriendStatusChanges This event is triggered when a friend changes their status message.
type OnFriendStatusChanges func(tox *Tox, friendnumber uint32, userstatus ToxUserStatus)

// OnFriendConnectionStatusChanges This event is triggered when a friend goes offline after having been online, or when a friend goes online.
// This callback is not called when adding friends. It is assumed that when adding friends, their connection status is initially offline.
type OnFriendConnectionStatusChanges func(tox *Tox, friendnumber uint32, connectionstatus ToxConnection)

// OnFriendTypingChanges This event is triggered when a friend starts or stops typing.
type OnFriendTypingChanges func(tox *Tox, friendnumber uint32, istyping bool)

// OnFriendReadReceipt This event is triggered when the friend receives the message with the corresponding message ID. */
type OnFriendReadReceipt func(tox *Tox, friendnumber uint32, messageid uint32)

// OnFriendRequest This event is triggered when a friend request is received.
type OnFriendRequest func(tox *Tox, publickey []byte, message []byte, length uint32)

// OnFriendMessage This event is triggered when a message from a friend is received.
type OnFriendMessage func(tox *Tox, friendnumber uint32, messagetype ToxMessageType, message []byte, length uint32)

// OnFileRecvControl This event is triggered when a file control command is received from a friend.
type OnFileRecvControl func(tox *Tox, friendnumber uint32, filenumber uint32, filecontrol ToxFileControl)

// OnFileChunkRequest This event is triggered when Core is ready to send more file data.
type OnFileChunkRequest func(tox *Tox, friendnumber uint32, filenumber uint32, position uint64, length uint64)

// OnFileRecv This event is triggered when a file transfer request is received.
type OnFileRecv func(tox *Tox, friendnumber uint32, filenumber uint32, kind ToxFileKind, filesize uint64, filename string, length uint32)

// OnFileRecvChunk This event is first triggered when a file transfer request is received, and subsequently when a chunk of file data for an accepted request was received.
type OnFileRecvChunk func(tox *Tox, friendnumber uint32, filenumber uint32, position uint64, data []byte, length uint32)

// OnFriendLossyPacket This event is triggered when a lossy packet is received from a friend.
type OnFriendLossyPacket func(tox *Tox, friendnumber uint32, data []byte, length uint32)

// OnFriendLosslessPacket This event is triggered when a lossless packet is received from a friend.
type OnFriendLosslessPacket func(tox *Tox, friendnumber uint32, data []byte, length uint32)

/*Conference callbacks*/

// OnConferenceInvite This event is triggered when the client is invited to join a conference.
type OnConferenceInvite func(tox *Tox, friendnumber uint32, conferencetype ToxConferenceType, cookie []byte)

// OnConferenceConnected This event is triggered when the client successfully connects to a conference after joining it with the tox_conference_join function.
type OnConferenceConnected func(tox *Tox, conferencenumber uint32)

// OnConferenceMessage This event is triggered when the client receives a conference message.
type OnConferenceMessage func(tox *Tox, conferencenumber uint32, peernumber uint32, messagetype ToxMessageType, message []byte, length uint32)

/*
 * Functions to register the callbacks.
 */

// CallbackSelfConnectionStatusChanges sets the function to be called when self connection status changed.
func (t *Tox) CallbackSelfConnectionStatusChanges(f OnSelfConnectionStatusChanges) {
	if t.Toxcore != nil {
		t.onSelfConnectionStatusChanges = f
		C.set_callback_self_connection_status(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendNameChanges sets the function to be called for friend's name changed.
func (t *Tox) CallbackFriendNameChanges(f OnFriendNameChanges) {
	if t.Toxcore != nil {
		t.onFriendNameChanges = f
		C.set_callback_friend_name(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendStatusMessageChanges sets the function to be called when friend's status message changed.
func (t *Tox) CallbackFriendStatusMessageChanges(f OnFriendStatusMessageChanges) {
	if t.Toxcore != nil {
		t.onFriendStatusMessageChanges = f
		C.set_callback_friend_status_message(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendStatusChanges sets the function to be called when friend's status changed.
func (t *Tox) CallbackFriendStatusChanges(f OnFriendStatusChanges) {
	if t.Toxcore != nil {
		t.onFriendStatusChanges = f
		C.set_callback_friend_status(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendConnectionStatusChanges sets the function to be called when friend's connection status changed.
func (t *Tox) CallbackFriendConnectionStatusChanges(f OnFriendConnectionStatusChanges) {
	if t.Toxcore != nil {
		t.onFriendConnectionStatusChanges = f
		C.set_callback_friend_connection_status(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendTypingChanges sets the function to be called when friend's typing changed.
func (t *Tox) CallbackFriendTypingChanges(f OnFriendTypingChanges) {
	if t.Toxcore != nil {
		t.onFriendTypingChanges = f
		C.set_callback_friend_typing(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendReadReceipt sets the function to be called when receiving read receipts.
func (t *Tox) CallbackFriendReadReceipt(f OnFriendReadReceipt) {
	if t.Toxcore != nil {
		t.onFriendReadReceipt = f
		C.set_callback_friend_read_receipt(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendRequest sets the function to be called when friend's request receipts.
func (t *Tox) CallbackFriendRequest(f OnFriendRequest) {
	if t.Toxcore != nil {
		t.onFriendRequest = f
		C.set_callback_friend_request(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendMessage sets the function to be called when receiving a friend message.
func (t *Tox) CallbackFriendMessage(f OnFriendMessage) {
	if t.Toxcore != nil {
		t.onFriendMessage = f
		C.set_callback_friend_message(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFileRecvControl sets the callback for file control requests.
func (t *Tox) CallbackFileRecvControl(f OnFileRecvControl) {
	if t.Toxcore != nil {
		t.onFileRecvControl = f
		C.set_callback_file_recv_control(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFileChunkRequest sets the callback to be called when tox is ready to send more file data.
func (t *Tox) CallbackFileChunkRequest(f OnFileChunkRequest) {
	if t.Toxcore != nil {
		t.onFileChunkRequest = f
		C.set_callback_file_chunk_request(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFileRecv sets the callback to be called when a file transfer request is received.
func (t *Tox) CallbackFileRecv(f OnFileRecv) {
	if t.Toxcore != nil {
		t.onFileRecv = f
		C.set_callback_file_recv(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFileRecvChunk sets the callback to be called when a file transfer request is received,
// and subsequently when a chunk of file data for an accepted request was received.
func (t *Tox) CallbackFileRecvChunk(f OnFileRecvChunk) {
	if t.Toxcore != nil {
		t.onFileRecvChunk = f
		C.set_callback_file_recv_chunk(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendLossyPacket sets the callback to be called when a lossy packet is received from a friend.
func (t *Tox) CallbackFriendLossyPacket(f OnFriendLossyPacket) {
	if t.Toxcore != nil {
		t.onFriendLossyPacket = f
		C.set_callback_friend_lossy_packet(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendLosslessPacket sets the callback to be called when a lossless packet is received from a friend.
func (t *Tox) CallbackFriendLosslessPacket(f OnFriendLosslessPacket) {
	if t.Toxcore != nil {
		t.onFriendLosslessPacket = f
		C.set_callback_friend_lossless_packet(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendLosslessPacket sets the callback to be called when the client is invited to join a conference.
func (t *Tox) CallbackConferenceInvite(f OnConferenceInvite) {
	if t.Toxcore != nil {
		t.onConferenceInvite = f
		C.set_callback_conference_invite(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackFriendLosslessPacket sets the callback to be called when the client receives a conference message.
func (t *Tox) CallbackConferenceMessage(f OnConferenceMessage) {
	if t.Toxcore != nil {
		t.onConferenceMessage = f
		C.set_callback_conference_message(t.Toxcore, unsafe.Pointer(t))
	}
}

// CallbackConferenceConnected sets the callback to be called when the client successfully connects to a conference
// after joining it with the tox_conference_join function.
func (t *Tox) CallbackConferenceConnected(f OnConferenceConnected) {
	if t.Toxcore != nil {
		t.onConferenceConnected = f
		C.set_callback_conference_connected(t.Toxcore, unsafe.Pointer(t))
	}
}
