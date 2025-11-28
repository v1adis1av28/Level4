package network

import (
	"encoding/json"
	"io"
	"mygrep/internal/config"
	"mygrep/internal/grep"
	"net/http"
)

func StartWorkerServer(port string) error {
	http.HandleFunc("/grep", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only post request allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid json data", http.StatusBadRequest)
			return
		}
		var req GrepRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		flags := &config.Flags{
			IgnoreFlag:         req.Flags.Ignore,
			InvertFlag:         req.Flags.Invert,
			StrictStringFlag:   req.Flags.Strict,
			LineNumberFlag:     req.Flags.LineNumbers,
			CountOfLineFlag:    req.Flags.CountOnly,
			AdditionalLineFlag: req.Flags.After,
			PreviousLineFlag:   req.Flags.Before,
			AroundLineFlag:     req.Flags.Around,
		}

		conf := &config.Config{
			Pattern: req.Pattern,
			Flags:   flags,
			Buffer:  make([]config.Line, 0),
		}

		lines, count, err := grep.GrepFromData(req.Data, conf)

		resp := GrepResponse{
			Lines: lines,
			Count: count,
			Error: "",
		}
		if err != nil {
			resp.Error = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	return http.ListenAndServe(":"+port, nil)
}
