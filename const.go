package dpc_tox

//#include <tox/tox.h>
import "C"
import "errors"

const (
	TOX_PUBLIC_KEY_SIZE           = C.TOX_PUBLIC_KEY_SIZE
	TOX_SECRET_KEY_SIZE           = C.TOX_SECRET_KEY_SIZE
	TOX_ADDRESS_SIZE              = C.TOX_ADDRESS_SIZE
	TOX_MAX_NAME_LENGTH           = C.TOX_MAX_NAME_LENGTH
	TOX_MAX_STATUS_MESSAGE_LENGTH = C.TOX_MAX_STATUS_MESSAGE_LENGTH
	TOX_MAX_FRIEND_REQUEST_LENGTH = C.TOX_MAX_FRIEND_REQUEST_LENGTH
	TOX_MAX_MESSAGE_LENGTH        = C.TOX_MAX_MESSAGE_LENGTH
	TOX_MAX_CUSTOM_PACKET_SIZE    = C.TOX_MAX_CUSTOM_PACKET_SIZE
	TOX_HASH_LENGTH               = C.TOX_HASH_LENGTH
	TOX_FILE_ID_LENGTH            = C.TOX_FILE_ID_LENGTH
	TOX_MAX_FILENAME_LENGTH       = C.TOX_MAX_FILENAME_LENGTH
)

type ToxUserStatus C.TOX_USER_STATUS

var (
	TOX_USERSTATUS_NONE ToxUserStatus = C.TOX_USER_STATUS_NONE //User is online and available.
	TOX_USERSTATUS_AWAY ToxUserStatus = C.TOX_USER_STATUS_AWAY //User is away. Clients can set this e.g. after a user defined inactivity time.
	TOX_USERSTATUS_BUSY ToxUserStatus = C.TOX_USER_STATUS_BUSY //User is busy. Signals to other clients that this client does not currently wish to communicate.
)

type ToxMessageType C.TOX_MESSAGE_TYPE

var (
	TOX_MESSAGE_TYPE_NORMAL ToxMessageType = C.TOX_MESSAGE_TYPE_NORMAL
	TOX_MESSAGE_TYPE_ACTION ToxMessageType = C.TOX_MESSAGE_TYPE_ACTION
)

type ToxProxyType C.TOX_PROXY_TYPE

var (
	TOX_PROXY_TYPE_NONE   ToxProxyType = C.TOX_PROXY_TYPE_NONE
	TOX_PROXY_TYPE_HTTP   ToxProxyType = C.TOX_PROXY_TYPE_HTTP
	TOX_PROXY_TYPE_SOCKS5 ToxProxyType = C.TOX_PROXY_TYPE_SOCKS5
)

type ToxSaveDataType C.TOX_SAVEDATA_TYPE

var (
	TOX_SAVEDATA_TYPE_NONE       ToxSaveDataType = C.TOX_SAVEDATA_TYPE_NONE
	TOX_SAVEDATA_TYPE_TOX_SAVE   ToxSaveDataType = C.TOX_SAVEDATA_TYPE_TOX_SAVE
	TOX_SAVEDATA_TYPE_SECRET_KEY ToxSaveDataType = C.TOX_SAVEDATA_TYPE_SECRET_KEY
)

type ToxConferenceType C.Tox_Conference_Type

var (
	TOX_CONFERENCE_TYPE_TEXT ToxConferenceType = C.TOX_CONFERENCE_TYPE_TEXT
	TOX_CONFERENCE_TYPE_AV   ToxConferenceType = C.TOX_CONFERENCE_TYPE_AV
)

type ToxErrOptionsNew C.TOX_ERR_OPTIONS_NEW

var (
	TOX_ERR_OPTIONS_NEW_OK     ToxErrOptionsNew = C.TOX_ERR_OPTIONS_NEW_OK
	TOX_ERR_OPTIONS_NEW_MALLOC ToxErrOptionsNew = C.TOX_ERR_OPTIONS_NEW_MALLOC
)

type ToxConnection C.Tox_Connection

