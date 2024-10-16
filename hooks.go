package dpc_tox

//#include <tox/tox.h>
import "C"
import "encoding/hex"
import "unsafe"
import "unicode/utf8"

//export hook_callback_self_connection_status
func hook_callback_self_connection_status(t unsafe.Pointer, status C.TOX_CONNECTION, tox unsafe.Pointer) {
	(*Tox)(tox).onSelfConnectionStatusChanges((*Tox)(tox), ToxConnection(status))
}

//export hook_callback_friend_name
func hook_callback_friend_name(t unsafe.Pointer, friendnumber C.uint32_t, name *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendNameChanges((*Tox)(tox), uint32(friendnumber), string(C.GoBytes(unsafe.Pointer(name), C.int(length))))
}

//export hook_callback_friend_status_message
func hook_callback_friend_status_message(t unsafe.Pointer, friendnumber C.uint32_t, message *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendStatusMessageChanges((*Tox)(tox), uint32(friendnumber), string(C.GoBytes(unsafe.Pointer(message), C.int(length))))
}

//export hook_callback_friend_status
func hook_callback_friend_status(t unsafe.Pointer, friendnumber C.uint32_t, status C.TOX_USER_STATUS, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendStatusChanges((*Tox)(tox), uint32(friendnumber), ToxUserStatus(status))
}

//export hook_callback_friend_connection_status
func hook_callback_friend_connection_status(t unsafe.Pointer, friendnumber C.uint32_t, status C.TOX_CONNECTION, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendConnectionStatusChanges((*Tox)(tox), uint32(friendnumber), ToxConnection(status))
}

//export hook_callback_friend_typing
func hook_callback_friend_typing(t unsafe.Pointer, friendnumber C.uint32_t, istyping C._Bool, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendTypingChanges((*Tox)(tox), uint32(friendnumber), bool(istyping))
}

//export hook_callback_friend_read_receipt
func hook_callback_friend_read_receipt(t unsafe.Pointer, friendnumber C.uint32_t, messageid C.uint32_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendReadReceipt((*Tox)(tox), uint32(friendnumber), uint32(messageid))
}

//export hook_callback_friend_request
func hook_callback_friend_request(t unsafe.Pointer, publicKey *C.uint8_t, message *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendRequest((*Tox)(tox), C.GoBytes((unsafe.Pointer)(publicKey), TOX_PUBLIC_KEY_SIZE), string(C.GoBytes(unsafe.Pointer(message), C.int(length))))
}

//export hook_callback_friend_message
func hook_callback_friend_message(t unsafe.Pointer, friendnumber C.uint32_t, messagetype C.TOX_MESSAGE_TYPE, message *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendMessage((*Tox)(tox), uint32(friendnumber), ToxMessageType(messagetype), string(C.GoBytes(unsafe.Pointer(message), C.int(length))))
}

//export hook_callback_file_recv_control
func hook_callback_file_recv_control(t unsafe.Pointer, friendnumber C.uint32_t, filenumber C.uint32_t, control C.TOX_FILE_CONTROL, tox unsafe.Pointer) {
	(*Tox)(tox).onFileRecvControl((*Tox)(tox), uint32(friendnumber), uint32(filenumber), ToxFileControl(control))
}

//export hook_callback_file_chunk_request
func hook_callback_file_chunk_request(t unsafe.Pointer, friendnumber C.uint32_t, filenumber C.uint32_t, position C.uint64_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFileChunkRequest((*Tox)(tox), uint32(friendnumber), uint32(filenumber), uint64(position), uint64(length))
}

//export hook_callback_file_recv
func hook_callback_file_recv(t unsafe.Pointer, friendnumber C.uint32_t, filenumber C.uint32_t, kind C.uint32_t, filesize C.uint64_t, filename *C.uint8_t, filenameLength C.size_t, tox unsafe.Pointer) {
	// convert the filename from CString to a GoString and encode hexadecimal if needed
	goFilenameBytes := C.GoBytes(unsafe.Pointer(filename), C.int(filenameLength))
	goFilename := string(goFilenameBytes)

	if !utf8.ValidString(goFilename) {
		goFilename = hex.EncodeToString(goFilenameBytes)
	}

	(*Tox)(tox).onFileRecv((*Tox)(tox), uint32(friendnumber), uint32(filenumber), ToxFileKind(kind), uint64(filesize), goFilename)
}

//export hook_callback_file_recv_chunk
func hook_callback_file_recv_chunk(t unsafe.Pointer, friendnumber C.uint32_t, filenumber C.uint32_t, position C.uint64_t, data *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFileRecvChunk((*Tox)(tox), uint32(friendnumber), uint32(filenumber), uint64(position), C.GoBytes((unsafe.Pointer)(data), C.int(length)))
}

//export hook_callback_friend_lossy_packet
func hook_callback_friend_lossy_packet(t unsafe.Pointer, friendnumber C.uint32_t, data *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendLossyPacket((*Tox)(tox), uint32(friendnumber), C.GoBytes((unsafe.Pointer)(data), C.int(length)))
}

//export hook_callback_friend_lossless_packet
func hook_callback_friend_lossless_packet(t unsafe.Pointer, friendnumber C.uint32_t, data *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onFriendLosslessPacket((*Tox)(tox), uint32(friendnumber), C.GoBytes((unsafe.Pointer)(data), (C.int)(length)))
}

//export hook_callback_conference_invite
func hook_callback_conference_invite(t unsafe.Pointer, friendnumber C.uint32_t, ctype C.Tox_Conference_Type, cookies *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onConferenceInvite((*Tox)(tox), uint32(friendnumber), ToxConferenceType(ctype), string(C.GoBytes(unsafe.Pointer(cookies), C.int(length))))
}

//export hook_callback_conference_connected
func hook_callback_conference_connected(t unsafe.Pointer, conferencenumber C.uint32_t, tox unsafe.Pointer) {
	(*Tox)(tox).onConferenceConnected((*Tox)(tox), uint32(conferencenumber))
}

//export hook_callback_conference_message
func hook_callback_conference_message(t unsafe.Pointer, conferencenumber C.uint32_t, peernumber C.uint32_t, messagetype C.Tox_Message_Type, message *C.uint8_t, length C.size_t, tox unsafe.Pointer) {
	(*Tox)(tox).onConferenceMessage((*Tox)(tox), uint32(conferencenumber), uint32(peernumber), ToxMessageType(messagetype), string(C.GoBytes(unsafe.Pointer(message), C.int(length))))
}

//export hook_callback_call
func hook_callback_call(t unsafe.Pointer, friendnumber C.uint32_t, audioenabled C._Bool, videoenabled C._Bool, toxav unsafe.Pointer) {
	(*ToxAV)(toxav).onCall((*ToxAV)(toxav), uint32(friendnumber), bool(audioenabled), bool(videoenabled))
}
