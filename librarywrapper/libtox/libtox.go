package libtox

//#cgo LDFLAGS: -ltoxcore
//#include <tox/tox.h>
//#include <stdlib.h>
import "C"
import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// Tox is Tox instance type.
/*
 * All the state associated with a connection is held within the instance. Multiple instances can exist and operate concurrently.
 * The maximum number of Tox instances that can exist on a single network device is limited.
 * Note that this is not just a per-process limit, since the limiting factor is the number of usable ports on a device.
 */
type Tox struct {
	cOptions *C.struct_Tox_Options
	Toxcore  *C.Tox
	mtx      sync.Mutex

	// Callbacks
	onSelfConnectionStatusChanges   OnSelfConnectionStatusChanges
	onFriendNameChanges             OnFriendNameChanges
	onFriendStatusMessageChanges    OnFriendStatusMessageChanges
	onFriendStatusChanges           OnFriendStatusChanges
	onFriendConnectionStatusChanges OnFriendConnectionStatusChanges
	onFriendTypingChanges           OnFriendTypingChanges
	onFriendReadReceipt             OnFriendReadReceipt
	onFriendRequest                 OnFriendRequest
	onFriendMessage                 OnFriendMessage
	onFileRecvControl               OnFileRecvControl
	onFileChunkRequest              OnFileChunkRequest
	onFileRecv                      OnFileRecv
	onFileRecvChunk                 OnFileRecvChunk
	onFriendLossyPacket             OnFriendLossyPacket
	onFriendLosslessPacket          OnFriendLosslessPacket

	onConferenceInvite    OnConferenceInvite
	onConferenceMessage   OnConferenceMessage
	onConferenceConnected OnConferenceConnected
}

// Options tox option params
type Options struct {
	/* The type of socket to create.
	 * If IPv6Enabled is true, both IPv6 and IPv4 connections are allowed.
	 */
	IPv6Enabled bool

	/* Enable the use of UDP communication when available.
	 *
	 * Setting this to false will force Tox to use TCP only. Communications will
	 * need to be relayed through a TCP relay node, potentially slowing them down.
	 * Disabling UDP support is necessary when using anonymous proxies or Tor.
	 */
	UDPEnabled bool

	/* The type of the proxy (PROXY_TYPE_NONE, PROXY_TYPE_HTTP or PROXY_TYPE_SOCKS5). */
	ProxyType ToxProxyType

	/* The IP address or DNS name of the proxy to be used. */
	ProxyHost string

	/* The port to use to connect to the proxy server. */
	ProxyPort uint16

	/* The start port of the inclusive port range to attempt to use. */
	StartPort uint16

	/* The end port of the inclusive port range to attempt to use. */
	EndPort uint16

	/* The port to use for the TCP server. If 0, the tcp server is disabled. */
	TcpPort uint16

	/* The type of savedata to load from. */
	SaveDataType ToxSaveDataType

	/* The savedata. */
	SaveData []byte
}

//=================
/* VersionMajor returns the major version number of the used Tox library */
func VersionMajor() uint32 {
	return uint32(C.tox_version_major())
}

/* VersionMinor returns the minor version number of the used Tox library */
func VersionMinor() uint32 {
	return uint32(C.tox_version_minor())
}

/* VersionPatch returns the patch number of the used Tox library */
func VersionPatch() uint32 {
	return uint32(C.tox_version_patch())
}

/* VersionIsCompatible returns whether the compiled Tox library version is
 * compatible with the passed version numbers. */
func VersionIsCompatible(major uint32, minor uint32, patch uint32) bool {
	return bool(C.tox_version_is_compatible((C.uint32_t)(major), (C.uint32_t)(minor), (C.uint32_t)(patch)))
}

/* New creates and initialises a new Tox instance and returns the corresponding
 * gotox instance. */
func New(options *Options) (*Tox, error) {
	var cTox *C.Tox
	var toxErrNew C.TOX_ERR_NEW
	var toxErrOptionsNew C.TOX_ERR_OPTIONS_NEW

	var cOptions *C.struct_Tox_Options = C.tox_options_new(&toxErrOptionsNew)
	if cOptions == nil || ToxErrOptionsNew(toxErrOptionsNew) != TOX_ERR_OPTIONS_NEW_OK {
		return nil, ErrFuncFail
	}

	if options == nil {
		cOptions = nil
	} else {
		// map options from Options to C.Tox_Options
		cOptions.ipv6_enabled = C.bool(options.IPv6Enabled)
		cOptions.udp_enabled = C.bool(options.UDPEnabled)

		var cProxyType C.TOX_PROXY_TYPE = C.TOX_PROXY_TYPE_NONE
		if options.ProxyType == TOX_PROXY_TYPE_HTTP {
			cProxyType = C.TOX_PROXY_TYPE_HTTP
		} else if options.ProxyType == TOX_PROXY_TYPE_SOCKS5 {
			cProxyType = C.TOX_PROXY_TYPE_SOCKS5
		}
		cOptions.proxy_type = cProxyType

		// max ProxyHost length is 255
		if len(options.ProxyHost) > 255 {
			return nil, ErrArgs
		}
		cProxyHost := C.CString(options.ProxyHost)
		cOptions.proxy_host = cProxyHost
		defer C.free(unsafe.Pointer(cProxyHost))

		cOptions.proxy_port = C.uint16_t(options.ProxyPort)
		cOptions.start_port = C.uint16_t(options.StartPort)
		cOptions.end_port = C.uint16_t(options.EndPort)
		cOptions.tcp_port = C.uint16_t(options.TcpPort)

		if options.SaveDataType == TOX_SAVEDATA_TYPE_TOX_SAVE {
			cOptions.savedata_type = C.TOX_SAVEDATA_TYPE_TOX_SAVE
		} else if options.SaveDataType == TOX_SAVEDATA_TYPE_SECRET_KEY {
			cOptions.savedata_type = C.TOX_SAVEDATA_TYPE_SECRET_KEY
		}

		if len(options.SaveData) > 0 {
			cOptions.savedata_data = (*C.uint8_t)(&options.SaveData[0])
		} else {
			cOptions.savedata_data = nil
		}

		cOptions.savedata_length = C.size_t(len(options.SaveData))
	}

	cTox = C.tox_new(cOptions, &toxErrNew)
	if cTox == nil || ToxErrNew(toxErrNew) != TOX_ERR_NEW_OK {
		C.tox_options_free(cOptions)
		switch ToxErrNew(toxErrNew) {
		case TOX_ERR_NEW_NULL:
			return nil, ErrArgs
		case TOX_ERR_NEW_MALLOC:
			return nil, ErrNewMalloc
		case TOX_ERR_NEW_PORT_ALLOC:
			return nil, ErrNewPortAlloc
		case TOX_ERR_NEW_PROXY_BAD_TYPE:
			return nil, ErrNewProxy
		case TOX_ERR_NEW_PROXY_BAD_HOST:
			return nil, ErrNewProxy
		case TOX_ERR_NEW_PROXY_BAD_PORT:
			return nil, ErrNewProxy
		case TOX_ERR_NEW_PROXY_NOT_FOUND:
			return nil, ErrNewProxy
		case TOX_ERR_NEW_LOAD_ENCRYPTED:
			return nil, ErrNewLoadEnc
		case TOX_ERR_NEW_LOAD_BAD_FORMAT:
			return nil, ErrNewLoadBadFormat
		}

		if cTox == nil {
			return nil, ErrToxNew
		}

		return nil, ErrUnknown
	}

	t := &Tox{Toxcore: cTox, cOptions: cOptions}
	return t, nil
}

