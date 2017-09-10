package client

import "encoding/json"

type APIResponse struct {
	Success   bool            `json:"success"`
	Result    json.RawMessage `json:"result"`
	Message   string          `json:"msg"`
	ErrorCode string          `json:"error_code"`
}
