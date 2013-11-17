package ioutil

import (
	"errors"
	"io"
)

type (
	multiReadSeeker struct {
		components   []multiReadSeekerComponent
		idx          int
		offset, size int64
		initialized  bool
	}
	multiReadSeekerComponent struct {
		io.ReadSeeker
		offset, size int64
	}
)

func (this *multiReadSeeker) Read(p []byte) (int, error) {
	// initialize
	err := this.init()
	if err != nil {
		return 0, err
	}

	total := 0

	for this.idx < len(this.components) {
		n, err := this.components[this.idx].Read(p)
		total += n
		// Ignore EOF errors
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			return total, err
		}
		this.offset += int64(n)
		this.components[this.idx].offset += int64(n)
		// Move to the next one if necessary
		if this.components[this.idx].offset == this.components[this.idx].size {
			this.idx++
		}
		p = p[n:]
		if len(p) == 0 {
			return total, nil
		}
	}

	if total > 0 {
		return total, nil
	}
	return 0, io.EOF
}

func (this *multiReadSeeker) Seek(offset int64, whence int) (ret int64, err error) {
	// initialize
	err = this.init()
	if err != nil {
		return
	}

	switch whence {
	case 0:
	case 1:
		offset += this.offset
	case 2:
		offset += this.size
	default:
		err = errors.New("Seek: invalid whence")
		return
	}

	if offset > this.size || offset < 0 {
		err = errors.New("Seek: invalid offset")
		return
	}

	for offset > this.offset {
		rel := offset - this.offset
		if rel > (this.components[this.idx].size - this.components[this.idx].offset) {
			rel = this.components[this.idx].size - this.components[this.idx].offset
		}
		var n int64
		n, err = this.components[this.idx].Seek(rel, 1)
		if err != nil {
			return
		}
		rel = n - this.components[this.idx].offset
		this.offset += rel
		this.components[this.idx].offset = n
		if offset > this.offset {
			this.idx++
		}
	}

	for offset < this.offset {
		if this.idx < len(this.components) {
			rel := this.components[this.idx].offset - (this.offset - offset)
			if rel < 0 {
				rel = 0
			}
			_, err = this.components[this.idx].Seek(rel, 0)
			if err != nil {
				return
			}
			this.offset -= this.components[this.idx].offset - rel
			this.components[this.idx].offset = rel
		}
		if offset < this.offset {
			this.idx--
		}
	}

	ret = this.offset

	return
}

func (this *multiReadSeeker) init() error {
	if !this.initialized {
		for i, component := range this.components {
			size, err := component.Seek(0, 2)
			if err != nil {
				return err
			}
			_, err = component.Seek(0, 0)
			if err != nil {
				return err
			}
			this.components[i].size = size
			this.size += size
		}
		this.initialized = true
	}
	return nil
}

func NewMultiReadSeeker(readSeekers ...io.ReadSeeker) io.ReadSeeker {
	components := make([]multiReadSeekerComponent, 0, len(readSeekers))
	for _, rdr := range readSeekers {
		components = append(components, multiReadSeekerComponent{rdr, 0, 0})
	}
	return &multiReadSeeker{
		components: components,
	}
}
