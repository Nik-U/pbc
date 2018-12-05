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
//

// +build darwin freebsd

package pbc

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <unistd.h>
#include <errno.h>
#include "memstream.h"

struct memstream_s {
	char*  buffer;
	size_t cap;
	size_t len;
	fpos_t cursor;
	FILE*  fd;
	int    closed;
};

static void memstream_realloc(memstream_t* m, size_t size) {
	m->cap = size;
	m->buffer = realloc(m->buffer, size);
	if (size < m->len) {
		m->len = size;
		m->cursor = m->len;
	}
}

static int memstream_read(void* cookie, char* buf, int nbytes) {
	memstream_t* m = (memstream_t*)cookie;
	if (m->closed) { errno = EBADF; return -1; }
	if (m->cursor >= m->len) return 0;
	size_t toRead = (m->len - m->cursor);
	if (toRead > (size_t)nbytes) toRead = (size_t)nbytes;
	memcpy(buf, &m->buffer[m->cursor], toRead);
	m->cursor += toRead;
	return toRead;
}

static int memstream_write(void* cookie, const char* buf, int nbytes) {
	memstream_t* m = (memstream_t*)cookie;
	if (m->closed) { errno = EBADF; return -1; }
	size_t zeros = m->cursor - m->len;
	size_t dataLen = (size_t)nbytes;
	size_t neededSpace = zeros + dataLen;
	if (neededSpace < dataLen) {
		errno = EFBIG;
		return -1;
	}
	size_t newLen = m->len + neededSpace;
	if (newLen < m->len) {
		errno = EFBIG;
		return -1;
	}
	if (newLen > m->cap) {
		size_t newCap = m->cap;
		do {
			if (SIZE_MAX - newCap < newCap) {
				newCap = SIZE_MAX;
			} else {
				newCap <<= 1;
			}
		} while (newLen > newCap);
		memstream_realloc(m, newCap);
	}
	if (zeros > 0) {
		memset(&m->buffer[m->len], 0, zeros);
	}
	memcpy(&m->buffer[m->cursor], buf, dataLen);
	m->cursor += dataLen;
	m->len = newLen;
	return neededSpace;
}

static fpos_t memstream_seek(void* cookie, fpos_t offset, int whence) {
	memstream_t* m = (memstream_t*)cookie;
	if (m->closed) { errno = EBADF; return -1; }
	fpos_t base = 0;
	switch (whence) {
		case SEEK_SET: base = 0; break;
		case SEEK_CUR: base = m->cursor; break;
		case SEEK_END: base = m->len; break;
		// SEEK_HOLE and SEEK_DATA are not supported on darwin
		default: errno = EINVAL; return -1;
	}
	fpos_t desired = base + offset;
	if (desired < 0) {
		errno = EINVAL;
		return -1;
	}
	if (offset > 0 && desired < base) {
		errno = EOVERFLOW;
		return -1;
	}
	return (m->cursor = desired);
}

static int memstream_close(void* cookie) {
	memstream_t* m = (memstream_t*)cookie;
	m->closed = 1;
	return 0;
}

memstream_t* pbc_open_memstream() {
	memstream_t* m = malloc(sizeof(memstream_t));
	m->buffer = NULL;
	memstream_realloc(m, 1024);
	m->len = 0;
	m->cursor = 0;
	m->closed = 0;
	m->fd = funopen(m, memstream_read, memstream_write, memstream_seek, memstream_close);
	return m;
}

FILE* pbc_memstream_to_fd(memstream_t* m) { return m->fd; }

int pbc_close_memstream(memstream_t* m, char** bufp, size_t* sizep) {
	fclose(m->fd);
	*bufp = m->buffer;
	*sizep = m->len;
	free(m);
	return 1;
}
*/
import "C"
