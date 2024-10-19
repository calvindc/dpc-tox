package libtoxav

/*
#include <tox/toxav.h>
#include "hooks-macro.c"
*/
import "C"
import "unsafe"

type OnCall func(toxav *ToxAV, friendnumber uint32, audioenabled bool, videoenabled bool)

type OnAudioReceiveFrame func(toxav *ToxAV, friendnumber uint32, pcm []byte, sample_count int, channels int, samplingrate int)

type OnCallState func(toxav *ToxAV, friendnumber uint32, state uint32)

func (tav *ToxAV) CallbackCall(f OnCall) {
	if tav.toxav != nil {
		tav.onCall = f
		C.set_callback_call(tav.toxav, unsafe.Pointer(tav), unsafe.Pointer(tav))
	}
}
func (tav *ToxAV) CallbackAudioReceiveFrame(f OnAudioReceiveFrame) {
	if tav.toxav != nil {
		tav.onAudioReceiveFrame = f
		C.set_callback_audio_receive_frame(tav.toxav, unsafe.Pointer(tav), unsafe.Pointer(tav))
	}
}
func (tav *ToxAV) CallbackCallState(f OnCallState) {
	if tav.toxav != nil {
		tav.onCallState = f
		C.set_callback_call_state(tav.toxav, unsafe.Pointer(tav), unsafe.Pointer(tav))
	}
}
