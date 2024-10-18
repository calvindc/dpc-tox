package libtoxav

import "C"
import "errors"

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

// ==== toxAV error ====
var (
	ErrToxAVNew  = errors.New("Error initializing ToxAV")
	ErrToxAVInit = errors.New("ToxAV not initialized")
)
var (
	ErrNewMultiple = errors.New("Not allow to create a second session")
)

type ToxavErrNew C.TOXAV_ERR_NEW

var (
	TOXAV_ERR_NEW_OK       = C.TOXAV_ERR_NEW_OK
	TOXAV_ERR_NEW_NULL     = C.TOXAV_ERR_NEW_NULL
	TOXAV_ERR_NEW_MALLOC   = C.TOXAV_ERR_NEW_MALLOC
	TOXAV_ERR_NEW_MULTIPLE = C.TOXAV_ERR_NEW_MULTIPLE
)

type ToxavErrCall C.TOXAV_ERR_CALL

var (
	TOXAV_ERR_CALL_OK                     = C.TOXAV_ERR_CALL_OK
	TOXAV_ERR_CALL_MALLOC                 = C.TOXAV_ERR_CALL_MALLOC
	TOXAV_ERR_CALL_SYNC                   = C.TOXAV_ERR_CALL_SYNC
	TOXAV_ERR_CALL_FRIEND_NOT_FOUND       = C.TOXAV_ERR_CALL_FRIEND_NOT_FOUND
	TOXAV_ERR_CALL_FRIEND_NOT_CONNECTED   = C.TOXAV_ERR_CALL_FRIEND_NOT_CONNECTED
	TOXAV_ERR_CALL_FRIEND_ALREADY_IN_CALL = C.TOXAV_ERR_CALL_FRIEND_ALREADY_IN_CALL
	TOXAV_ERR_CALL_INVALID_BIT_RATE       = C.TOXAV_ERR_CALL_INVALID_BIT_RATE
)

type ToxavErrAnswer C.TOXAV_ERR_ANSWER

var (
	TOXAV_ERR_ANSWER_OK                   = C.TOXAV_ERR_ANSWER_OK
	TOXAV_ERR_ANSWER_SYNC                 = C.TOXAV_ERR_ANSWER_SYNC
	TOXAV_ERR_ANSWER_CODEC_INITIALIZATION = C.TOXAV_ERR_ANSWER_CODEC_INITIALIZATION
	TOXAV_ERR_ANSWER_FRIEND_NOT_FOUND     = C.TOXAV_ERR_ANSWER_FRIEND_NOT_FOUND
	TOXAV_ERR_ANSWER_FRIEND_NOT_CALLING   = C.TOXAV_ERR_ANSWER_FRIEND_NOT_CALLING
	TOXAV_ERR_ANSWER_INVALID_BIT_RATE     = C.TOXAV_ERR_ANSWER_INVALID_BIT_RATE
)
