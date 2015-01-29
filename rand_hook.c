#include "_cgo_export.h"
#include <pbc/pbc.h>

void pbc_init_random();

void goRandomHook(mpz_t out, mpz_t limit, void* data) {
	UNUSED_VAR(data);
	goGenerateRandom(out, limit);
}

void installRandomHook() {
	pbc_random_set_function(goRandomHook, NULL);
}

void uninstallRandomHook() {
	pbc_init_random();
}