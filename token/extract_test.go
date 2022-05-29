package token

import (
	"fmt"
	"testing"
)

func TestExtractionCollection_ParseLine(t *testing.T) {
	tokenPath := "token/path"
	tokenSubKey := "subKey"
	expectedToken := fmt.Sprint("%%%", tokenPath, ":", tokenSubKey,"%%%")
	testPhrase := fmt.Sprintf("hello, this is a test line with %s", expectedToken)

	e := NewExtractionCollection()

	// make sure that we don't get any errors on parsing
	if err := e.ParseLine(testPhrase); err != nil {
		t.Errorf("Unexpected Error Parsing Line %s", err)
	}

	// ensure that the expected token is present
	value, ok := e.Tokens[expectedToken]
	if !ok {
		t.Errorf("couldn't detect token %s when we should have been able to", expectedToken)
	}

	// validate that the path is parsed correctly
	if value.TokenPath != tokenPath {
		t.Errorf("unexpected token path, expected %s got %s", tokenPath, value.TokenPath)
	}

	// validate the subkey is parsed correctly
	if value.TokenSubKey != tokenSubKey {
		t.Errorf("unexpected token subkey, expected %s got %s", tokenSubKey, value.TokenSubKey)
	}
}
