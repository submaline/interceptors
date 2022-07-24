package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type genTokenRequest struct {
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	ReturnSecureToken bool   `json:"return_secure_token,omitempty"`
}

type genTokenResponse struct {
	Kind         string `json:"kind"`
	LocalId      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IdToken      string `json:"IdToken"`
	Registered   bool   `json:"registered"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"ExpiresIn"`
}

// todo : error handle
type genTokenError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
		Status  string `json:"status"`
		Details []struct {
			Type     string `json:"@type"`
			Reason   string `json:"reason"`
			Domain   string `json:"domain"`
			Metadata struct {
				Service string `json:"service"`
			} `json:"metadata"`
		} `json:"details"`
	} `json:"error"`
}

type TokenData struct {
	IdToken   string
	Refresh   string
	ExpiresAt time.Time
	ExpiresIn string // 互換
	UID       string // 互換
}

func GenToken(email, password string) (*TokenData, error) {
	url_ := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s",
		os.Getenv("FIREBASE_WEB_API_KEY"))
	bin := genTokenRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
	dataBin, err := json.Marshal(bin)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(url_, "application/json", bytes.NewBuffer(dataBin))
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var result genTokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	expiresIn, err := strconv.ParseInt(result.ExpiresIn, 10, 64)
	if err != nil {
		var resultErr genTokenError
		_ = json.Unmarshal(body, &resultErr)
		return nil, errors.Wrap(err, resultErr.Error.Message)
	}

	now := time.Now()
	expiresAt := now.Add(time.Second * time.Duration(expiresIn))
	return &TokenData{
		IdToken:   result.IdToken,
		Refresh:   result.RefreshToken,
		ExpiresAt: expiresAt,
		ExpiresIn: result.ExpiresIn,
		UID:       result.LocalId,
	}, nil
}
