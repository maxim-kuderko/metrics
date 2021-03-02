package drivers

import (
	"os"
)

func NewStdout() *Writer {
	return NewWriter(os.Stdout)
}
