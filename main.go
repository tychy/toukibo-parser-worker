package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

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

type Houjin struct {
	ToukiboCreatedAt   time.Time `json:"toukibo_created_at"`
	HoujinName         string    `json:"houjin_name"`
	HoujinKaku         string    `json:"houjin_kaku"`
	HoujinAddress      string    `json:"houjin_address"`
	HoujinBankruptedAt string    `json:"houjin_bankrupted_at"`
	HoujinDissolvedAt  string    `json:"houjin_dissolved_at"`
	HoujinCapital      int       `json:"houjin_capital"`
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
		houjin := &Houjin{
			ToukiboCreatedAt:   h.GetToukiboCreatedAt(),
			HoujinName:         h.GetHoujinName(),
			HoujinKaku:         h.GetHoujinKaku(),
			HoujinAddress:      h.GetHoujinAddress(),
			HoujinBankruptedAt: h.GetHoujinBankruptedAt(),
			HoujinDissolvedAt:  h.GetHoujinDissolvedAt(),
			HoujinCapital:      h.GetHoujinCapital(),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(houjin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	workers.Serve(nil) // use http.DefaultServeMux
}
