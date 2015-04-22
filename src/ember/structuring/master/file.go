package master

import (
	"bytes"
	"io"
	"os"
	"sync"
	"bufio"
)

func (p *RawFile) Write(buf []byte, str string) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.buf.WriteString(str)
	p.cache+= 1
	if p.cache >= p.cacheLine {
		p.Flush()
	}
	return err
}

func (p *RawFile) Read() (ret string, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	r := bufio.NewReader(p.file)
	for {
		line, err := r.ReadString('\n')
		if io.EOF == err || nil != err {
			break
		}
		ret = ret + line
	}
	return ret, err
}

func (p *RawFile) Close() (err error) {
	p.file.Close()
	return err
}

func (p *RawFile) Flush() (err error) {
	// TODO lock 
	p.file.Truncate(0)
	io.Copy(p.file, p.buf)
	p.buf.Reset()
	p.cache= 0
	return err
}

func NewRawFile(fileName string) (file RawFile, err error){
	file.file, err = os.OpenFile(fileName, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		return
	}
	file.buf = bytes.NewBuffer([]byte{})
	file.cache = 0
	file.cacheLine = 1

	return file, err
}

type RawFile struct {
	file *os.File
	locker sync.Mutex
	cache int
	cacheLine int
	buf *bytes.Buffer
}