var (
	TOX_CONNECTION_NONE ToxConnection = C.TOX_CONNECTION_NONE
	TOX_CONNECTION_TCP  ToxConnection = C.TOX_CONNECTION_TCP
	TOX_CONNECTION_UDP  ToxConnection = C.TOX_CONNECTION_UDP
)

type ToxFileKind C.TOX_FILE_KIND

var (
	TOX_FILE_KIND_DATA   ToxFileKind = C.TOX_FILE_KIND_DATA
	TOX_FILE_KIND_AVATAR ToxFileKind = C.TOX_FILE_KIND_AVATAR
)

type ToxFileControl C.TOX_FILE_CONTROL

var (
	TOX_FILE_CONTROL_RESUME ToxFileControl = C.TOX_FILE_CONTROL_RESUME
	TOX_FILE_CONTROL_PAUSE  ToxFileControl = C.TOX_FILE_CONTROL_PAUSE
	TOX_FILE_CONTROL_CANCEL ToxFileControl = C.TOX_FILE_CONTROL_CANCEL
)

/* === Errors === */
// General errors
var (
	ErrToxNew   = errors.New("Error initializing Tox")
	ErrToxInit  = errors.New("Tox not initialized")
	ErrArgs     = errors.New("Nil arguments or wrong size")
	ErrFuncFail = errors.New("Function failed")
	ErrUnknown  = errors.New("An unknown error occoured")
)

var (
	ErrNewMalloc        = errors.New("Memory allocation failed")
	ErrNewPortAlloc     = errors.New("Could not bind to port")
	ErrNewProxy         = errors.New("Invalid proxy configuration")
	ErrNewLoadEnc       = errors.New("The savedata is encrypted")
	ErrNewLoadBadFormat = errors.New("The savedata format is invalid")
)

type ToxErrNew C.TOX_ERR_NEW

var (
	TOX_ERR_NEW_OK              ToxErrNew = C.TOX_ERR_NEW_OK
	TOX_ERR_NEW_NULL            ToxErrNew = C.TOX_ERR_NEW_NULL
	TOX_ERR_NEW_MALLOC          ToxErrNew = C.TOX_ERR_NEW_MALLOC
	TOX_ERR_NEW_PORT_ALLOC      ToxErrNew = C.TOX_ERR_NEW_PORT_ALLOC
	TOX_ERR_NEW_PROXY_BAD_TYPE  ToxErrNew = C.TOX_ERR_NEW_PROXY_BAD_TYPE
	TOX_ERR_NEW_PROXY_BAD_HOST  ToxErrNew = C.TOX_ERR_NEW_PROXY_BAD_HOST
	TOX_ERR_NEW_PROXY_BAD_PORT  ToxErrNew = C.TOX_ERR_NEW_PROXY_BAD_PORT
	TOX_ERR_NEW_PROXY_NOT_FOUND ToxErrNew = C.TOX_ERR_NEW_PROXY_NOT_FOUND
	TOX_ERR_NEW_LOAD_ENCRYPTED  ToxErrNew = C.TOX_ERR_NEW_LOAD_ENCRYPTED
	TOX_ERR_NEW_LOAD_BAD_FORMAT ToxErrNew = C.TOX_ERR_NEW_LOAD_BAD_FORMAT
)

type ToxErrBootstrap C.TOX_ERR_BOOTSTRAP

var (
	TOX_ERR_BOOTSTRAP_OK       ToxErrBootstrap = C.TOX_ERR_BOOTSTRAP_OK
	TOX_ERR_BOOTSTRAP_NULL     ToxErrBootstrap = C.TOX_ERR_BOOTSTRAP_NULL
	TOX_ERR_BOOTSTRAP_BAD_HOST ToxErrBootstrap = C.TOX_ERR_BOOTSTRAP_BAD_HOST
	TOX_ERR_BOOTSTRAP_BAD_PORT ToxErrBootstrap = C.TOX_ERR_BOOTSTRAP_BAD_PORT
)

var (
	ErrFriendAddTooLong      = errors.New("Message too long")
	ErrFriendAddNoMessage    = errors.New("Empty message")
	ErrFriendAddOwnKey       = errors.New("Own key")
	ErrFriendAddAlreadySent  = errors.New("Already sent")
	ErrFriendAddBadChecksum  = errors.New("Bad checksum in address")
	ErrFriendAddSetNewNospam = errors.New("Different nospam")
	ErrFriendAddNoMem        = errors.New("Failed increasing friend list")
)

