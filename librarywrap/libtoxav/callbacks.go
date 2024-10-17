package libtoxav

import "unsafe"

type OnCall func(toxav *ToxAV, friendnumber uint32, audioenabled bool, videoenabled bool, userdata unsafe.Pointer)
