#include <tox/toxav.h>

/* Macro defined:
 * Creates the C function to directly register a given callback from toxav.h
 */
#define CREATE_HOOK(x) \
static void set_##x(ToxAV *toxav, void *t, void *user_data) { \
toxav_##x(toxav, hook_##x, user_data); \
}

//Tag: Headers for the exported GO functions in /libtoxav/hooks.go

//typedef void toxav_call_cb(ToxAV *av, uint32_t friend_number, bool audio_enabled, bool video_enabled, void *user_data);
void hook_callback_call(ToxAV*, uint32_t, bool, bool, void*);

//typedef void toxav_audio_receive_frame_cb(ToxAV *av, uint32_t friend_number, const int16_t pcm[],
//size_t sample_count,uint8_t channels, uint32_t sampling_rate, void *user_data);
void hook_callback_audio_receive_frame(ToxAV*, uint32_t, const int16_t*, size_t, uint8_t, uint32_t, void*);

//typedef void toxav_call_state_cb(ToxAV *av, uint32_t friend_number, uint32_t state, void *user_data);
void hook_callback_call_state(ToxAV*, uint32_t, uint32_t, void*);

//toxav callback functions
CREATE_HOOK(callback_call)
CREATE_HOOK(callback_audio_receive_frame)
CREATE_HOOK(callback_call_state)