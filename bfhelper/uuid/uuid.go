// Package uuid uuid生成
package uuid

import "github.com/google/uuid"

// NewUUID generates a new UUID v4
func NewUUID() string{
	return uuid.New().String()
}