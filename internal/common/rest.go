package common

import (
	"encoding/json"
	"io"
	"net/http"
)

type RestResponse struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg,omitempty"`
	Data   interface{} `json:"data"`
}

func ResponseOk(w http.ResponseWriter, statusCode int, data interface{}) {
	res := &RestResponse{Status: "ok", Data: data}
	b, err := json.Marshal(res)
	if err == nil {
		ResponseError(w, http.StatusInternalServerError, nil, "invalid response")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	io.WriteString(w, string(b))
}

func ResponseError(w http.ResponseWriter, statusCode int, data interface{}, msg string) {
	res := &RestResponse{Status: "error", Msg: msg, Data: data}
	b, _ := json.Marshal(res)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	io.WriteString(w, string(b))
}
