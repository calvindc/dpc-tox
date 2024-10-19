package libtoxav

//#include <tox/toxav.h>
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
	TOXAV_ERR_NEW_OK       ToxavErrNew = C.TOXAV_ERR_NEW_OK
	TOXAV_ERR_NEW_NULL     ToxavErrNew = C.TOXAV_ERR_NEW_NULL
	TOXAV_ERR_NEW_MALLOC   ToxavErrNew = C.TOXAV_ERR_NEW_MALLOC
	TOXAV_ERR_NEW_MULTIPLE ToxavErrNew = C.TOXAV_ERR_NEW_MULTIPLE
)

type ToxavErrCall C.TOXAV_ERR_CALL

var (
	TOXAV_ERR_CALL_OK                     ToxavErrCall = C.TOXAV_ERR_CALL_OK
	TOXAV_ERR_CALL_MALLOC                 ToxavErrCall = C.TOXAV_ERR_CALL_MALLOC
	TOXAV_ERR_CALL_SYNC                   ToxavErrCall = C.TOXAV_ERR_CALL_SYNC
	TOXAV_ERR_CALL_FRIEND_NOT_FOUND       ToxavErrCall = C.TOXAV_ERR_CALL_FRIEND_NOT_FOUND
	TOXAV_ERR_CALL_FRIEND_NOT_CONNECTED   ToxavErrCall = C.TOXAV_ERR_CALL_FRIEND_NOT_CONNECTED
	TOXAV_ERR_CALL_FRIEND_ALREADY_IN_CALL ToxavErrCall = C.TOXAV_ERR_CALL_FRIEND_ALREADY_IN_CALL
	TOXAV_ERR_CALL_INVALID_BIT_RATE       ToxavErrCall = C.TOXAV_ERR_CALL_INVALID_BIT_RATE
)

var (
	ErrAnswerSync               = errors.New("TOXAV_ERR_ANSWER_SYNC")
	ErrAnswerCodecInitialzation = errors.New("TOXAV_ERR_ANSWER_CODEC_INITIALIZATION")
	ErrAnswerFriendNotFound     = errors.New("TOXAV_ERR_ANSWER_FRIEND_NOT_FOUND")
	ErrAnswerFriendNotCalling   = errors.New("TOXAV_ERR_ANSWER_FRIEND_NOT_CALLING")
	ErrAnswerInvalidBitRate     = errors.New("TOXAV_ERR_ANSWER_INVALID_BIT_RATE")
)

type ToxavErrAnswer C.TOXAV_ERR_ANSWER

var (
	TOXAV_ERR_ANSWER_OK                   ToxavErrAnswer = C.TOXAV_ERR_ANSWER_OK
	TOXAV_ERR_ANSWER_SYNC                 ToxavErrAnswer = C.TOXAV_ERR_ANSWER_SYNC
	TOXAV_ERR_ANSWER_CODEC_INITIALIZATION ToxavErrAnswer = C.TOXAV_ERR_ANSWER_CODEC_INITIALIZATION
	TOXAV_ERR_ANSWER_FRIEND_NOT_FOUND     ToxavErrAnswer = C.TOXAV_ERR_ANSWER_FRIEND_NOT_FOUND
	TOXAV_ERR_ANSWER_FRIEND_NOT_CALLING   ToxavErrAnswer = C.TOXAV_ERR_ANSWER_FRIEND_NOT_CALLING
	TOXAV_ERR_ANSWER_INVALID_BIT_RATE     ToxavErrAnswer = C.TOXAV_ERR_ANSWER_INVALID_BIT_RATE
)

type ToxavErrSendFrame C.TOXAV_ERR_SEND_FRAME

var (
	TOXAV_ERR_SEND_FRAME_OK                    ToxavErrSendFrame = C.TOXAV_ERR_SEND_FRAME_OK
	TOXAV_ERR_SEND_FRAME_NULL                  ToxavErrSendFrame = C.TOXAV_ERR_SEND_FRAME_NULL
	TOXAV_ERR_SEND_FRAME_FRIEND_NOT_FOUND      ToxavErrSendFrame = C.TOXAV_ERR_SEND_FRAME_FRIEND_NOT_FOUND
	TOXAV_ERR_SEND_FRAME_FRIEND_NOT_IN_CALL    ToxavErrSendFrame = C.TOXAV_ERR_SEND_FRAME_FRIEND_NOT_IN_CALL
	TOXAV_ERR_SEND_FRAME_SYNC                  ToxavErrSendFrame = C.TOXAV_ERR_SEND_FRAME_SYNC
	TOXAV_ERR_SEND_FRAME_INVALID               ToxavErrSendFrame = C.TOXAV_ERR_SEND_FRAME_INVALID
	TOXAV_ERR_SEND_FRAME_PAYLOAD_TYPE_DISABLED ToxavErrSendFrame = C.TOXAV_ERR_SEND_FRAME_PAYLOAD_TYPE_DISABLED
	TOXAV_ERR_SEND_FRAME_RTP_FAILED            ToxavErrSendFrame = C.TOXAV_ERR_SEND_FRAME_RTP_FAILED
)
