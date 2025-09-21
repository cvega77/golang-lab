package handler

import (
	"github.com/google/uuid"
)

func IsValidUUID(uuidString string) bool {
	_, err := uuid.Parse(uuidString)
	return err == nil
}
