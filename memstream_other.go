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

// +build !linux,!darwin,!freebsd

package pbc

/*
#include <stdlib.h>
#include "memstream.h"

struct memstream_s {
	FILE* fd;
};

memstream_t* pbc_open_memstream() {
	FILE* fd = tmpfile();
	if (fd == NULL) return NULL;
	memstream_t* result = malloc(sizeof(memstream_t));
	result->fd = fd;
	return result;
}

FILE* pbc_memstream_to_fd(memstream_t* m) { return m->fd; }

int pbc_close_memstream(memstream_t* m, char** bufp, size_t* sizep) {
	*bufp = NULL;
	*sizep = 0;

	FILE* fd = m->fd;
	m->fd = NULL;
	free(m);
	m = NULL;

	if (!ferror(fd)) {
		fseek(fd, 0, SEEK_END);
		*sizep = (size_t)ftell(fd);
		rewind(fd);
		*bufp = malloc(*sizep + 1);
		size_t readBytes = fread(*bufp, 1, *sizep, fd);
		if (readBytes < *sizep || ferror(fd)) {
			free(*bufp);
			*bufp = NULL;
			*sizep = 0;
		}
		bufp[*sizep] = '\0';
	}
	fclose(fd);
	return (*bufp != NULL);
}
*/
import "C"
