package token

import (
	"fmt"
	"token-replace/vault"
	"testing"
)

func TestReplacementCollection_ProcessLine(t *testing.T) {
	tokenName := "%%%path/to:subkey%%%"
	tokenValue := "Hello World"

	testString := "We should have this replaced with %s, is it?"
	tokenizedString := fmt.Sprintf(testString, tokenName)
	finalString := fmt.Sprintf(testString, tokenValue)

	// create a extraction collection
	e := NewExtractionCollection()
	if err := e.ParseLine(tokenizedString); err != nil {
		t.Errorf("error parsing tokenized string: %s", err)
	}

	// mock a secret collection, yes this really should be abstracted
	s := vault.NewSecretCollection()
	s.AddSecret("path/to", &vault.Secret{
		RequestId: "1234",
		Contents: map[string]string{
			"subkey": tokenValue,
		},
	})

	r := NewReplacementCollection(e, s)
	line, err := r.ProcessLine(tokenizedString)
	if err != nil {
		t.Errorf("unexpected error processing line: %s", err)
	}

	if line != finalString {
		t.Errorf("unexpected line, expected `%s`, got `%s`", finalString, line)
	}
}
