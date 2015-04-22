package master

import (
	"encoding/binary"
	"hash/crc32"
)

func (p *Data) read(nid int) (str string, err error) {
	// TODO read raw file and check crc32
	return p.file.Read()
}

func (p *Data) write(b string, nid int) (err error) {
	// TODO write crc32 check
	buf := make([]byte, 8)

	size_buf := buf[:4]
	binary.LittleEndian.PutUint32(size_buf, uint32(len(b)))

	crc_buf := buf[4:]
	binary.LittleEndian.PutUint32(crc_buf, uint32(crc32.ChecksumIEEE([]byte(b))))

	p.file.Write(buf, b)
	return err
}

func NewData(fileName string) Data{
	file, err := NewRawFile(fileName)
	if err != nil {
		println(err.Error())
	}
	return Data { file }
}

type Data struct {
	file RawFile
}
