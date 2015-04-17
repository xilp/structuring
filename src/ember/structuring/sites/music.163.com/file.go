package m1c

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

func (p *RawFile) Write(buf []byte, str string) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.buf.WriteString(string(buf) + str)
	p.backupBuf.WriteString(str)
	p.cache+= 1
	if p.cache >= p.cacheLine {
		p.Flush()
	}
	return err
}

func (p *RawFile) Read() (ret string, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	return ret, err
}

func (p *RawFile) Close() (err error) {
	p.file.Close()
	err = p.backupFile.Close()
	return err
}

func (p *RawFile) Flush() (err error) {
	// TODO lock 
	io.Copy(p.file, p.buf)
	io.Copy(p.backupFile, p.backupBuf)
	p.buf.Reset()
	p.backupBuf.Reset()
	p.cache= 0
	return err
}

func NewRawFile() (file RawFile, err error){
	file.file, err = os.OpenFile("music.163.com." + "binlog.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		return
	}
	file.backupFile, err = os.OpenFile("music.163.com." + "binlog.backup", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		return
	}
	file.buf = bytes.NewBuffer([]byte{})
	file.backupBuf = bytes.NewBuffer([]byte{})
	file.cache = 0
	file.cacheLine = 1

	fmt.Printf("[file.file:%#v]\n", file.file)
	stat, err := file.file.Stat()
	fmt.Printf("[file.Stat():%#v]\n", stat)
	fd := file.file.Fd()
	fmt.Printf("[file.Fd():%#v]\n", fd)

	seek, err := file.file.Seek(-2048, 2)
	fmt.Printf("[file.Seek():%#v]\n", seek)
	if err != nil {
		fmt.Printf("[err:%v]\n", err.Error())
	}
	//magicNum size

	return file, err
}

type RawFile struct {
	file *os.File
	backupFile *os.File
	locker sync.Mutex
	cache int
	cacheLine int
	buf *bytes.Buffer
	backupBuf *bytes.Buffer
}
