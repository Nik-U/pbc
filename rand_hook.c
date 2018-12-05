// Copyright © 2018 Nik Unger
//
// This file is part of The PBC Go Wrapper.
//
// The PBC Go Wrapper is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// The PBC Go Wrapper is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public
// License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with The PBC Go Wrapper. If not, see <http://www.gnu.org/licenses/>.
//
// The PBC Go Wrapper makes use of The PBC library. The PBC Library and its use
// are covered under the terms of the GNU Lesser General Public License
// version 3, or (at your option) any later version.

#include "_cgo_export.h"
#include <pbc/pbc.h>

void pbc_init_random();

void goRandomHook(mpz_t out, mpz_t limit, void* data) {
	UNUSED_VAR(data);
	goGenerateRandom(&out, &limit);
}

void installRandomHook() {
	pbc_random_set_function(goRandomHook, NULL);
}

void uninstallRandomHook() {
	pbc_init_random();
}
