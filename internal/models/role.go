package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Role string

type Roles []Role

func (r *Roles) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan into Roles")
	}

	return json.Unmarshal(bytes, r)
}

func (r Roles) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}
