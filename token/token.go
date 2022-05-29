package token

import (
	"errors"
	"strings"
)

const tokenPattern = "%%%.+?%%%"



type Token struct {
	TokenName string
	TokenPath string
	TokenSubKey string
}

func NewToken(rawToken string) (*Token, error) {
	// tokens come in wrapped in %%%tokenPath:tokenSubKey%%%, strip that off
	tokenValue := strings.Replace(rawToken, "%%%", "", -1)
	if tokenValue == "" {
		return nil, errors.New("blank token provided")
	}

	// vault based tokens come in with the form of <path in vault>:<sub key in json>, split on the colon
	tokenPath := strings.SplitN(tokenValue, ":", 2)

	// do some sanity checks
	if tokenPath[0] == "" {
		return nil, errors.New("blank vault path provided")
	}

	if len(tokenPath) < 2 || tokenPath[1] == "" {
		return nil, errors.New("blank vault sub-key provided")
	}

	return &Token{
		TokenName: rawToken,
		TokenPath: tokenPath[0],
		TokenSubKey: tokenPath[1],
	}, nil
}
