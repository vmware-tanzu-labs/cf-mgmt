package ios

import (
	"io"
	"os"
)

func CloseReaderIgnoreError(r io.ReadCloser) {
	_ = r.Close()
}

func CloseWriterIgnoreError(w io.WriteCloser) {
	_ = w.Close()
}

func CloseIgnoreError(c io.Closer) {
	_ = c.Close()
}

func CleanupTempFile(f *os.File) {
	_ = f.Close()
	_ = os.Remove(f.Name())
}
