package render

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err := w.Write(buf.Bytes())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
