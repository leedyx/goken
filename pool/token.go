package pool

import (
	"encoding/json"
	"fmt"
)

type Token struct {
	Id              int64  `json:"id"`
	ExpireTimestamp int64  `json:"expireTimestamp"`
	Sig             string `json:"sig"`
	NcToken         string `json:"nc_token"`
	SessionId       string `json:"nc_cessionid"`
}

func (token *Token) toJson() ([]byte, error) {

	data, err := json.Marshal(token)
	if err != nil {
		return nil, fmt.Errorf("to json error ! %w", err)
	}

	return data, nil
}
