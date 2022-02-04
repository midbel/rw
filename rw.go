package rw

import (
  "errors"
  "io"
  "os"
)

var ErrTooLong = errors.New("rw: write too long")

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
  if w.r <= 0 {
    return 0, ErrTooLong
  }

  return 0, nil
}

var zeroes = make([]byte, 4096)

type zero struct{}

func Zero() io.Reader {
  return zero{}
}

func (_ zero) Read(b []byte) (int, error) {
  var n int
  for n < len(b) {
    x := copy(b[n:], zeroes)
    n += x
  }
  return n, nil
}

type Pipe struct {
  R *os.File
  W *os.File
}

func PipeFile() (*Pipe, error) {
  var (
    p Pipe
    err error
  )
  p.R, p.W, err = os.Pipe()
  return &p, err
}

func (p *Pipe) Read(b []byte) (int, error) {
  return p.R.Read(b)
}

func (p *Pipe) Write(b []byte) (int, error) {
  return p.W.Write(b)
}
