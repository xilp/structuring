package slave

import (
)

func (p *Data) write(url string, ret []string, nid int) (err error) {
	/*
		存储方式：01\turl\t0\t1\t2\t3\t4\t5\t6\n 典型的行列结构
	*/
	version := "01"
	str := version + "\t" + url
	for _, v := range ret {
		str = str + "\t" +  v
	}
	str = str + "\n"
	p.file.write(str)
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
