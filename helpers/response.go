package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(w http.ResponseWriter, httpCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	response := Response{
		Status:  "Success",
		Message: "Success",
		Data:    payload,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}

}

func ResponseError(w http.ResponseWriter, httpCode int, errStatus, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	response := Response{
		Status:  errStatus,
		Message: errMsg,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err)
	}
}
