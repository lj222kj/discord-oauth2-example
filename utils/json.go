package utils

import (
	"encoding/json"
	"net/http"
)

func ParseBodyResponse(r *http.Response, target any) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&target)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}
