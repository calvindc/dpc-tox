#include <tox/tox.h>

/* Macro defined:
 * Creates the C function to directly register a given callback from tox.h
 */
#define CREATE_HOOK(x) \
static void set_##x(Tox *tox, void *t) { \
  tox_##x(tox, hook_##x, t); \
}

//Tag: Headers for the exported GO functions in /libtox/hooks.go

/*
 * general callback functions
 */

//typedef void tox_self_connection_status_cb(Tox *tox, Tox_Connection connection_status, void *user_data);
void hook_callback_self_connection_status(Tox*, Tox_Connection, void*);

//typedef void tox_friend_name_cb(Tox *tox, Tox_Friend_Number friend_number, const uint8_t name[], size_t length, void *user_data);
void hook_callback_friend_name(Tox*, Tox_Friend_Number, const uint8_t*, size_t, void*);

//typedef void tox_friend_status_message_cb(Tox *tox, Tox_Friend_Number friend_number,const uint8_t message[], size_t length, void *user_data);
void hook_callback_friend_status_message(Tox*, Tox_Friend_Number, const uint8_t*, size_t, void*);

//typedef void tox_friend_status_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_User_Status status, void *user_data);
void hook_callback_friend_status(Tox*, Tox_Friend_Number, TOX_USER_STATUS, void*);

//typedef void tox_friend_connection_status_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_Connection connection_status, void *user_data);
void hook_callback_friend_connection_status(Tox*, Tox_Friend_Number, Tox_Connection, void*);

//typedef void tox_friend_typing_cb(Tox *tox, Tox_Friend_Number friend_number, bool typing, void *user_data);
void hook_callback_friend_typing(Tox*, Tox_Friend_Number, bool, void*);

//typedef void tox_friend_read_receipt_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_Friend_Message_Id message_id, void *user_data);
void hook_callback_friend_read_receipt(Tox*, Tox_Friend_Number, Tox_Friend_Message_Id, void*);

//typedef void tox_friend_request_cb(Tox *tox, const uint8_t public_key[TOX_PUBLIC_KEY_SIZE], const uint8_t message[], size_t length,void *user_data);
void hook_callback_friend_request(Tox*, const uint8_t*, const uint8_t*, size_t, void*);

//typedef void tox_friend_message_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_Message_Type type, const uint8_t message[], size_t length, void *user_data);
void hook_callback_friend_message(Tox*, Tox_Friend_Number, Tox_Message_Type, const uint8_t*, size_t, void*);

//typedef void tox_file_recv_control_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_File_Number file_number, Tox_File_Control control, void *user_data);
void hook_callback_file_recv_control(Tox*, Tox_Friend_Number, Tox_File_Number, TOX_FILE_CONTROL, void*);

//typedef void tox_file_chunk_request_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_File_Number file_number, uint64_t position,size_t length, void *user_data);
void hook_callback_file_chunk_request(Tox*, Tox_Friend_Number, Tox_File_Number, uint64_t, size_t, void*);

//typedef void tox_file_recv_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_File_Number file_number, uint32_t kind, uint64_t file_size, const uint8_t filename[], size_t filename_length, void *user_data);
void hook_callback_file_recv(Tox*, Tox_Friend_Number, Tox_File_Number, uint32_t, uint64_t, const uint8_t*, size_t, void*);

//typedef void tox_file_recv_chunk_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_File_Number file_number, uint64_t position, const uint8_t data[], size_t length, void *user_data);
void hook_callback_file_recv_chunk(Tox*, uint32_t, uint32_t, uint64_t, const uint8_t*, size_t, void*);

//typedef void tox_friend_lossy_packet_cb(Tox *tox, Tox_Friend_Number friend_number, const uint8_t data[], size_t length, void *user_data);
void hook_callback_friend_lossy_packet(Tox*, Tox_Friend_Number, const uint8_t*, size_t, void*);

//typedef void tox_friend_lossless_packet_cb(Tox *tox, Tox_Friend_Number friend_number,const uint8_t data[], size_t length,void *user_data);
void hook_callback_friend_lossless_packet(Tox*, Tox_Friend_Number, const uint8_t*, size_t, void*);

/*
 * conference callback functions
 */
//typedef void tox_conference_invite_cb(Tox *tox, Tox_Friend_Number friend_number, Tox_Conference_Type type, const uint8_t cookie[], size_t length, void *user_data);
void hook_callback_conference_invite(Tox*, Tox_Friend_Number, Tox_Conference_Type, const uint8_t*, size_t, void*);

//typedef void tox_conference_connected_cb(Tox *tox, Tox_Conference_Number conference_number, void *user_data);
void hook_callback_conference_connected(Tox*, Tox_Conference_Number, void*);

//typedef void tox_conference_message_cb(Tox *tox, Tox_Conference_Number conference_number, Tox_Conference_Peer_Number peer_number, Tox_Message_Type type, const uint8_t message[], size_t length, void *user_data);
void hook_callback_conference_message(Tox*, Tox_Conference_Number, Tox_Conference_Peer_Number, Tox_Message_Type, const uint8_t*, size_t, void*);

CREATE_HOOK(callback_self_connection_status)
CREATE_HOOK(callback_friend_name)
CREATE_HOOK(callback_friend_status_message)
CREATE_HOOK(callback_friend_status)
CREATE_HOOK(callback_friend_connection_status)
CREATE_HOOK(callback_friend_typing)
CREATE_HOOK(callback_friend_read_receipt)
CREATE_HOOK(callback_friend_request)
CREATE_HOOK(callback_friend_message)
CREATE_HOOK(callback_file_recv_control)
CREATE_HOOK(callback_file_chunk_request)
CREATE_HOOK(callback_file_recv)
CREATE_HOOK(callback_file_recv_chunk)
CREATE_HOOK(callback_friend_lossy_packet)
CREATE_HOOK(callback_friend_lossless_packet)

CREATE_HOOK(callback_conference_invite)
CREATE_HOOK(callback_conference_connected)
CREATE_HOOK(callback_conference_message)