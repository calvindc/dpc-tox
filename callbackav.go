package dpc_tox

/*
#include <tox/toxav.h>
#include "hookav.c"
*/
import "C"
import "unsafe"

type OnCall func(toxav *ToxAV, friendnumber uint32, audioenabled bool, videoenabled bool, userdata unsafe.Pointer)

func (tav *ToxAV) CallbackCall(f OnCall, userData unsafe.Pointer) {
	if tav.toxav != nil {
		tav.onCall = f
		tav.onCallUserData = userData
		C.set_callback_call(tav.toxav, userData)
	}
}