/* Kill releases all resources associated with the Tox instance and disconnects
 * from the network.
 * After calling this function `t *TOX` becomes invalid. Do not use it again! */
func (t *Tox) Kill() error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	C.tox_options_free(t.cOptions)
	C.tox_kill(t.Toxcore)

	return nil
}

/* GetSaveDataSize returns the size of the savedata returned by GetSavedata. */
func (t *Tox) GetSaveDataSize() (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	return uint32(C.tox_get_savedata_size(t.Toxcore)), nil
}

/* GetSavedata returns a byte slice of all information associated with the tox
 * instance. */
func (t *Tox) GetSavedata() ([]byte, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}
	size, err := t.GetSaveDataSize()
	if err != nil || size == 0 {
		return nil, ErrFuncFail
	}

	data := make([]byte, size)

	if size > 0 {
		C.tox_get_savedata(t.Toxcore, (*C.uint8_t)(&data[0]))
	}

	return data, nil
}

/* Bootstrap sends a "get nodes" request to the given bootstrap node with IP,
 * port, and public key to setup connections. */
func (t *Tox) Bootstrap(address string, port uint16, publickey []byte) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return ErrArgs
	}

	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))

	var toxErrBootstrap C.TOX_ERR_BOOTSTRAP
	success := C.tox_bootstrap(t.Toxcore, caddr, (C.uint16_t)(port), (*C.uint8_t)(&publickey[0]), &toxErrBootstrap)

	switch ToxErrBootstrap(toxErrBootstrap) {
	case TOX_ERR_BOOTSTRAP_OK:
		return nil
	case TOX_ERR_BOOTSTRAP_NULL:
		return ErrArgs
	case TOX_ERR_BOOTSTRAP_BAD_HOST:
		return ErrFuncFail
	case TOX_ERR_BOOTSTRAP_BAD_PORT:
		return ErrFuncFail
	}

	if !bool(success) {
		return ErrFuncFail
	}

	return ErrUnknown
}

/* AddTCPRelay adds the given node with IP, port, and public key without using
 * it as a boostrap node. */
func (t *Tox) AddTCPRelay(address string, port uint16, publickey []byte) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return ErrArgs
	}

	caddr := C.CString(address)
	defer C.free(unsafe.Pointer(caddr))

	var toxErrBootstrap C.TOX_ERR_BOOTSTRAP
	success := C.tox_add_tcp_relay(t.Toxcore, caddr, (C.uint16_t)(port), (*C.uint8_t)(&publickey[0]), &toxErrBootstrap)

	switch ToxErrBootstrap(toxErrBootstrap) {
	case TOX_ERR_BOOTSTRAP_OK:
		return nil
	case TOX_ERR_BOOTSTRAP_NULL:
		return ErrArgs
	case TOX_ERR_BOOTSTRAP_BAD_HOST:
		return ErrFuncFail
	case TOX_ERR_BOOTSTRAP_BAD_PORT:
		return ErrFuncFail
	}

	if !bool(success) {
		return ErrFuncFail
	}

	return ErrUnknown
}

/* SelfGetConnectionStatus returns true if Tox is connected to the DHT. */
func (t *Tox) SelfGetConnectionStatus() (ToxConnection, error) {
	if t.Toxcore == nil {
		return TOX_CONNECTION_NONE, ErrToxInit
	}

	return ToxConnection(C.tox_self_get_connection_status(t.Toxcore)), nil
}

/* IterationInterval returns the time in milliseconds before Iterate() should be
 * called again. */
func (t *Tox) IterationInterval() (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	ret := C.tox_iteration_interval(t.Toxcore)

	return uint32(ret), nil
}

/* Iterate is the main loop. It needs to be called every IterationInterval()
 * milliseconds. */
func (t *Tox) Iterate() error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	t.mtx.Lock()
	C.tox_iterate(t.Toxcore, unsafe.Pointer(t))
	t.mtx.Unlock()

	return nil
}

/* SelfGetAddress returns the public address to give to others. */
func (t *Tox) SelfGetAddress() ([]byte, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}

	address := make([]byte, TOX_ADDRESS_SIZE)
	C.tox_self_get_address(t.Toxcore, (*C.uint8_t)(&address[0]))

	return address, nil
}

/* SelfSetNospam sets the nospam of your ID. */
func (t *Tox) SelfSetNospam(nospam uint32) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	C.tox_self_set_nospam(t.Toxcore, (C.uint32_t)(nospam))
	return nil
}

