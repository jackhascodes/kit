package common_adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIPresenter struct {
	writer     http.ResponseWriter
}

func NewAPIPresenter(writer http.ResponseWriter) *APIPresenter {
	return &APIPresenter{writer}
}

func (a *APIPresenter) PresentError(topic string, err interface{}) error {
	j, _ := json.Marshal(err)
	fmt.Println(fmt.Sprintf("writing error to http writer: %s", j))
	a.writer.Header().Set("Content-Type", "application/json")
	a.writer.WriteHeader(http.StatusInternalServerError)
	written, werr := a.writer.Write(j)
	fmt.Println(fmt.Sprintf("wrote %d bytes with error: %s", written, werr))


	return nil
}

func (a *APIPresenter) PresentData(topic string, data interface{}) error {
	j, _ := json.Marshal(data)

	a.writer.Header().Set("Content-Type", "application/json")
	a.writer.Write(j)
	if f, ok := a.writer.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}
func (a *APIPresenter) PresentNotFound(topic string, data interface{}) error {
	fmt.Println("writing not found response")
	a.writer.WriteHeader(http.StatusNotFound)
	return nil
}
func (a *APIPresenter) PresentUnauthorized(topic string, data interface{}) error {
	fmt.Println("writing unauthorized response")
	a.writer.WriteHeader(http.StatusUnauthorized)
	return nil
}