package vault

import "testing"

func TestNewSecretCollection(t *testing.T) {
	c := NewSecretCollection()

	testSecretKey := "secret/path"
	testSubKey := "subkey"
	testValue := "value!"

	// no secret is there, this should be false
	if c.HasSecret(testSecretKey) {
		t.Errorf("unexpected secret entry at %s", testSecretKey)
	}

	// add a secret entry
	c.AddSecret(testSecretKey, &Secret{
		RequestId: "1234",
		Contents: map[string]string{
			testSubKey: testValue,
		},
	})

	// make sure we have the secret
	if !c.HasSecret(testSecretKey) {
		t.Errorf("expected secret entry of %s", testSecretKey)
	}

	// fetch the secret out and make sure it matches
	storedSecret, err := c.Secrets[testSecretKey].GetSubKey(testSubKey)
	if err != nil {
		t.Errorf("unexpected error fetching subkey %s", err)
	}

	if storedSecret != testValue {
		t.Errorf("unexpected secret value, expected `%s` got `%s`", testValue, storedSecret)
	}
}