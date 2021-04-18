package share

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status int         `json:"-"`
	Code   int         `json:"code,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Msg    string      `json:"msg,omitempty"`
}

func (r *Response) JsonBytes() []byte {
	d, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	return d
}

func (r *Response) WriteJson(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	if r.Status != 0 {
		w.WriteHeader(r.Status)
	}
	_, err := w.Write(r.JsonBytes())
	return err
}
