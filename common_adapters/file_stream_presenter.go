package common_adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FileStreamPresenter struct {
	writer     http.ResponseWriter
}

type File interface{
	Stream() io.Reader
	ContentLength() int64
	Name() string
	Type() string
}

func NewFileStreamPresenter(writer http.ResponseWriter) *FileStreamPresenter {
	return &FileStreamPresenter{writer}
}

func (f *FileStreamPresenter) PresentError(topic string, err interface{}) error {
	j, _ := json.Marshal(err)
	f.writer.Header().Set("Content-Type", "application/json")
	f.writer.WriteHeader(http.StatusInternalServerError)
	f.writer.Write(j)
	return nil
}

func (f *FileStreamPresenter) PresentData(topic string, data interface{}) error {
	file := data.(File)
	f.writer.Header().Set("Content-Length", fmt.Sprintf("%d", file.ContentLength()))
	f.writer.Header().Set("Content-Disposition", "inline")
	f.writer.Header().Set("Content-Type", file.Type())
	io.Copy(f.writer, file.Stream())
	return nil
}
func (f *FileStreamPresenter) PresentNotFound(topic string, data interface{}) error {
	f.writer.WriteHeader(http.StatusNotFound)
	return nil
}
func (f *FileStreamPresenter) PresentUnauthorized(topic string, data interface{}) error {
	f.writer.WriteHeader(http.StatusUnauthorized)
	return nil
}
