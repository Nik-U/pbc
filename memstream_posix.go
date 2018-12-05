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

// +build linux

package pbc

/*
#include <stdlib.h>
#include <stdio.h>
#include "memstream.h"

struct memstream_s {
	char*  buf;
	size_t size;
	FILE*  fd;
};

memstream_t* pbc_open_memstream() {
	memstream_t* result = malloc(sizeof(memstream_t));
	result->fd = open_memstream(&result->buf, &result->size);
	if (result->fd == NULL) {
		free(result);
		result = NULL;
	}
	return result;
}

FILE* pbc_memstream_to_fd(memstream_t* m) { return m->fd; }

int pbc_close_memstream(memstream_t* m, char** bufp, size_t* sizep) {
	fclose(m->fd);
	*bufp = m->buf;
	*sizep = m->size;
	free(m);
	return 1;
}
*/
import "C"
