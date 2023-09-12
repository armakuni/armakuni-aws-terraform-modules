package test

import (
	"encoding/json"
	"fmt"
)

func structToMap(input interface{}) (map[string]interface{}, error) {
	var output map[string]interface{}
	jsonStr, err := json.Marshal(input)
	err = json.Unmarshal(jsonStr, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %s", err.Error())
	}
	return output, nil
}
