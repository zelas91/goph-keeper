package request

import (
	"encoding/json"
)

func prettyJSON(data []byte) (string, error) {
	var result []interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	pretty, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return "", err
	}
	return string(pretty), nil
}
