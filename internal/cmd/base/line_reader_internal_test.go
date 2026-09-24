// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

package base

import (
	"io"
	"strings"
	"testing"
)

func TestLineReaderReturnsOneLinePerRead(t *testing.T) {
	lr := lineReader{r: strings.NewReader("first\nsecond\nlast")}
	buf := make([]byte, 64)
	for _, want := range []string{"first\n", "second\n", "last"} {
		n, err := lr.Read(buf)
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
		if got := string(buf[:n]); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
	if _, err := lr.Read(buf); err != io.EOF {
		t.Errorf("err = %v, want io.EOF", err)
	}
}

func TestLineReaderSmallBuffer(t *testing.T) {
	lr := lineReader{r: strings.NewReader("abcdef\n")}
	buf := make([]byte, 4)
	n, _ := lr.Read(buf)
	if string(buf[:n]) != "abcd" {
		t.Errorf("first read = %q, want abcd", buf[:n])
	}
	n, _ = lr.Read(buf)
	if string(buf[:n]) != "ef\n" {
		t.Errorf("second read = %q, want ef\\n", buf[:n])
	}
}
