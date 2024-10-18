package libtoxav

//#cgo LDFLAGS: -ltoxcore
//#include <tox/toxav.h>
//#include <vpx/vpx_image.h>
import "C"

import (
	"errors"
	"github.com/calvindc/dpc-tox/librarywrapper/libtox"
	"sync"
	"unsafe"
)

type ToxAV struct {
	tox   libtox.Tox
	toxav *C.ToxAV
	mtx   sync.Mutex

	// session datas
	out_image  []byte
	out_width  uint16
	out_hegith uint16
	in_image   *C.vpx_image_t
	in_width   uint16
	in_height  uint16

	// Callbacks
	onCall         OnCall
	onCallUserData unsafe.Pointer
}

/**
 * Start new A/V session. There can only be only one session per Tox instance.
 * ToxAV *toxav_new(Tox *tox, Toxav_Err_New *error);
 */
func NewToxAV(tox *libtox.Tox) (*ToxAV, error) {
	var cToxAV *C.ToxAV
	var toxAVErrNew C.TOXAV_ERR_NEW
	cToxAV = C.toxav_new(tox.Toxcore, &toxAVErrNew)
	if cToxAV == nil || ToxavErrNew(toxAVErrNew) != TOXAV_ERR_NEW_OK {
		switch ToxavErrNew(toxAVErrNew) {
		case TOXAV_ERR_NEW_NULL:
			return nil, ErrArgs
		case TOXAV_ERR_NEW_MALLOC:
			return nil, ErrNewMalloc
		case TOXAV_ERR_NEW_MULTIPLE:
			return nil, ErrNewMultiple
		}
		if cToxAV == nil {
			return nil, ErrToxAVNew
		}
		return nil, ErrUnknown
	}
	tav := &ToxAV{tox: tox.Toxcore, toxav: cToxAV}

	return tav, nil
}

/**
 * Releases all resources associated with the A/V session.
 * If any calls were ongoing, these will be forcibly terminated without
 * notifying peers. After calling this function, no other functions may be
 * called and the av pointer becomes invalid.
 */
func (tav *ToxAV) Kill() {
	C.toxav_kill(tav.toxav)
}

/**
 * Returns the Tox instance the A/V object was created for.
 */
func (tav *ToxAV) GetTox() *libtox.Tox {
	return tav.toxav
}

/**
 * A/V event loop, single thread
 * Returns the interval in milliseconds when the next toxav_iterate call should
 * be. If no call is active at the moment, this function returns 200.
 * This function MUST be called from the same thread as toxav_iterate.
 */
func (tav *ToxAV) IterationInterval() uint32 {
	return uint32(C.toxav_iteration_interval(tav.toxav))
}

/**
 * Main loop for the session. This function needs to be called in intervals of
 * `toxav_iteration_interval()` milliseconds. It is best called in the separate
 * thread from tox_iterate.
 */
func (tav *ToxAV) Iterate() {
	C.toxav_iterate(tav.toxav)
}

/**
 * A/V event loop, multiple threads
 * Returns the interval in milliseconds when the next toxav_audio_iterate call
 * should be. If no call is active at the moment, this function returns 200.
 * This function MUST be called from the same thread as toxav_audio_iterate.
 */
func (tav *ToxAV) AudioIterationInterval() uint32 {
	return uint32(C.toxav_audio_iteration_interval(tav.toxav))
}

func (tav *ToxAV) AudioIterate() {
	C.toxav_audio_iterate(tav.toxav)
}

func (tav *ToxAV) Call(friendNumber uint32, audioBitRate uint32, videoBitRate uint32) (bool, error) {
	var toxavErrCall C.TOXAV_ERR_CALL
	ret := C.toxav_call(tav.toxav, (C.uint32_t)(friendNumber), (C.uint32_t)(audioBitRate), (C.uint32_t)(videoBitRate), &toxavErrCall)
	if ToxavErrCall(toxavErrCall) != TOXAV_ERR_CALL_OK {
		return bool(ret), errors.New(string(toxavErrCall)) //todo: write to error discription
	}
	return bool(ret), nil
}

func (tav *ToxAV) Answer(friendNumber uint32, audioBitRate uint32, videoBitRate uint32) (bool, error) {
	var toxavErrAnswer C.TOXAV_ERR_ANSWER
	ret := C.toxav_answer(tav.toxav, C.uint32_t(friendNumber), C.uint32_t(audioBitRate), C.uint32_t(videoBitRate), &toxavErrAnswer)
	if ToxavErrAnswer(toxavErrAnswer) != TOXAV_ERR_ANSWER_OK {
		return bool(ret), errors.New(string(toxavErrAnswer))
	}
	return bool(ret), nil
}
