package master

import (
	"encoding/binary"
	"hash/crc32"
)

func (p *Data) read(nid int) (str string, err error) {
	/*
	buf := bytes.NewBuffer(b) // b is []byte
	myfirstint, err := binary.ReadVarint(buf)
	anotherint, err := binary.ReadVarint(buf)
	var i int16 = 41
	err := binary.Write(w, binary.LittleEndian, i)
	var i int16 = 41
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	*/
	/*
	for _, byte := range size_buf {
		fmt.Printf("%02X ", byte)
	}
	println()
	for _, byte := range tmp_02 {
		fmt.Printf("%02X ", byte)
	}
	println()
	*/
	return p.file.Read()
}

func (p *Data) write(b string, nid int) (err error) {
	buf := make([]byte, 8)

	size_buf := buf[:4]
	binary.LittleEndian.PutUint32(size_buf, uint32(len(b)))

	crc_buf := buf[4:]
	binary.LittleEndian.PutUint32(crc_buf, uint32(crc32.ChecksumIEEE([]byte(b))))

	//p.file.Write(append(buf, b))
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
