package common

import "github.com/google/uuid"

func GenerateUUIDStr() string {
	return uuid.New().String()
}
