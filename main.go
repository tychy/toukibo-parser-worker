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

type HoujinExecutive struct {
	Name     string `json:"氏名"`
	Position string `json:"役職"`
}

type HoujinPreferredStock struct {
	Type   string `yaml:"Type"`
	Amount int    `yaml:"Amount"`
}

type Houjin struct {
	ToukiboCreatedAt      time.Time              `json:"登記簿作成時刻"`
	HoujinName            string                 `json:"法人名"`
	HoujinKaku            string                 `json:"法人格"`
	HoujinAddress         string                 `json:"住所"`
	HoujinCapital         int                    `json:"資本金"`
	HoujinStock           int                    `json:"発行済み株式数"`
	HoujinPreferredStock  []HoujinPreferredStock `json:"各種の株式の数"`
	HoujinExecutives      []HoujinExecutive      `json:"役員"`
	HoujinExecutiveNames  []string               `json:"役員氏名"`
	HoujinRepresentatives []string               `json:"代表者氏名"`
	HoujinBankruptedAt    string                 `json:"破産日"`
	HoujinDissolvedAt     string                 `json:"解散日"`
	HoujinContinuedAt     string                 `json:"会社継続日"`
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
		representativeNames, err := h.GetHoujinRepresentativeNames()
		if err != nil {
			panic(err)
		}
		exectiveNames, err := h.GetHoujinExecutiveNames()
		if err != nil {
			panic(err)
		}

		exectives, err := h.GetHoujinExecutives()
		if err != nil {
			panic(err)
		}
		var houjinExecutives []HoujinExecutive
		for _, e := range exectives {
			houjinExecutives = append(houjinExecutives, HoujinExecutive{
				Name:     e.Name,
				Position: e.Position,
			})
		}

		stock := h.GetHoujinStock()

		var houjinPreferredStock []HoujinPreferredStock
		for _, p := range stock.Preferred {
			houjinPreferredStock = append(houjinPreferredStock, HoujinPreferredStock{
				Type:   p.Type,
				Amount: p.Amount,
			})
		}

		houjin := &Houjin{
			ToukiboCreatedAt:      h.GetToukiboCreatedAt(),
			HoujinName:            h.GetHoujinName(),
			HoujinKaku:            h.GetHoujinKaku(),
			HoujinAddress:         h.GetHoujinAddress(),
			HoujinCapital:         h.GetHoujinCapital(),
			HoujinStock:           stock.Total,
			HoujinPreferredStock:  houjinPreferredStock,
			HoujinExecutives:      houjinExecutives,
			HoujinExecutiveNames:  exectiveNames,
			HoujinRepresentatives: representativeNames,
			HoujinBankruptedAt:    h.GetHoujinBankruptedAt(),
			HoujinDissolvedAt:     h.GetHoujinDissolvedAt(),
			HoujinContinuedAt:     h.GetHoujinContinuedAt(),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(houjin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	workers.Serve(nil) // use http.DefaultServeMux
}
