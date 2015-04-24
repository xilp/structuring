package m1c

import (
	"encoding/binary"
	"hash/crc32"
)

func (p *Data) load(b string, nid int) (err error) {
	// TODO check file, move offsite of file ptr
	return err
}

func (p *Data) write(line []byte, nid int) (err error) {
	head := make([]byte, 8)

	size := len(line)
	crc32Line := uint32(crc32.ChecksumIEEE([]byte(line)))
	binary.LittleEndian.PutUint32(head[:4], uint32(size))
	binary.LittleEndian.PutUint32(head[4:], crc32Line)

	err = p.file.Write(head, line)
	if err != nil {
		return
	}
	return
}

func NewData(path string) Data{
	file, err := NewRawFile(path)
	if err != nil {
		println(err.Error())
	}
	return Data { file }
}

type Data struct {
	file RawFile
}
