package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/priykumar/notification-service/model"
)

func GenerateResponse(w http.ResponseWriter, code int, message string) {
	resp := model.APIResponse{
		Code:    code,
		Message: message,
	}

	json.NewEncoder(w).Encode(resp)
	jsonBytes, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("Generating response: %s\n\n", string(jsonBytes))
}
