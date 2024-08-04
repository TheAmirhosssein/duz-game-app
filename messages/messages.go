package messages

import (
	"encoding/json"
	"errors"
	"strings"
)

func isValidType(checkType string) bool {
	validTypes := [3]string{"join_match", "move", "remove"}
	for _, validType := range validTypes {
		if checkType == validType {
			return true
		}
	}
	return false
}

func toCamelCase(snakeStr string) string {
	parts := strings.Split(snakeStr, "_")
	for i, part := range parts {
		if i != 0 && len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

func validateKeys(validKeys []string, mapToCheck *map[string]string) bool {
	var validKeysCount int
	for key := range *mapToCheck {
		for _, validKey := range validKeys {
			if key == validKey {
				validKeysCount++
			}
		}
	}

	CamelCaseMap := make(map[string]string)
	for key, value := range *mapToCheck {
		CamelCaseMap[toCamelCase(key)] = value
		delete(*mapToCheck, key)
	}
	*mapToCheck = CamelCaseMap
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

func ParseMessage(validKeys []string, message []byte) (map[string]string, error) {
	result := make(map[string]string)
	json.Unmarshal(message, &result)
	if !validateKeys(validKeys, &result) {
		return nil, errors.New("json is invalid")
	}
	return result, nil
}

func GenerateMessage(messageType, userId, gameId string, message any) []byte {
	generatedMessage := map[string]any{"type": messageType, "user_id": userId, "message": message, "game_id": gameId}
	serializedMessage, _ := json.Marshal(generatedMessage)
	return serializedMessage
}
