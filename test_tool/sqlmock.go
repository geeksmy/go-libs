package test_tool

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type UUIDArg struct{}

func (arg UUIDArg) Match(value driver.Value) bool {
	_, err := uuid.Parse(value.(string))
	return err == nil
}
