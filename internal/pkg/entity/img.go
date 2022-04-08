package entity

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"io/ioutil"
	"mime/multipart"
	"strings"
)

type Image struct {
	Id   string
	Name string
	Ext  string
	Data []byte
}

func NewImage(file multipart.File, fileName string) *Image {
	id, _ := uuid.NewV4()
	data, err := ioutil.ReadAll(file)
	if err != nil {

	}
	return &Image{
		Id:   id.String(),
		Name: fileName,
		Ext:  resolveExt(fileName),
		Data: data,
	}

}

func (i *Image) ToBytes() []byte {
	data, err := json.Marshal(i)
	if err != nil {
		zap.S().Errorf("Error during serializing img: %s_%s", i.Id, i.Name)
	}
	return data
}

func resolveExt(fileName string) string {
	r := strings.Split(fileName, ".")
	if len(r) < 2 {
		return ""
	} else if len(r) == 2 && r[0] == "" {
		return ""
	}
	return r[len(r)-1]
}
