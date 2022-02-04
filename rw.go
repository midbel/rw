package rw

import (
  "io"
)

type LimitedWriter struct{
  W io.Writer
  N int64
  r int64
}

func LimitWriter(w io.Writer, n int64) io.Writer {
  return &LimitedWriter{
    W: w,
    N: n,
    r: n,
  }
}

func (w *LimitedWriter) Write(b []byte) (int, error) {
  return 0, nil
}
