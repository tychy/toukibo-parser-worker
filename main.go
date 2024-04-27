package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/syumai/workers"
	"github.com/tychy/toukibo-parser/pdf"
	"github.com/tychy/toukibo-parser/toukibo"
)

func readPdf(data []byte) (string, error) {
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		msg := "Hello!"
		w.Write([]byte(msg))
	})
	http.HandleFunc("/parse", func(w http.ResponseWriter, req *http.Request) {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		text, err := readPdf(b)
		if err != nil {
			panic(err)
		}
		h, err := toukibo.Parse(text)
		if err != nil {
			panic(err)
		}
		kaku := h.GetHoujinKaku()
		if err != nil {
			panic(err)
		}
		io.Copy(w, bytes.NewReader([]byte(kaku)))
		// curl -X POST -H "Content-Type: application/pdf" --data-binary "@sample.pdf" https://go-worker.a2sin2a2ko1115.workers.dev/parse
		//  http://localhost:8787
		// curl -X POST -H "Content-Type: application/pdf" --data-binary "@sample.pdf"  http://localhost:8787/parse
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