/* SelfGetNospam returns the nospam of your ID. */
func (t *Tox) SelfGetNospam() (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	n := C.tox_self_get_nospam(t.Toxcore)
	return uint32(n), nil
}

/* SelfGetPublicKey returns the publickey of your profile. */
func (t *Tox) SelfGetPublicKey() ([]byte, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}

	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)

	C.tox_self_get_public_key(t.Toxcore, (*C.uint8_t)(&publickey[0]))
	return publickey, nil
}

/* SelfGetSecretKey returns the secretkey of your profile. */
func (t *Tox) SelfGetSecretKey() ([]byte, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}

	secretkey := make([]byte, TOX_SECRET_KEY_SIZE)

	C.tox_self_get_secret_key(t.Toxcore, (*C.uint8_t)(&secretkey[0]))
	return secretkey, nil
}

/* SelfSetName sets your nickname. The maximum name length is MAX_NAME_LENGTH. */
func (t *Tox) SelfSetName(name string) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var cName (*C.uint8_t)

	if len(name) == 0 {
		cName = nil
	} else {
		cName = (*C.uint8_t)(&[]byte(name)[0])
	}

	var setInfoError C.TOX_ERR_SET_INFO = C.TOX_ERR_SET_INFO_OK
	success := C.tox_self_set_name(t.Toxcore, cName, (C.size_t)(len(name)), &setInfoError)
	if !bool(success) || ToxErrSetInfo(setInfoError) != TOX_ERR_SET_INFO_OK {
		return ErrFuncFail
	}

	return nil
}

/* SelfGetNameSize returns the length of your name. */
func (t *Tox) SelfGetNameSize() (int64, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	ret := C.tox_self_get_name_size(t.Toxcore)

	return int64(ret), nil
}

/* SelfGetName returns your nickname. */
func (t *Tox) SelfGetName() (string, error) {
	if t.Toxcore == nil {
		return "", ErrToxInit
	}

	length, err := t.SelfGetNameSize()
	if err != nil {
		return "", ErrFuncFail
	}

	name := make([]byte, length)

	if length > 0 {
		C.tox_self_get_name(t.Toxcore, (*C.uint8_t)(&name[0]))
	}

	return string(name), nil
}

/* SelfSetStatusMessage sets your status message.
 * The maximum status length is MAX_STATUS_MESSAGE_LENGTH. */
func (t *Tox) SelfSetStatusMessage(status string) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var cStatus (*C.uint8_t)

	if len(status) == 0 {
		cStatus = nil
	} else {
		cStatus = (*C.uint8_t)(&[]byte(status)[0])
	}

	var setInfoError C.TOX_ERR_SET_INFO = C.TOX_ERR_SET_INFO_OK
	C.tox_self_set_status_message(t.Toxcore, cStatus, (C.size_t)(len(status)), &setInfoError)

	if ToxErrSetInfo(setInfoError) != TOX_ERR_SET_INFO_OK {
		return ErrFuncFail
	}

	return nil
}

/* SelfGetStatusMessageSize returns the size of your status message. */
func (t *Tox) SelfGetStatusMessageSize() (int64, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	ret := C.tox_self_get_status_message_size(t.Toxcore)

	return int64(ret), nil
}

/* SelfGetStatusMessage returns your status message. */
func (t *Tox) SelfGetStatusMessage() (string, error) {
	if t.Toxcore == nil {
		return "", ErrToxInit
	}

	length, err := t.SelfGetStatusMessageSize()
	if err != nil {
		return "", ErrFuncFail
	}

	statusMessage := make([]byte, length)

	if length > 0 {
		C.tox_self_get_status_message(t.Toxcore, (*C.uint8_t)(&statusMessage[0]))
	}

	return string(statusMessage), nil
}

/* SelfSetStatus sets your userstatus. */
func (t *Tox) SelfSetStatus(userstatus ToxUserStatus) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	C.tox_self_set_status(t.Toxcore, (C.TOX_USER_STATUS)(userstatus))

	return nil
}

/* SelfGetStatus returns your status. */
func (t *Tox) SelfGetStatus() (ToxUserStatus, error) {
	if t.Toxcore == nil {
		return TOX_USERSTATUS_NONE, ErrToxInit
	}

	n := C.tox_self_get_status(t.Toxcore)

	return ToxUserStatus(n), nil
}

/* FriendAdd adds a friend by sending a friend request containing the given
 * message.
 * Returns the friend number on success, or a ToxErrFriendAdd on failure.
 */
func (t *Tox) FriendAdd(address []byte, message string) (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	if len(address) != TOX_ADDRESS_SIZE || len(message) == 0 {
		return 0, ErrArgs
	}

	caddr := (*C.uint8_t)(&address[0])
	cmessage := (*C.uint8_t)(&[]byte(message)[0])

	var toxErrFriendAdd C.TOX_ERR_FRIEND_ADD
	ret := C.tox_friend_add(t.Toxcore, caddr, cmessage, (C.size_t)(len(message)), &toxErrFriendAdd)

	switch ToxErrFriendAdd(toxErrFriendAdd) {
	case TOX_ERR_FRIEND_ADD_OK:
		return uint32(ret), nil
	case TOX_ERR_FRIEND_ADD_NULL:
		return uint32(ret), ErrArgs
	case TOX_ERR_FRIEND_ADD_TOO_LONG:
		return uint32(ret), ErrFriendAddTooLong
	case TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		return uint32(ret), ErrFriendAddNoMessage
	case TOX_ERR_FRIEND_ADD_OWN_KEY:
		return uint32(ret), ErrFriendAddOwnKey
	case TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		return uint32(ret), ErrFriendAddAlreadySent
	case TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		return uint32(ret), ErrFriendAddBadChecksum
	case TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		return uint32(ret), ErrFriendAddSetNewNospam
	case TOX_ERR_FRIEND_ADD_MALLOC:
		return uint32(ret), ErrFriendAddNoMem
	default:
		return uint32(ret), ErrFuncFail
	}

	return uint32(ret), ErrUnknown
}

/* FriendAddNorequest adds a friend without sending a friend request.
 * Returns the friend number on success.
 */
