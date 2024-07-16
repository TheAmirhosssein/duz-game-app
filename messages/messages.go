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

func validateKeys(validKeys []string, mapToCheck map[string]string) bool {
	var validKeysCount int
	for key := range mapToCheck {
		for _, validKey := range validKeys {
			if key == validKey {
				validKeysCount++
			}
		}
	}
	return validKeysCount == len(validKeys)
}

func GetMessageType(message *[]byte) (string, error) {
	result := make(map[string]string)
	json.Unmarshal(*message, &result)
	if !isValidType(result["type"]) {
		return "", errors.New("type is not valid")
	}
	messageType := result["type"]
	delete(result, "type")
	*message, _ = json.Marshal(result)
	return messageType, nil
}

func ParseMessageForMatch(message []byte) (map[string]string, error) {
	validKeys := []string{"game_id", "user_id"}
	result := make(map[string]string)
	json.Unmarshal(message, &result)
	if !validateKeys(validKeys, result) {
		return nil, errors.New("json is invalid")
	}
	return result, nil
}
