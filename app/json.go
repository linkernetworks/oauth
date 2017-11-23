package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func toJson(w http.ResponseWriter, val interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Error("json encode error", err)
	}
	return err
	// http.Error(w, err.Error(), 200)
}
