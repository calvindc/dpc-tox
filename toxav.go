package dpc_tox

import "C"

import "sync"

//#cgo LDFLAGS: -ltoxcore
//#include <tox/toxav.h>
//#include <vpx/vpx_image.h>

type ToxAV struct {
	tox   *Tox
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
	onCall OnCall
}
