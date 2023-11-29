package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type KernelDepositView struct {
	Chain        string `json:"chain"`
	DepositHash  string `json:"deposit_hash"`
	DepositIndex int64  `json:"deposit_index"`
}

type Output struct {
	Type               string    `json:"type"`
	OutputID           string    `json:"output_id"`
	TransactionHash    string    `json:"transaction_hash"`
	OutputIndex        uint      `json:"output_index"`
	Asset              string    `json:"asset"`
	Amount             string    `json:"amount"`
	Mask               string    `json:"mask"`
	Keys               []string  `json:"keys"`
	SendersHash        string    `json:"senders_hash"`
	SendersThreshold   int64     `json:"senders_threshold"`
	Senders            []string  `json:"senders"`
	ReceiversHash      string    `json:"receivers_hash"`
	ReceiversThreshold int64     `json:"receivers_threshold"`
	Receivers          []string  `json:"receivers"`
	Extra              string    `json:"extra"`
	State              string    `json:"state"`
	Sequence           int64     `json:"sequence"`
	CreatedAt          time.Time `json:"created_at"`
	Signers            []string  `json:"signers"`
	SignedBy           string    `json:"signed_by"`

	Deposit   *KernelDepositView `json:"deposit,omitempty"`
	RequestId string             `json:"request_id,omitempty"`
}

func ListUnspentOutputs(ctx context.Context, membersHash string, threshold byte, assetId string, u *SafeUser) ([]*Output, error) {
	method, path := "GET", fmt.Sprintf("/safe/outputs?members=%s&threshold=%d&asset=%s&state=unspent", membersHash, threshold, assetId)
	token, err := SignAuthenticationToken(u.UserId, u.SessionId, u.SessionKey, method, path, "")
	body, err := Request(ctx, method, path, []byte{}, token)
	if err != nil {
		return nil, ServerError(ctx, err)
	}
	var resp struct {
		Data  []*Output `json:"data"`
		Error Error     `json:"error"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, BadDataError(ctx)
	}
	if resp.Error.Code > 0 {
		return nil, resp.Error
	}
	return resp.Data, nil
}

func ListOutputs(ctx context.Context, membersHash string, threshold byte, assetId, state, offset, order string, limit int, u *SafeUser) ([]*Output, error) {
	v := url.Values{}
	v.Set("members", membersHash)
	v.Set("threshold", fmt.Sprint(threshold))
	v.Set("limit", strconv.Itoa(limit))
	if offset != "" {
		v.Set("offset", offset)
	}
	if assetId != "" {
		v.Set("asset", assetId)
	}
	if state != "" {
		v.Set("state", state)
	}
	if order != "" {
		v.Set("order", order)
	}
	method, path := "GET", fmt.Sprintf("/safe/outputs?"+v.Encode())
	token, err := SignAuthenticationToken(u.UserId, u.SessionId, u.SessionKey, method, path, "")
	body, err := Request(ctx, method, path, []byte{}, token)
	if err != nil {
		return nil, ServerError(ctx, err)
	}
	var resp struct {
		Data  []*Output `json:"data"`
		Error Error     `json:"error"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, BadDataError(ctx)
	}
	if resp.Error.Code > 0 {
		return nil, resp.Error
	}
	return resp.Data, nil
}