func (t *Tox) FriendAddNorequest(publickey []byte) (uint32, error) {
	if t.Toxcore == nil {
		return C.UINT32_MAX, ErrToxInit
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return C.UINT32_MAX, ErrArgs
	}

	var toxErrFriendAdd C.TOX_ERR_FRIEND_ADD
	ret := C.tox_friend_add_norequest(t.Toxcore, (*C.uint8_t)(&publickey[0]), &toxErrFriendAdd)
	if ret == C.UINT32_MAX {
		return C.UINT32_MAX, ErrFuncFail
	}

	switch ToxErrFriendAdd(toxErrFriendAdd) {
	case TOX_ERR_FRIEND_ADD_OK:
		return uint32(ret), nil
	case TOX_ERR_FRIEND_ADD_NULL:
		return uint32(ret), ErrArgs
	case TOX_ERR_FRIEND_ADD_TOO_LONG:
		return uint32(ret), ErrFriendAddTooLong
	case TOX_ERR_FRIEND_ADD_NO_MESSAGE:
		return uint32(ret), ErrFriendAddNoMessage
	case TOX_ERR_FRIEND_ADD_OWN_KEY:
		return uint32(ret), ErrFriendAddOwnKey
	case TOX_ERR_FRIEND_ADD_ALREADY_SENT:
		return uint32(ret), ErrFriendAddAlreadySent
	case TOX_ERR_FRIEND_ADD_BAD_CHECKSUM:
		return uint32(ret), ErrFriendAddBadChecksum
	case TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM:
		return uint32(ret), ErrFriendAddSetNewNospam
	case TOX_ERR_FRIEND_ADD_MALLOC:
		return uint32(ret), ErrFriendAddNoMem
	default:
		return uint32(ret), ErrFuncFail
	}

	return uint32(ret), ErrUnknown
}

/* FriendDelete removes a friend. */
func (t *Tox) FriendDelete(friendNumber uint32) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var toxErrFriendDelete C.TOX_ERR_FRIEND_DELETE = C.TOX_ERR_FRIEND_DELETE_OK
	C.tox_friend_delete(t.Toxcore, (C.uint32_t)(friendNumber), &toxErrFriendDelete)

	switch ToxErrFriendDelete(toxErrFriendDelete) {
	case TOX_ERR_FRIEND_DELETE_OK:
		return nil
	case TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND:
		return ErrArgs
	default:
		return ErrFuncFail
	}

	return ErrUnknown
}

/* FriendByPublicKey returns the friend number associated to a given publickey. */
func (t *Tox) FriendByPublicKey(publickey []byte) (uint32, error) {
	if t.Toxcore == nil {
		return C.UINT32_MAX, ErrToxInit
	}

	if len(publickey) != TOX_PUBLIC_KEY_SIZE {
		return C.UINT32_MAX, ErrArgs
	}

	var toxErrFriendByPublicKey C.TOX_ERR_FRIEND_BY_PUBLIC_KEY
	n := C.tox_friend_by_public_key(t.Toxcore, (*C.uint8_t)(&publickey[0]), &toxErrFriendByPublicKey)

	switch ToxErrFriendByPublicKey(toxErrFriendByPublicKey) {
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK:
		return uint32(n), nil
	case TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL:
		return uint32(n), ErrArgs
	default:
		return uint32(n), ErrFuncFail
	}

	return uint32(n), ErrUnknown
}

/* FriendExists returns true if a friend exists with given friendNumber. */
func (t *Tox) FriendExists(friendNumber uint32) (bool, error) {
	if t.Toxcore == nil {
		return false, ErrToxInit
	}

	success := C.tox_friend_exists(t.Toxcore, (C.uint32_t)(friendNumber))

	return bool(success), nil
}

/* SelfGetFriendlistSize returns the number of friends on the friendlist. */
func (t *Tox) SelfGetFriendlistSize() (int64, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}
	n := C.tox_self_get_friend_list_size(t.Toxcore)

	return int64(n), nil
}

/* SelfGetFriendlist returns a slice of uint32 containing the friendNumbers. */
func (t *Tox) SelfGetFriendlist() ([]uint32, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}

	size, err := t.SelfGetFriendlistSize()
	if err != nil {
		return nil, ErrFuncFail
	}

	friendlist := make([]uint32, size)

	if size > 0 {
		C.tox_self_get_friend_list(t.Toxcore, (*C.uint32_t)(&friendlist[0]))
	}

	return friendlist, nil
}

/* FriendGetPublickey returns the publickey associated to that friendNumber. */
func (t *Tox) FriendGetPublickey(friendNumber uint32) ([]byte, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}
	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)
	var toxErrFriendGetPublicKey C.TOX_ERR_FRIEND_GET_PUBLIC_KEY = C.TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK
	C.tox_friend_get_public_key(t.Toxcore, (C.uint32_t)(friendNumber), (*C.uint8_t)(&publickey[0]), &toxErrFriendGetPublicKey)

	switch ToxErrFriendGetPublicKey(toxErrFriendGetPublicKey) {
	case TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK:
		return publickey, nil
	case TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND:
		return nil, ErrArgs
	default:
		return nil, ErrFuncFail
	}

	return nil, ErrUnknown
}

/* FriendGetLastOnline returns the timestamp of the last time the friend with
 * the given friendNumber was seen online. */
func (t *Tox) FriendGetLastOnline(friendNumber uint32) (time.Time, error) {
	if t.Toxcore == nil {
		return time.Time{}, ErrToxInit
	}

	var toxErrFriendGetLastOnline C.TOX_ERR_FRIEND_GET_LAST_ONLINE = C.TOX_ERR_FRIEND_GET_LAST_ONLINE_OK
	ret := C.tox_friend_get_last_online(t.Toxcore, (C.uint32_t)(friendNumber), &toxErrFriendGetLastOnline)

	if ret == C.INT64_MAX || ToxErrFriendGetLastOnline(toxErrFriendGetLastOnline) != TOX_ERR_FRIEND_GET_LAST_ONLINE_OK {
		return time.Time{}, ErrFuncFail
	}

	last := time.Unix(int64(ret), 0)

	return last, nil
}