type ToxErrFriendAdd C.TOX_ERR_FRIEND_ADD

var (
	TOX_ERR_FRIEND_ADD_OK             ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_OK
	TOX_ERR_FRIEND_ADD_NULL           ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_NULL
	TOX_ERR_FRIEND_ADD_TOO_LONG       ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_TOO_LONG
	TOX_ERR_FRIEND_ADD_NO_MESSAGE     ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_NO_MESSAGE
	TOX_ERR_FRIEND_ADD_OWN_KEY        ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_OWN_KEY
	TOX_ERR_FRIEND_ADD_ALREADY_SENT   ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_ALREADY_SENT
	TOX_ERR_FRIEND_ADD_BAD_CHECKSUM   ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_BAD_CHECKSUM
	TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM
	TOX_ERR_FRIEND_ADD_MALLOC         ToxErrFriendAdd = C.TOX_ERR_FRIEND_ADD_MALLOC
)

type ToxErrFriendByPublicKey C.TOX_ERR_FRIEND_BY_PUBLIC_KEY

var (
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK        ToxErrFriendByPublicKey = C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL      ToxErrFriendByPublicKey = C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_NOT_FOUND ToxErrFriendByPublicKey = C.TOX_ERR_FRIEND_BY_PUBLIC_KEY_NOT_FOUND
)

type ToxErrFriendGetPublicKey C.TOX_ERR_FRIEND_GET_PUBLIC_KEY

var (
	TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK               ToxErrFriendGetPublicKey = C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK
	TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND ToxErrFriendGetPublicKey = C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND
)

type ToxErrFriendDelete C.TOX_ERR_FRIEND_DELETE

var (
	TOX_ERR_FRIEND_DELETE_OK               ToxErrFriendDelete = C.TOX_ERR_FRIEND_DELETE_OK
	TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND ToxErrFriendDelete = C.TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND
)

type ToxErrFriendQuery C.TOX_ERR_FRIEND_QUERY

var (
	TOX_ERR_FRIEND_QUERY_OK               ToxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_OK
	TOX_ERR_FRIEND_QUERY_NULL             ToxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_NULL
	TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND ToxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
)

type ToxErrSetInfo C.TOX_ERR_SET_INFO

var (
	TOX_ERR_SET_INFO_OK       ToxErrSetInfo = C.TOX_ERR_SET_INFO_OK
	TOX_ERR_SET_INFO_NULL     ToxErrSetInfo = C.TOX_ERR_SET_INFO_NULL
	TOX_ERR_SET_INFO_TOO_LONG ToxErrSetInfo = C.TOX_ERR_SET_INFO_TOO_LONG
)

type ToxErrSetTyping C.TOX_ERR_SET_TYPING

var (
	TOX_ERR_SET_TYPING_OK               ToxErrSetTyping = C.TOX_ERR_SET_TYPING_OK
	TOX_ERR_SET_TYPING_FRIEND_NOT_FOUND ToxErrSetTyping = C.TOX_ERR_SET_TYPING_FRIEND_NOT_FOUND
)

type ToxErrFriendSendMessage C.TOX_ERR_FRIEND_SEND_MESSAGE

var (
	TOX_ERR_FRIEND_SEND_MESSAGE_OK                   ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_OK
	TOX_ERR_FRIEND_SEND_MESSAGE_NULL                 ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_NULL
	TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_FOUND     ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_FOUND
	TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_CONNECTED ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_CONNECTED
	TOX_ERR_FRIEND_SEND_MESSAGE_SENDQ                ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_SENDQ
	TOX_ERR_FRIEND_SEND_MESSAGE_TOO_LONG             ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_TOO_LONG
	TOX_ERR_FRIEND_SEND_MESSAGE_EMPTY                ToxErrFriendSendMessage = C.TOX_ERR_FRIEND_SEND_MESSAGE_EMPTY
)

type ToxErrFriendGetLastOnline C.TOX_ERR_FRIEND_GET_LAST_ONLINE

