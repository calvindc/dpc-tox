#include <tox/toxav.h>

/* Convenient macro:
 * Creates the C function to directly register a given callback */
#define CREATE_HOOK(x) \
static void set_##x(ToxAV *toxav, void *user_data) { \
toxav_##x(toxav, hook_##x, user_data); \
}

// Headers for the exported GO functions in hookav.go
void hook_callback_call(ToxAV*, uint32_t, bool, bool, void*);



CREATE_HOOK(callback_call)