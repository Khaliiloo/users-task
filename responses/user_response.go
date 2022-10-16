package responses

import (
	"encoding/json"
	"users-task/configs"
)

type UserResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

func (response UserResponse) Log() *UserResponse {
	jsonResponse, _ := json.Marshal(response)
	configs.Logger.Info(string(jsonResponse))
	return &response
}
