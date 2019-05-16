package gexto

import "io"

// A limitedWriter writes to W but limits the amount of
// data written to just N bytes. Each call to Write
// updates N to reflect the new amount remaining.
// Write returns EOF when N <= 0 or when the underlying W returns EOF.
type limitedWriter struct {
	W io.Writer
	N int64
}

func LimitWriter(w io.Writer, n int64) io.Writer { return &limitedWriter{w, n} }

func (lw *limitedWriter) Write(p []byte) (n int, err error) {
	if lw.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > lw.N {
		p = p[0:lw.N]
	}
	n, err = lw.W.Write(p)
	lw.N -= int64(n)
	return
}