var (
	TOX_ERR_FRIEND_GET_LAST_ONLINE_OK               ToxErrFriendGetLastOnline = C.TOX_ERR_FRIEND_GET_LAST_ONLINE_OK
	TOX_ERR_FRIEND_GET_LAST_ONLINE_FRIEND_NOT_FOUND ToxErrFriendGetLastOnline = C.TOX_ERR_FRIEND_GET_LAST_ONLINE_FRIEND_NOT_FOUND
)

type ToxErrFileControl C.TOX_ERR_FILE_CONTROL

var (
	TOX_ERR_FILE_CONTROL_OK                   ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_OK
	TOX_ERR_FILE_CONTROL_FRIEND_NOT_FOUND     ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_FRIEND_NOT_FOUND
	TOX_ERR_FILE_CONTROL_FRIEND_NOT_CONNECTED ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_CONTROL_NOT_FOUND            ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_NOT_FOUND
	TOX_ERR_FILE_CONTROL_NOT_PAUSED           ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_NOT_PAUSED
	TOX_ERR_FILE_CONTROL_DENIED               ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_DENIED
	TOX_ERR_FILE_CONTROL_ALREADY_PAUSED       ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_ALREADY_PAUSED
	TOX_ERR_FILE_CONTROL_SENDQ                ToxErrFileControl = C.TOX_ERR_FILE_CONTROL_SENDQ
)

type ToxErrFileSeek C.TOX_ERR_FILE_SEEK

var (
	TOX_ERR_FILE_SEEK_OK                   ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_OK
	TOX_ERR_FILE_SEEK_FRIEND_NOT_FOUND     ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEEK_FRIEND_NOT_CONNECTED ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEEK_NOT_FOUND            ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_NOT_FOUND
	TOX_ERR_FILE_SEEK_DENIED               ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_DENIED
	TOX_ERR_FILE_SEEK_INVALID_POSITION     ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_INVALID_POSITION
	TOX_ERR_FILE_SEEK_SENDQ                ToxErrFileSeek = C.TOX_ERR_FILE_SEEK_SENDQ
)

type ToxErrFileGet C.TOX_ERR_FILE_GET

var (
	TOX_ERR_FILE_GET_OK               ToxErrFileGet = C.TOX_ERR_FILE_GET_OK
	TOX_ERR_FILE_GET_NULL             ToxErrFileGet = C.TOX_ERR_FILE_GET_NULL
	TOX_ERR_FILE_GET_FRIEND_NOT_FOUND ToxErrFileGet = C.TOX_ERR_FILE_GET_FRIEND_NOT_FOUND
	TOX_ERR_FILE_GET_NOT_FOUND        ToxErrFileGet = C.TOX_ERR_FILE_GET_NOT_FOUND
)

var (
	ErrFileSendInvalidFileID = errors.New("The size of the given FileID is invalid.")
)

type ToxErrFileSend C.TOX_ERR_FILE_SEND

var (
	TOX_ERR_FILE_SEND_OK                   ToxErrFileSend = C.TOX_ERR_FILE_SEND_OK
	TOX_ERR_FILE_SEND_NULL                 ToxErrFileSend = C.TOX_ERR_FILE_SEND_NULL
	TOX_ERR_FILE_SEND_FRIEND_NOT_FOUND     ToxErrFileSend = C.TOX_ERR_FILE_SEND_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEND_FRIEND_NOT_CONNECTED ToxErrFileSend = C.TOX_ERR_FILE_SEND_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEND_NAME_TOO_LONG        ToxErrFileSend = C.TOX_ERR_FILE_SEND_NAME_TOO_LONG
	TOX_ERR_FILE_SEND_TOO_MANY             ToxErrFileSend = C.TOX_ERR_FILE_SEND_TOO_MANY
)

type ToxErrFileSendChunk C.TOX_ERR_FILE_SEND_CHUNK

