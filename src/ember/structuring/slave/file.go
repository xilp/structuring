package slave

import (
	"bytes"
	"io"
	"os"
	"sync"
)

func (p *RawFile) write(str string) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.buf.WriteString(str)
	p.cache+= 1
	if p.cache >= p.cacheLine {
		p.flush()
	}
	return err
}

func (p *RawFile) flush() (err error) {
	io.Copy(p.file, p.buf)
	p.buf.Reset()
	p.cache= 0
	return err
}

func NewRawFile() (file RawFile, err error){
	file.file, err = os.OpenFile("binlog.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		return
	}
	file.buf = bytes.NewBuffer([]byte{})
	file.cache = 0
	file.cacheLine = 100
	return file, err
}

type RawFile struct {
	file *os.File
	locker sync.Mutex
	cache int
	cacheLine int
	buf *bytes.Buffer
}
