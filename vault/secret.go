package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Secret struct {
	RequestId string
	Contents map[string]string
}

func (s *Secret) GetSubKey(subKeyName string) (string, error) {
	if secretValue, ok := s.Contents[subKeyName]; ok {
		return secretValue, nil
	} else {
		return "", errors.New(fmt.Sprint("invalid subKey: `",subKeyName,"`"))
	}
}


func SecretFetch(path string, client SecretGetter) (*Secret, error) {
	// build the full vault api path from the provided path, which just ends up being
	// a split on "/" with the following structure v1/<first entry>/data</any/other/things/here>
	pathParts := strings.Split(path, "/")
	apiPath := "/" + strings.Join(append([]string{"v1", pathParts[0], "data"}, pathParts[1:]...), "/")

	// now that the api path is correct, use the client to fetch the secret
	result, err := client.Get(apiPath)
	if err != nil {
		return nil, err
	}

	var SecretResponse struct {
		RequestId string `json:"request_id"`
		Data struct {
			Data map[string]string `json:"data"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(*result), &SecretResponse); err != nil {
		return nil, errors.New(fmt.Sprint("error unmarshalling secret response", err))
	}

	return &Secret{
		RequestId: SecretResponse.RequestId,
		Contents: SecretResponse.Data.Data,
	}, nil
}