/* FriendGetNameSize returns the length of the name of friendNumber. */
func (t *Tox) FriendGetNameSize(friendNumber uint32) (int64, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	ret := C.tox_friend_get_name_size(t.Toxcore, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return 0, ErrFuncFail
	}

	return int64(ret), nil
}

/* FriendGetName returns the name of friendNumber. */
func (t *Tox) FriendGetName(friendNumber uint32) (string, error) {
	if t.Toxcore == nil {
		return "", ErrToxInit
	}

	length, err := t.FriendGetNameSize(friendNumber)
	if err != nil {
		return "", ErrFuncFail
	}

	name := make([]byte, length)

	if length > 0 {
		var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
		success := C.tox_friend_get_name(t.Toxcore, (C.uint32_t)(friendNumber), (*C.uint8_t)(&name[0]), &toxErrFriendQuery)

		if success != true || ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(name), nil
}

/* FriendGetStatusMessageSize returns the size of the status of a friend with
 * the given friendNumber.
 */
func (t *Tox) FriendGetStatusMessageSize(friendNumber uint32) (int64, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	ret := C.tox_friend_get_status_message_size(t.Toxcore, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return 0, ErrFuncFail
	}

	return int64(ret), nil
}

/* FriendGetStatusMessage returns the status message of friend with the given
 * friendNumber.
 */
func (t *Tox) FriendGetStatusMessage(friendNumber uint32) (string, error) {
	if t.Toxcore == nil {
		return "", ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK

	size, error := t.FriendGetStatusMessageSize(friendNumber)
	if error != nil {
		return "", ErrFuncFail
	}

	statusMessage := make([]byte, size)

	if size > 0 {
		toxErrFriendQuery = C.TOX_ERR_FRIEND_QUERY_OK
		n := C.tox_friend_get_status_message(t.Toxcore, (C.uint32_t)(friendNumber), (*C.uint8_t)(&statusMessage[0]), &toxErrFriendQuery)

		if n != true || ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(statusMessage), nil
}

/* FriendGetStatus returns the status of friendNumber. */
func (t *Tox) FriendGetStatus(friendNumber uint32) (ToxUserStatus, error) {
	if t.Toxcore == nil {
		return TOX_USERSTATUS_NONE, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	status := C.tox_friend_get_status(t.Toxcore, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return TOX_USERSTATUS_NONE, ErrFuncFail
	}

	return ToxUserStatus(status), nil
}

/* FriendGetConnectionStatus returns true if the friend is connected. */
func (t *Tox) FriendGetConnectionStatus(friendNumber uint32) (ToxConnection, error) {
	if t.Toxcore == nil {
		return TOX_CONNECTION_NONE, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	status := C.tox_friend_get_connection_status(t.Toxcore, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return TOX_CONNECTION_NONE, ErrFuncFail
	}

	return ToxConnection(status), nil
}

/* FriendGetTyping returns true if friendNumber is typing. */
func (t *Tox) FriendGetTyping(friendNumber uint32) (bool, error) {
	if t.Toxcore == nil {
		return false, ErrToxInit
	}

	var toxErrFriendQuery C.TOX_ERR_FRIEND_QUERY = C.TOX_ERR_FRIEND_QUERY_OK
	istyping := C.tox_friend_get_typing(t.Toxcore, (C.uint32_t)(friendNumber), &toxErrFriendQuery)

	if ToxErrFriendQuery(toxErrFriendQuery) != TOX_ERR_FRIEND_QUERY_OK {
		return false, ErrFuncFail
	}

	return bool(istyping), nil
}

/* SelfSetTyping sets your typing status to a friend. */
func (t *Tox) SelfSetTyping(friendNumber uint32, typing bool) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var toxErrSetTyping C.TOX_ERR_SET_TYPING = C.TOX_ERR_SET_TYPING_OK
	success := C.tox_self_set_typing(t.Toxcore, (C.uint32_t)(friendNumber), (C._Bool)(typing), &toxErrSetTyping)

	if !bool(success) || ToxErrSetTyping(toxErrSetTyping) != TOX_ERR_SET_TYPING_OK {
		return ErrFuncFail
	}

	return nil
}

/* FriendSendMessage sends a message to a friend if he/she is online.
 * Maximum message length is MAX_MESSAGE_LENGTH.
 * messagetype is the type of the message (normal, action, ...).
 * Returns the message ID if successful, an error otherwise.
 */
func (t *Tox) FriendSendMessage(friendNumber uint32, messagetype ToxMessageType, message []byte) (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	if len(message) == 0 {
		return 0, ErrArgs
	}

	var cMessageType C.TOX_MESSAGE_TYPE
	if messagetype == TOX_MESSAGE_TYPE_ACTION {
		cMessageType = C.TOX_MESSAGE_TYPE_ACTION
	} else {
		cMessageType = C.TOX_MESSAGE_TYPE_NORMAL
	}

	cMessage := (*C.uint8_t)(&[]byte(message)[0])

	var toxFriendSendMessageError C.TOX_ERR_FRIEND_SEND_MESSAGE = C.TOX_ERR_FRIEND_SEND_MESSAGE_OK
	n := C.tox_friend_send_message(t.Toxcore, (C.uint32_t)(friendNumber), cMessageType, cMessage, (C.size_t)(len(message)), &toxFriendSendMessageError)

	if ToxErrFriendSendMessage(toxFriendSendMessageError) != TOX_ERR_FRIEND_SEND_MESSAGE_OK {
		return 0, ErrFuncFail
	}

	return uint32(n), nil
}

/* Hash generates a cryptographic hash of the given data (can be used to cache
 * avatars). */
func (t *Tox) Hash(data []byte) ([]byte, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	hash := make([]byte, TOX_HASH_LENGTH)

	success := C.tox_hash((*C.uint8_t)(&hash[0]), cData, C.size_t(len(data)))
	if !bool(success) {
		return nil, ErrFuncFail
	}

	return hash, nil
}

/* FileControl sends a FileControl to a friend with the given friendNumber. */
func (t *Tox) FileControl(friendNumber uint32, fileNumber uint32, fileControl ToxFileControl) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var cFileControl C.TOX_FILE_CONTROL
	switch ToxFileControl(fileControl) {
	case TOX_FILE_CONTROL_RESUME:
		cFileControl = C.TOX_FILE_CONTROL_RESUME
	case TOX_FILE_CONTROL_PAUSE:
		cFileControl = C.TOX_FILE_CONTROL_PAUSE
	case TOX_FILE_CONTROL_CANCEL:
		cFileControl = C.TOX_FILE_CONTROL_CANCEL
	}

	var toxErrFileControl C.TOX_ERR_FILE_CONTROL
	success := C.tox_file_control(t.Toxcore, (C.uint32_t)(friendNumber), (C.uint32_t)(fileNumber), cFileControl, &toxErrFileControl)

	if !bool(success) || ToxErrFileControl(toxErrFileControl) != TOX_ERR_FILE_CONTROL_OK {
		return ErrFuncFail
	}

	return nil
}

/* FileSeek sends a file seek control command to a friend for a given file
 * transfer. */
func (t *Tox) FileSeek(friendNumber uint32, fileNumber uint32, position uint64) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var toxErrFileSeek C.TOX_ERR_FILE_SEEK
	success := C.tox_file_seek(t.Toxcore, C.uint32_t(friendNumber), C.uint32_t(fileNumber), C.uint64_t(position), &toxErrFileSeek)

	if !bool(success) || ToxErrFileSeek(toxErrFileSeek) != TOX_ERR_FILE_SEEK_OK {
		return ErrFuncFail
	}

	return nil
}

/* FileGetFileId returns the file id associated to the file transfer. */
func (t *Tox) FileGetFileId(friendNumber uint32, fileNumber uint32) ([]byte, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}

	fileId := make([]byte, TOX_FILE_ID_LENGTH)

	var toxErrFileGet C.TOX_ERR_FILE_GET
	success := C.tox_file_get_file_id(t.Toxcore, C.uint32_t(friendNumber), C.uint32_t(fileNumber), (*C.uint8_t)(&fileId[0]), &toxErrFileGet)
	if !bool(success) || ToxErrFileGet(toxErrFileGet) != TOX_ERR_FILE_GET_OK {
		return nil, ErrFuncFail
	}

	return fileId, nil
}

/* FileSend sends a file transmission request. */
func (t *Tox) FileSend(friendNumber uint32, fileKind ToxFileKind, fileLength uint64, fileID []byte, fileName string) (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var cFileKind = C.TOX_FILE_KIND_DATA
	switch ToxFileKind(fileKind) {
	case TOX_FILE_KIND_AVATAR:
		cFileKind = C.TOX_FILE_KIND_AVATAR
	case TOX_FILE_KIND_DATA:
		cFileKind = C.TOX_FILE_KIND_DATA
	}

	var cFileID *C.uint8_t

	if fileID == nil {
		cFileID = nil
	} else {
		if len(fileID) != TOX_FILE_ID_LENGTH {
			return 0, ErrFileSendInvalidFileID
		}

		cFileID = (*C.uint8_t)(&[]byte(fileID)[0])
	}

	if len(fileName) == 0 {
		return 0, ErrArgs
	}

	cFileName := (*C.uint8_t)(&[]byte(fileName)[0])

	var toxErrFileSend C.TOX_ERR_FILE_SEND
	n := C.tox_file_send(t.Toxcore, (C.uint32_t)(friendNumber), (C.uint32_t)(cFileKind), (C.uint64_t)(fileLength), cFileID, cFileName, (C.size_t)(len(fileName)), &toxErrFileSend)

	if n == C.UINT32_MAX || ToxErrFileSend(toxErrFileSend) != TOX_ERR_FILE_SEND_OK {
		return 0, ErrFuncFail
	}
	return uint32(n), nil
}

/* FileSendChunk sends a chunk of file data to a friend. */
func (t *Tox) FileSendChunk(friendNumber uint32, fileNumber uint32, position uint64, data []byte) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	var toxErrFileSendChunk C.TOX_ERR_FILE_SEND_CHUNK
	success := C.tox_file_send_chunk(t.Toxcore, (C.uint32_t)(friendNumber), (C.uint32_t)(fileNumber), (C.uint64_t)(position), cData, (C.size_t)(len(data)), &toxErrFileSendChunk)

	if !bool(success) || ToxErrFileSendChunk(toxErrFileSendChunk) != TOX_ERR_FILE_SEND_CHUNK_OK {
		return ErrFuncFail
	}
	return nil
}

