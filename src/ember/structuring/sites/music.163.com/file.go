package m1c

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

	return file, err
}

func (p *RawFile) OpenScanner() (scanner Scanner, err error) {
	f, err := os.OpenFile("music.163.com." + "binlog.backup", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		return scanner, err
	}
	return Scanner{bufio.NewScanner(f), f}, err
}

func (p *Scanner) Close() ( err error) {
	p.f.Close()
	return
}

type Scanner struct {
	scanner *bufio.Scanner
	f *os.File
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
