package common

import (
	"crypto/rand"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

func GenerateUniqueString() string {
	uuid1 := strings.ReplaceAll(uuid.New().String(), "-", "")
	uuid2 := strings.ReplaceAll(uuid.New().String(), "-", "")
	return uuid1 + uuid2[:16]
}

func ToJSON(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func FromJSON(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateBase62ShortURL() string {
	return generateRandomBase62(7)
}

func generateRandomBase62(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return generateBase62FromUUID(length)
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = base62Chars[int(bytes[i])%62]
	}
	return string(result)
}

func generateBase62FromUUID(length int) string {
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex := int(uuidStr[i%len(uuidStr)]) % 62
		result[i] = base62Chars[charIndex]
	}
	return string(result)
}
