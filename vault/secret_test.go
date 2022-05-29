package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type MockVault struct {
	RequestID string
	Data map[string]string
}

func (v *MockVault) Get(path string) (*string, error) {
	mockData, err := v.getDataJson()
	if err != nil {
		return nil, err
	}

	mockResponse := fmt.Sprintf(`{"request_id": "%s", "data": {"data": %s}}`, v.RequestID, mockData)
	return &mockResponse, nil
}

func (v *MockVault) getDataJson() (string, error) {
	result, err := json.Marshal(v.Data)
	if err != nil {
		return "", errors.New("error marshalling test data")
	}

	return string(result), nil
}

func TestSecretFetch(t *testing.T) {
	mockVault := &MockVault{
		RequestID: "abcdef",
		Data: map[string]string{
			"KEY_1": "VALUE_1",
			"KEY_2": "VALUE_2",
		},
	}

	secretResult, err := SecretFetch("example/secret", mockVault)
	if err != nil {
		t.Errorf("unexpected error doing mock secret fetch:  %s", err)
	}

	if secretResult.RequestId != mockVault.RequestID {
		t.Errorf("unexpected request id, wanted `%s` but got `%s`", mockVault.RequestID, secretResult.RequestId)
	}

	for key := range mockVault.Data {
		result, err := secretResult.GetSubKey(key)
		if err != nil {
			t.Errorf("unexpected error fetching subkey `%s`", key)
		}

		if result != mockVault.Data[key] {
			t.Errorf("unexpected secret value, wanted `%s` but got `%s`", mockVault.Data[key], result)
		}
	}


}