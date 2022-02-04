package rw

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestLimitedWriter(t *testing.T) {
	t.Run("max", testLimitedWriterMax)
	t.Run("short", testLimitedWriterShort)
	t.Run("error", testLimitedWriterError)
	t.Run("chunk", testLimitedWriterChunk)
}

func testLimitedWriterError(t *testing.T) {
	var (
		err error
		wrt = limitWriter(32)
	)
	_, err = io.WriteString(wrt, strings.Repeat("0", 64))
	if err == nil {
		t.Fatalf("expected error writing more bytes than allow! but got none")
	}
}

func testLimitedWriterChunk(t *testing.T) {
  var (
    err error
    wrt = limitWriter(32)
  )
  sizes := []struct{
    Size int
    Err  bool
  } {
    {Size: 16, Err: false},
    {Size: 8, Err: false},
    {Size: 16, Err: true},
  }
  for _, s := range sizes {
    _, err = io.WriteString(wrt, strings.Repeat("0", s.Size))
    switch {
    case s.Err && err == nil:
      t.Fatalf("expected error but got none writing %d bytes", s.Size)
    case !s.Err && err != nil:
      t.Fatalf("expected no error but got one writing %d bytes: %s", s.Size, err)
    }
  }

}

func testLimitedWriterMax(t *testing.T) {
	var (
		err error
		wrt = limitWriter(32)
	)
	_, err = io.WriteString(wrt, strings.Repeat("0", 32))
	if err != nil {
		t.Fatalf("unexpected error writing max bytes in writer")
	}
}

func testLimitedWriterShort(t *testing.T) {
	var (
		err error
		wrt = limitWriter(32)
	)
	_, err = io.WriteString(wrt, strings.Repeat("0", 16))
	if err != nil {
		t.Fatalf("unexpected error writing less bytes in writer")
	}
}

func limitWriter(max int64) io.Writer {
	var (
		buf bytes.Buffer
		wrt = LimitWriter(&buf, max)
	)
	return wrt
}
