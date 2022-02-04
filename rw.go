package rw

import (
	"errors"
	"io"
	"os"
)

var ErrTooLong = errors.New("rw: write too long")

type LimitedWriter struct {
	W io.Writer
	N int64
	r int64
}

func LimitWriter(w io.Writer, n int64) io.Writer {
	return &LimitedWriter{
		W: w,
		N: n,
	}
}

func (w *LimitedWriter) Write(b []byte) (int, error) {
	if w.N <= 0 {
		return 0, ErrTooLong
	}
	var (
		remain int64
		err    error
	)
	if remain = w.N - int64(len(b)); remain < 0 {
		b = b[:w.N]
		err = ErrTooLong
	}
	n, err1 := w.W.Write(b)
	if err == nil {
		err = err1
	}
  w.N -= int64(n)
	return n, err
}

func (w *LimitedWriter) Available() int64 {
	return w.N
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
		p   Pipe
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

func (p *Pipe) Close() error {
  errw := p.W.Close()
  errr := p.R.Close()
  return hasErrpr(errw, errr)
}

func hasError(errs ...error) error {
  for _, e := range errs {
    if e != nil {
      return e
    }
  }
  return nil
}