/* FriendSendLossyPacket sends a custom lossy packet to a friend.
 * The first byte of data must be in the range 200-254. Maximum length of a
 * custom packet is TOX_MAX_CUSTOM_PACKET_SIZE. */
func (t *Tox) FriendSendLossyPacket(friendNumber uint32, data []byte) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	var toxErrFriendCustomPacket C.TOX_ERR_FRIEND_CUSTOM_PACKET
	C.tox_friend_send_lossy_packet(t.Toxcore, C.uint32_t(friendNumber), cData, C.size_t(len(data)), &toxErrFriendCustomPacket)

	switch ToxErrFriendCustomPacket(toxErrFriendCustomPacket) {
	case TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
		return nil
	case TOX_ERR_FRIEND_CUSTOM_PACKET_NULL:
		return ErrArgs
	default:
		return ErrFuncFail
	}

	return ErrUnknown
}

/* FriendSendLosslessPacket sends a custom lossless packet to a friend.
 * The first byte of data must be in the range 160-191. Maximum length of a
 * custom packet is TOX_MAX_CUSTOM_PACKET_SIZE. */
func (t *Tox) FriendSendLosslessPacket(friendNumber uint32, data []byte) error {
	if t.Toxcore == nil {
		return ErrToxInit
	}

	var cData *C.uint8_t

	if len(data) == 0 {
		cData = nil
	} else {
		cData = (*C.uint8_t)(&data[0])
	}

	var toxErrFriendCustomPacket C.TOX_ERR_FRIEND_CUSTOM_PACKET
	C.tox_friend_send_lossless_packet(t.Toxcore, C.uint32_t(friendNumber), cData, C.size_t(len(data)), &toxErrFriendCustomPacket)

	switch ToxErrFriendCustomPacket(toxErrFriendCustomPacket) {
	case TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
		return nil
	case TOX_ERR_FRIEND_CUSTOM_PACKET_NULL:
		return ErrArgs
	default:
		return ErrFuncFail
	}

	return ErrUnknown
}

