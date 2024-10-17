#include <tox/toxav.h>

/* Macro defined:
 * Creates the C function to directly register a given callback from toxav.h
 */
#define CREATE_HOOK(x) \
static void set_##x(ToxAV *toxav, void *t, void *user_data) { \
toxav_##x(toxav, hook_##x, t, user_data); \
}

//Tag: Headers for the exported GO functions in /libtoxav/hooks.go

//typedef void toxav_call_cb(ToxAV *av, uint32_t friend_number, bool audio_enabled, bool video_enabled, void *user_data);
void hook_callback_call(ToxAV*, uint32_t, bool, bool, void*);

//toxav callback functions
CREATE_HOOK(callback_call)