package common

import (
	"encoding/json"
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
	a.writer.WriteHeader(http.StatusInternalServerError)
	a.writer.Write(j)
	return nil
}

func (a *APIPresenter) PresentData(topic string, data interface{}) error {
	j, _ := json.Marshal(data)
	a.writer.Write(j)
	return nil
}
func (a *APIPresenter) PresentNotFound(topic string, data interface{}) error {
	a.writer.WriteHeader(http.StatusNotFound)
	return nil
}
func (a *APIPresenter) PresentUnauthorized(topic string, data interface{}) error {
	a.writer.WriteHeader(http.StatusUnauthorized)
	return nil
}