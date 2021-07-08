package common_adapters

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
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

func NewFileStreamPresenter(writer http.ResponseWriter) *APIPresenter {
	return &APIPresenter{writer}
}

func (f *FileStreamPresenter) PresentError(topic string, err interface{}) error {
	j, _ := json.Marshal(err)
	f.writer.WriteHeader(http.StatusInternalServerError)
	f.writer.Write(j)
	return nil
}

func (f *FileStreamPresenter) PresentData(topic string, file File) error {
	f.writer.Header().Set("Content-Length", spew.Sprint("%d,", file.ContentLength()))
	f.writer.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
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
