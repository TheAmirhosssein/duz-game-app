package messages

import (
	"encoding/json"
	"errors"
)

func isValidType(checkType string) bool {
	validTypes := [1]string{"join_match"}
	for _, validType := range validTypes {
		if checkType == validType {
			return true
		}
	}
	return false
}

func GetMessageType(message []byte) (string, error) {
	result := make(map[string]string)
	json.Unmarshal(message, &result)
	if !isValidType(result["type"]) {
		return "", errors.New("type is not valid")
	}
	return result["type"], nil
}
