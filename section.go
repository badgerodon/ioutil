package ioutil

import (
	"errors"
	"io"
)

type (
	sectionReader struct {
		src              io.ReadSeeker
		base, off, limit int64
	}
)

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func NewSectionReader(src io.ReadSeeker, offset, length int64) io.ReadSeeker {
	return &sectionReader{src, offset, offset, offset + length}
}

func (this *sectionReader) Read(p []byte) (n int, err error) {
	if this.off >= this.limit {
		return 0, io.EOF
	}
	_, err = this.src.Seek(this.off, 0)
	if err != nil {
		return
	}
	if max := this.limit - this.off; int64(len(p)) > max {
		p = p[:max]
	}
	n, err = this.src.Read(p)
	this.off += int64(n)
	return
}
func (this *sectionReader) Seek(offset int64, whence int) (n int64, err error) {
	switch whence {
	default:
		return 0, errWhence
	case 0:
		offset += this.base
	case 1:
		offset += this.off
	case 2:
		offset += this.limit
	}
	if offset < this.base || offset > this.limit {
		return 0, errOffset
	}
	this.off = offset
	n, err = this.src.Seek(this.off, 0)
	n -= this.off
	return
}
