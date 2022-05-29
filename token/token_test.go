package token

import (
	"fmt"
	"testing"
)

func TestNewToken(t *testing.T) {
	tokenPath := "path/to/the/token"
	tokenSubPath := "sub_token_key"
	rawToken := fmt.Sprint("%%%", tokenPath, ":", tokenSubPath, "%%%")

	parsedToken, err := NewToken(rawToken)
	if err != nil {
		t.Errorf("unexpected error parsing token: %s", err)
	}

	// make sure the outputs are expected
	if parsedToken.TokenName != rawToken {
		t.Errorf("unexpected token name, expected `%s`, got `%s`", rawToken, parsedToken.TokenName)
	}

	if parsedToken.TokenPath != tokenPath {
		t.Errorf("unexpected token path, expected `%s`, got `%s`", tokenPath, parsedToken.TokenPath)
	}

	if parsedToken.TokenSubKey != tokenSubPath {
		t.Errorf("unexpected token subkey, expected `%s`, got `%s`", tokenSubPath, parsedToken.TokenSubKey)
	}
}