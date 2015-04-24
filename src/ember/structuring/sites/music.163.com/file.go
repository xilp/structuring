package m1c

import (
	"bufio"
	"bytes"
	"hash/crc32"
	"errors"
	"encoding/binary"
	"io"
	"os"
	"sync"
)

func (p *RawFile) Write(head, line []byte) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.buf.Write(head)
	p.buf.Write(line)

	p.backupBuf.Write(line)

	p.cache += 1
	if p.cache >= p.cacheLine {
		p.Flush()
	}

	return err
}

func (p *RawFile) Read() (ret string, err error) {
	//TODO read raw file
	return ret, err
}

func (p *RawFile) Flush() (err error) {
	n, err := io.Copy(p.fd, p.buf)
	_ = n
	if err != nil {
		println(err.Error())
		return
	}
	n, err = io.Copy(p.backupFd, p.backupBuf)
	if err != nil {
		println(err.Error())
		return
	}
	p.buf.Reset()
	p.backupBuf.Reset()
	p.cache= 0
	return err
}

func NewRawFile(path string) (file RawFile, err error) {
	file.fileName = path + "/binlog.txt"
	file.fd, err = os.OpenFile(file.fileName, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		println(err.Error())
		return
	}

	file.backupFileName = path + "/binlog.backup"
	file.backupFd, err = os.OpenFile(file.backupFileName, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		println(err.Error())
		return
	}

	file.buf = bytes.NewBuffer([]byte{})
	file.backupBuf = bytes.NewBuffer([]byte{})
	file.cache = 0
	file.cacheLine = 1

	return
}

func (p *RawFile) Close() (err error) {
	err = p.fd.Close()
	if err != nil {
		return errors.New("close " + p.fileName + " error")
	}
	err = p.backupFd.Close()
	if err != nil {
		return errors.New("close " + p.backupFileName + " error")
	}
	return err
}

func (p *RawFile) OpenScanner() (scanner Scanner, err error) {
	fd, err := os.OpenFile(p.fileName, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		println(err.Error())
		return scanner,err
	}
	stat, err := fd.Stat()
	if err != nil {
		return scanner, err
	}
	return Scanner {
		scanner:bufio.NewReaderSize(fd, 1024*1024),
		fd:fd,
		head:make([]byte, 8),
		line:make([]byte,1024*8),
		size:stat.Size(), idx:0 }, err
}

func (p *Scanner) ReadNByte(buf []byte, size uint32) (err error) {
	var n int
	var sum uint32
	p.idx += int64(size)
	if p.idx > p.size {
		return errors.New("EOF")
	}

	for sum = 0; sum < size; sum = sum + uint32(n) {
		n = int(sum)
		n, err = p.scanner.Read(buf[n:size])
		if err != nil {
			return err
		}
	}
	return
}

func (p *Scanner) Scan() (buf []byte, err error) {
	for {
		err = p.ReadNByte(p.head, 8)
		if err != nil {
			return buf, err
		}

		size := binary.LittleEndian.Uint32(p.head[0:4])
		crc32Line := binary.LittleEndian.Uint32(p.head[4:])

		err = p.ReadNByte(p.line, size)
		if err != nil {
			return buf, err
		}

		crc32Check := uint32(crc32.ChecksumIEEE([]byte(p.line[:size])))
		if crc32Check == crc32Line {
			return p.line[:size], err
		}
		if crc32Check != crc32Line {
			println("crc32 check fault")
			continue
		}
	}
	return
}

func (p *Scanner) Close() ( err error) {
	p.fd.Close()
	return
}

type Scanner struct {
	scanner *bufio.Reader
	fd *os.File
	head []byte
	line []byte
	size int64
	idx int64
}

type RawFile struct {
	fd *os.File
	backupFd *os.File
	fileName string
	backupFileName string
	locker sync.Mutex
	cache int
	cacheLine int
	buf *bytes.Buffer
	backupBuf *bytes.Buffer
}
