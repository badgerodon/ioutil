package ioutil

import (
	"io"
)

type multiReadSeeker struct {
	readSeekers []io.ReadSeeker
	current int
}

func (this *multiReadSeeker) Read(p []byte) (n int, err error) {
	for this.current < len(this.readSeekers) {
		n, err = this.readSeekers[this.current].Read(p)
		if n > 0 || err != io.EOF {
			if err == io.EOF {
				err = nil
			}
			return
		}
		this.current++
	}
	err = io.EOF
	return
}

func (this *multiReadSeeker) Seek(offset int64, whence int) (ret int64, err error) {
	return
}

func NewMultiReadSeeker(readSeekers ...io.ReadSeeker) io.ReadSeeker {
	return &multiReadSeeker{readSeekers, 0}
}
