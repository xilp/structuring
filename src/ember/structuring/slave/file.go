package slave

import (
	"os"
)

func (p *RawFile) write(str string) (err error) {
	_, err = p.file.WriteString(str)
	if err != nil {
		return err
	}
	return err
}

func NewRawFile() (file RawFile, err error){
	file.file, err = os.OpenFile("binglog.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0640)
	if err != nil {
		return
	}
	return file, err
}

type RawFile struct {
	file *os.File
}
