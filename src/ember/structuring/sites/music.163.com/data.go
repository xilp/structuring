package m1c

import (
	"encoding/binary"
	"hash/crc32"
)

func (p *Data) load(b string, nid int) (err error) {
	// TODO check file, move offsite of file ptr
	return err
}

func (p *Data) write(b string, nid int) (err error) {
	buf := make([]byte, 8)

	size_buf := buf[:4]
	binary.LittleEndian.PutUint32(size_buf, uint32(len(b)))

	crc_buf := buf[4:]
	binary.LittleEndian.PutUint32(crc_buf, uint32(crc32.ChecksumIEEE([]byte(b))))

	err = p.file.Write(buf, b)
	if err != nil {
		println(err.Error())
		return
	}
	return err
}

func NewData() Data{
	file, err := NewRawFile()
	if err != nil {
		println(err.Error())
	}
	return Data { file }
}

type Data struct {
	file RawFile
}
