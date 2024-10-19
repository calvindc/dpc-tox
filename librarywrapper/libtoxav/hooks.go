package libtoxav

//#include <tox/toxav.h>
import "C"
import (
	"unsafe"
)

//export hook_callback_call
func hook_callback_call(t unsafe.Pointer, friendnumber C.uint32_t, audioenabled C._Bool, videoenabled C._Bool, toxav unsafe.Pointer) {
	(*ToxAV)(toxav).onCall((*ToxAV)(toxav), uint32(friendnumber), bool(audioenabled), bool(videoenabled))
}

//export hook_callback_audio_receive_frame
func hook_callback_audio_receive_frame(t unsafe.Pointer, friendnumber C.uint32_t, pcm *C.uint8_t, samplecount C.size_t, channels C.uint8_t, samplingrate C.uint32_t, toxav unsafe.Pointer) {
	(*ToxAV)(toxav).onAudioReceiveFrame((*ToxAV)(toxav), uint32(friendnumber), C.GoBytes(unsafe.Pointer(pcm), C.int(samplecount*C.size_t(channels)*2)), int(samplecount), int(channels), int(samplingrate))
}

//export hook_callback_call_state
func hook_callback_call_state(t unsafe.Pointer, friendnumber C.uint32_t, state C.uint32_t, toxav unsafe.Pointer) {
	(*ToxAV)(toxav).onCallState((*ToxAV)(toxav), uint32(friendnumber), uint32(state))
}
