package trpc

import (
	"encoding/json"
)

type TRPCError struct {
	Code   int `json:"-"`
	Detail any `json:"detail"`
}

func (e *TRPCError) Error() string {
	out, err := json.Marshal(e)
	if err != nil {
		return ""
	}

	return string(out)
}