/* SelfGetDhtId returns the temporary DHT public key of this instance. */
func (t *Tox) SelfGetDhtId() ([]byte, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}

	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)

	C.tox_self_get_dht_id(t.Toxcore, (*C.uint8_t)(&publickey[0]))
	return publickey, nil
}

/* SelfGetUDPPort returns the UDP port the Tox instance is bound to. */
func (t *Tox) SelfGetUDPPort() (uint16, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var toxErrGetPort C.TOX_ERR_GET_PORT
	port := C.tox_self_get_udp_port(t.Toxcore, &toxErrGetPort)

	if ToxErrGetPort(toxErrGetPort) != TOX_ERR_GET_PORT_OK {
		return 0, ErrFuncFail
	}

	return uint16(port), nil
}

/* SelfGetTCPPort returns the TCP port the Tox instance is bound to. This is
 * only relevant if the instance is acting as a TCP relay. */
func (t *Tox) SelfGetTCPPort() (uint16, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var toxErrGetPort C.TOX_ERR_GET_PORT
	port := C.tox_self_get_tcp_port(t.Toxcore, &toxErrGetPort)

	if ToxErrGetPort(toxErrGetPort) != TOX_ERR_GET_PORT_OK {
		return 0, ErrFuncFail
	}

	return uint16(port), nil
}

// =================
// ConferenceNew creates and connects to a new text conference.
func (t *Tox) ConferenceNew() (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var toxErrConferenceNew C.Tox_Err_Conference_New
	conferenceNumber := C.tox_conference_new(t.Toxcore, &toxErrConferenceNew)

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
	if t.Toxcore == nil {
		return false, ErrToxInit
	}

	var toxErrConferenceDelete C.Tox_Err_Conference_Delete
	ret := C.tox_conference_delete(t.Toxcore, (C.uint32_t)(conferenceNumber), &toxErrConferenceDelete)
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
	if t.Toxcore == nil {
		return "", ErrToxInit
	}
	length, err := t.ConferencePeerGetNameSize(conferenceNumber, peerNumber)
	if err != nil {
		return "", ErrFuncFail
	}
	name := make([]byte, length)
	if length > 0 {
		var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query = C.TOX_ERR_CONFERENCE_PEER_QUERY_OK
		ret := C.tox_conference_peer_get_name(t.Toxcore, (C.uint32_t)(conferenceNumber), (C.uint32_t)(peerNumber), (*C.uint8_t)(&name[0]), &toxErrConferencePeerQuery)
		if ret != true || ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
			return "", ErrFuncFail
		}
	}

	return string(name), nil
}

func (t *Tox) ConferencePeerGetNameSize(conferenceNumber, peerNumber uint32) (int64, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}
	var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query = C.TOX_ERR_CONFERENCE_PEER_QUERY_OK
	ret := C.tox_conference_peer_get_name_size(t.Toxcore, (C.uint32_t)(conferenceNumber), (C.uint32_t)(peerNumber), &toxErrConferencePeerQuery)
	if ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
		return 0, ErrFuncFail
	}
	return int64(ret), nil
}

func (t *Tox) ConferencePeerGetPublicKey(conferenceNumber uint32, peerNumber uint32) (string, error) {
	if t.Toxcore == nil {
		return "", ErrToxInit
	}
	publickey := make([]byte, TOX_PUBLIC_KEY_SIZE)
	var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query
	r := C.tox_conference_peer_get_public_key(t.Toxcore, (C.uint32_t)(conferenceNumber), (C.uint32_t)(peerNumber), (*C.uint8_t)(&publickey[0]), &toxErrConferencePeerQuery)
	if bool(r) != true || ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
		return "", ErrFuncFail
	}

	pubkey := strings.ToUpper(hex.EncodeToString(publickey[:]))
	return pubkey, nil
}

func (t *Tox) ConferenceInvite(friendNumber uint32, conferenceNumber uint32) (int, error) {
	if t.Toxcore == nil {
		return -2, ErrToxInit
	}
	// if give a friendNumber which not exists,the tox_invite_friend has a strange behaive: cause other tox_* call failed
	// and the call will return true, but only strange thing accurs so just precheck the friendNumber and then go
	friendExist, err := t.FriendExists(friendNumber)
	if err != nil || friendExist == false {
		return -1, errors.New(fmt.Sprintf("friend not exists: %d", friendNumber))
	}

	var toxErrConferenceInvite C.Tox_Err_Conference_Invite
	r := C.tox_conference_invite(t.Toxcore, (C.uint32_t)(friendNumber), (C.uint32_t)(conferenceNumber), &toxErrConferenceInvite)
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
	r := C.tox_friend_exists(t.Toxcore, (C.uint32_t)(friendNumber))
	return bool(r)
}*/

func (t *Tox) ConferenceJoin(friendNumber uint32, cookie []byte) (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	if string(cookie) == "" || len(cookie) < 20 {
		return 0, errors.New("Invalid cookie:" + string(cookie))
	}

	//data, err := hex.DecodeString(cookie)
	data := []byte(cookie)
	fmt.Println(len(data))
	//fmt.Println(err)
	if data == nil { // || len(data) < 10
		return 0, errors.New("Invalid data: " + string(cookie))
	}

	var toxErrConferenceJoin C.Tox_Err_Conference_Join
	ret := C.tox_conference_join(t.Toxcore, (C.uint32_t)(friendNumber), (*C.uint8_t)(&data[0]), (C.size_t)(len(data)), &toxErrConferenceJoin)
	if ret == C.UINT32_MAX {
		return uint32(ret), errors.New(fmt.Sprintf("join group chat failed: %d", toxErrConferenceJoin))
	}

	return uint32(ret), nil
}