var (
	TOX_ERR_FILE_SEND_CHUNK_OK                   ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_OK
	TOX_ERR_FILE_SEND_CHUNK_NULL                 ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_NULL
	TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_FOUND     ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_CONNECTED ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEND_CHUNK_NOT_FOUND            ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_NOT_FOUND
	TOX_ERR_FILE_SEND_CHUNK_NOT_TRANSFERRING     ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_NOT_TRANSFERRING
	TOX_ERR_FILE_SEND_CHUNK_INVALID_LENGTH       ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_INVALID_LENGTH
	TOX_ERR_FILE_SEND_CHUNK_SENDQ                ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_SENDQ
	TOX_ERR_FILE_SEND_CHUNK_WRONG_POSITION       ToxErrFileSendChunk = C.TOX_ERR_FILE_SEND_CHUNK_WRONG_POSITION
)

type ToxErrFriendCustomPacket C.TOX_ERR_FRIEND_CUSTOM_PACKET

var (
	TOX_ERR_FRIEND_CUSTOM_PACKET_OK                   ToxErrFriendCustomPacket = C.TOX_ERR_FRIEND_CUSTOM_PACKET_OK
	TOX_ERR_FRIEND_CUSTOM_PACKET_NULL                 ToxErrFriendCustomPacket = C.TOX_ERR_FRIEND_CUSTOM_PACKET_NULL
	TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_FOUND     ToxErrFriendCustomPacket = C.TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_FOUND
	TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_CONNECTED ToxErrFriendCustomPacket = C.TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_CONNECTED
	TOX_ERR_FRIEND_CUSTOM_PACKET_INVALID              ToxErrFriendCustomPacket = C.TOX_ERR_FRIEND_CUSTOM_PACKET_INVALID
	TOX_ERR_FRIEND_CUSTOM_PACKET_EMPTY                ToxErrFriendCustomPacket = C.TOX_ERR_FRIEND_CUSTOM_PACKET_EMPTY
	TOX_ERR_FRIEND_CUSTOM_PACKET_TOO_LONG             ToxErrFriendCustomPacket = C.TOX_ERR_FRIEND_CUSTOM_PACKET_TOO_LONG
	TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ                ToxErrFriendCustomPacket = C.TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ
)

type ToxErrGetPort C.TOX_ERR_GET_PORT

var (
	TOX_ERR_GET_PORT_OK        ToxErrGetPort = C.TOX_ERR_GET_PORT_OK
	TOX_ERR_GET_PORT_NOT_BOUND ToxErrGetPort = C.TOX_ERR_GET_PORT_NOT_BOUND
)

var (
	ErrConferenceNewFailedInitialize = errors.New("conference instance failed to initialize")
)

type ToxErrConferenceNew C.Tox_Err_Conference_New

var (
	TOX_ERR_CONFERENCE_NEW_OK   ToxErrConferenceNew = C.TOX_ERR_CONFERENCE_NEW_OK   //The function returned successfully.
	TOX_ERR_CONFERENCE_NEW_INIT ToxErrConferenceNew = C.TOX_ERR_CONFERENCE_NEW_INIT //The conference instance failed to initialize.
)

var (
	ErrConferenceDeleteFailed             = errors.New("delete conference failed")
	ErrConferenceDeleteConferenceNotFound = errors.New("the conference number passed did not designate a valid conference.")
)

type ToxErrConferenceDelete C.Tox_Err_Conference_Delete

var (
	TOX_ERR_CONFERENCE_DELETE_OK                   ToxErrConferenceDelete = C.TOX_ERR_CONFERENCE_DELETE_OK
	TOX_ERR_CONFERENCE_DELETE_CONFERENCE_NOT_FOUND ToxErrConferenceDelete = C.TOX_ERR_CONFERENCE_DELETE_CONFERENCE_NOT_FOUND //The conference number passed did not designate a valid conference.
)

type ToxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query

var (
	TOX_ERR_CONFERENCE_PEER_QUERY_OK                   ToxErrConferencePeerQuery = C.TOX_ERR_CONFERENCE_PEER_QUERY_OK
	TOX_ERR_CONFERENCE_PEER_QUERY_CONFERENCE_NOT_FOUND ToxErrConferencePeerQuery = C.TOX_ERR_CONFERENCE_PEER_QUERY_CONFERENCE_NOT_FOUND
	TOX_ERR_CONFERENCE_PEER_QUERY_PEER_NOT_FOUND       ToxErrConferencePeerQuery = C.TOX_ERR_CONFERENCE_PEER_QUERY_PEER_NOT_FOUND
	TOX_ERR_CONFERENCE_PEER_QUERY_NO_CONNECTION        ToxErrConferencePeerQuery = C.TOX_ERR_CONFERENCE_PEER_QUERY_NO_CONNECTION
)

