package libtoxav

//#include <tox/toxav.h>
import "C"
import "unsafe"

//export hook_callback_call
func hook_callback_call(t unsafe.Pointer, friendnumber C.uint32_t, audioenabled C._Bool, videoenabled C._Bool, toxav unsafe.Pointer, userdata unsafe.Pointer) {
	(*ToxAV)(toxav).onCall((*ToxAV)(toxav), uint32(friendnumber), bool(audioenabled), bool(videoenabled), userdata)
}
