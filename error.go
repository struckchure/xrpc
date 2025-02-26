package xrpc

import (
	"encoding/json"
)

type XRPCError struct {
	Code   int `json:"-"`
	Detail any `json:"detail"`
}

func (e *XRPCError) Error() string {
	out, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(out)
}