var (
	ErrConferenceInviteConferenceNotFound = errors.New("The conference number passed did not designate a valid conference")
	ErrConferenceInviteFailSend           = errors.New("The invite packet failed to send")
	ErrConferenceInviteNoConnection       = errors.New("The client is not connected to the conference")
)

type ToxErrConferenceInvite C.Tox_Err_Conference_Invite

var (
	TOX_ERR_CONFERENCE_INVITE_OK                   ToxErrConferenceInvite = C.TOX_ERR_CONFERENCE_INVITE_OK
	TOX_ERR_CONFERENCE_INVITE_CONFERENCE_NOT_FOUND ToxErrConferenceInvite = C.TOX_ERR_CONFERENCE_INVITE_CONFERENCE_NOT_FOUND
	TOX_ERR_CONFERENCE_INVITE_FAIL_SEND            ToxErrConferenceInvite = C.TOX_ERR_CONFERENCE_INVITE_FAIL_SEND
	TOX_ERR_CONFERENCE_INVITE_NO_CONNECTION        ToxErrConferenceInvite = C.TOX_ERR_CONFERENCE_INVITE_NO_CONNECTION
)

type ToxErrConferenceSendMessage C.Tox_Err_Conference_Send_Message

var (
	TOX_ERR_CONFERENCE_SEND_MESSAGE_OK                   ToxErrConferenceSendMessage = C.TOX_ERR_CONFERENCE_SEND_MESSAGE_OK
	TOX_ERR_CONFERENCE_SEND_MESSAGE_CONFERENCE_NOT_FOUND ToxErrConferenceSendMessage = C.TOX_ERR_CONFERENCE_SEND_MESSAGE_CONFERENCE_NOT_FOUND
	TOX_ERR_CONFERENCE_SEND_MESSAGE_TOO_LONG             ToxErrConferenceSendMessage = C.TOX_ERR_CONFERENCE_SEND_MESSAGE_TOO_LONG
	TOX_ERR_CONFERENCE_SEND_MESSAGE_NO_CONNECTION        ToxErrConferenceSendMessage = C.TOX_ERR_CONFERENCE_SEND_MESSAGE_NO_CONNECTION
	TOX_ERR_CONFERENCE_SEND_MESSAGE_FAIL_SEND            ToxErrConferenceSendMessage = C.TOX_ERR_CONFERENCE_SEND_MESSAGE_FAIL_SEND
)

type ToxErrConferenceTitle C.Tox_Err_Conference_Title

var (
	TOX_ERR_CONFERENCE_TITLE_OK                   ToxErrConferenceTitle = C.TOX_ERR_CONFERENCE_TITLE_OK
	TOX_ERR_CONFERENCE_TITLE_CONFERENCE_NOT_FOUND ToxErrConferenceTitle = C.TOX_ERR_CONFERENCE_TITLE_CONFERENCE_NOT_FOUND
	TOX_ERR_CONFERENCE_TITLE_INVALID_LENGTH       ToxErrConferenceTitle = C.TOX_ERR_CONFERENCE_TITLE_INVALID_LENGTH
	TOX_ERR_CONFERENCE_TITLE_FAIL_SEND            ToxErrConferenceTitle = C.TOX_ERR_CONFERENCE_TITLE_FAIL_SEND
)

type ToxErrConferenceGetType C.Tox_Err_Conference_Get_Type

var (
	TOX_ERR_CONFERENCE_GET_TYPE_OK                   ToxErrConferenceGetType = C.TOX_ERR_CONFERENCE_GET_TYPE_OK
	TOX_ERR_CONFERENCE_GET_TYPE_CONFERENCE_NOT_FOUND ToxErrConferenceGetType = C.TOX_ERR_CONFERENCE_GET_TYPE_CONFERENCE_NOT_FOUND
)
