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

#include <stdio.h>

// memstream_s is a structure that provides platform-independent conversion from
// file descriptor writes to strings
typedef struct memstream_s memstream_t;

// pbc_open_memstream returns a memstream that can be used for writing data
memstream_t* pbc_open_memstream();

// pbc_memstream_to_fd retrieves the file descriptor for a memstream
FILE* pbc_memstream_to_fd(memstream_t* m);

// pbc_close_memstream closes the memstream and returns the written data
int pbc_close_memstream(memstream_t* m, char** bufp, size_t* sizep);
