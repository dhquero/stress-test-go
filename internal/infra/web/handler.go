package web

import (
	"encoding/json"
	"net/http"
	"time"
)

type Output struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	output := Output{
		Status: "OK",
		Time:   time.Now().Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