func (t *Tox) ConferenceSendMessage(conferenceNumber uint32, messageType ToxMessageType, message []byte) (bool, error) {
	if t.Toxcore == nil {
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
	ret := C.tox_conference_send_message(t.Toxcore, (C.uint32_t)(conferenceNumber), cMessageType, cMessage, (C.size_t)(len(message)), &toxErrConferenceSendMessage)
	if ret == false {
		return false, errors.New(fmt.Sprintf("group send message failed: %d", toxErrConferenceSendMessage))
	}
	if ToxErrConferenceSendMessage(toxErrConferenceSendMessage) != TOX_ERR_CONFERENCE_SEND_MESSAGE_OK {
		return false, ErrFuncFail
	}

	return bool(ret), nil
}

func (t *Tox) ConferenceSetTitle(conferenceNumber uint32, title string) (bool, error) {
	if t.Toxcore == nil {
		return false, ErrToxInit
	}
	var cTitle (*C.uint8_t)
	if len(title) == 0 {
		cTitle = nil
	} else {
		cTitle = (*C.uint8_t)(&[]byte(title)[0])
	}
	var toxErrConferenceTitle C.Tox_Err_Conference_Title
	success := C.tox_conference_set_title(t.Toxcore, (C.uint32_t)(conferenceNumber), cTitle, (C.size_t)(len(title)), &toxErrConferenceTitle)
	if !bool(success) || ToxErrConferenceTitle(toxErrConferenceTitle) != TOX_ERR_CONFERENCE_TITLE_OK {
		return false, ErrFuncFail
	}
	return true, nil
}

func (t *Tox) ConferenceGetTitle(conferenceNumber uint32) (string, error) {
	if t.Toxcore == nil {
		return "", ErrToxInit
	}
	length, err := t.ConferenceGetTitleSize(conferenceNumber)
	if err != nil {
		return "", ErrFuncFail
	}
	title := make([]byte, length)
	var toxErrConferenceTitle C.Tox_Err_Conference_Title
	success := C.tox_conference_get_title(t.Toxcore, (C.uint32_t)(conferenceNumber), (*C.uint8_t)(&title[0]), &toxErrConferenceTitle)
	if !bool(success) || ToxErrConferenceTitle(toxErrConferenceTitle) != TOX_ERR_CONFERENCE_TITLE_OK {
		return "", ErrFuncFail
	}

	return string(title), nil
}

func (t *Tox) ConferenceGetTitleSize(conferenceNumber uint32) (int64, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var toxErrConferenceTitle C.Tox_Err_Conference_Title
	ret := C.tox_conference_get_title_size(t.Toxcore, (C.uint32_t)(conferenceNumber), &toxErrConferenceTitle)
	if ToxErrConferenceTitle(toxErrConferenceTitle) != TOX_ERR_CONFERENCE_TITLE_OK {
		return 0, ErrFuncFail
	}
	return int64(ret), nil
}

func (t *Tox) ConferencePeerNumberIsOurs(conferenceNumber, peerNumber uint32) (bool, error) {
	if t.Toxcore == nil {
		return false, ErrToxInit
	}
	var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query
	ret := C.tox_conference_peer_number_is_ours(t.Toxcore, (C.uint32_t)(conferenceNumber), (C.uint32_t)(conferenceNumber), &toxErrConferencePeerQuery)
	if !ret || ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
		return false, ErrFuncFail
	}
	return bool(ret), nil
}

func (t *Tox) ConferencePeerCount(conferenceNumber uint32) (uint32, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var toxErrConferencePeerQuery C.Tox_Err_Conference_Peer_Query
	ret := C.tox_conference_peer_count(t.Toxcore, (C.uint32_t)(conferenceNumber), &toxErrConferencePeerQuery)
	if ToxErrConferencePeerQuery(toxErrConferencePeerQuery) != TOX_ERR_CONFERENCE_PEER_QUERY_OK {
		return 0, ErrFuncFail
	}
	return uint32(ret), nil
}

// extra combined api
func (t *Tox) ConferenceGetNames(conferenceNumber uint32) ([]string, error) {
	if t.Toxcore == nil {
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
	if t.Toxcore == nil {
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
	if t.Toxcore == nil {
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
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}
	ret := C.tox_conference_get_chatlist_size(t.Toxcore)
	return uint32(ret), nil
}

func (t *Tox) ConferenceGetChatlist() ([]uint32, error) {
	if t.Toxcore == nil {
		return nil, ErrToxInit
	}

	size, err := t.ConferenceGetChatlistSize()
	if err != nil {
		return nil, err
	}

	chatList := make([]uint32, size)

	if size > 0 {
		C.tox_conference_get_chatlist(t.Toxcore, (*C.uint32_t)(&chatList[0]))
	}

	return chatList, nil
}

func (t *Tox) ConferenceGetType(conferenceNumber uint32) (int, error) {
	if t.Toxcore == nil {
		return 0, ErrToxInit
	}

	var toxErrConferenceGetType C.Tox_Err_Conference_Get_Type
	ret := C.tox_conference_get_type(t.Toxcore, (C.uint32_t)(conferenceNumber), &toxErrConferenceGetType)
	if ToxErrConferenceGetType(toxErrConferenceGetType) != TOX_ERR_CONFERENCE_GET_TYPE_OK {
		return int(ret), ErrFuncFail
	}

	return int(ret), nil
}

func (t *Tox) ConferenceGetIdentifier(conferenceNumber uint32) (string, error) {
	if t.Toxcore == nil {
		return "", ErrToxInit
	}

	idbuf := [1 + C.TOX_PUBLIC_KEY_SIZE]byte{}

	C.tox_conference_get_id(t.Toxcore, (C.uint32_t)(conferenceNumber), (*C.uint8_t)(&idbuf[0]))

	identifier := strings.ToUpper(hex.EncodeToString(idbuf[:]))
	identifier = identifier[2:] // 1B(type)+32B(identifier)

	return identifier, nil
}
