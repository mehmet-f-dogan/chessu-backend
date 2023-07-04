package lib

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type UUID struct {
	pgtype.UUID
}

type UUIDArray struct {
	pgtype.UUIDArray
}

func (a *UUIDArray) Append(uuid string) error {
	uuidValue := pgtype.UUID{}
	err := uuidValue.Set(uuid)
	if err != nil {
		return err
	}

	newArray := append(a.Elements, uuidValue)
	return a.UUIDArray.Set(newArray)
}

func (u UUID) MarshalJSON() ([]byte, error) {
	baseUUID, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return nil, err
	}
	return json.Marshal(baseUUID.String())
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		return err
	}
	pgUUID := pgtype.UUID{}
	err = pgUUID.Set(parsedUUID.String())
	if err != nil {
		return err
	}
	u.UUID = pgUUID
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return &u.UUID, nil
}

func (u *UUID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	pgUUID := &pgtype.UUID{}
	err := pgUUID.Scan(value)
	if err != nil {
		return err
	}
	u.UUID = *pgUUID
	return nil
}

// Implement the Value method for the Valuer interface
func (ua UUIDArray) Value() (driver.Value, error) {
	// Convert the UUIDArray to a format compatible with the database
	// and return the result
	return pgtype.UUIDArray(ua.UUIDArray).Value()
}

// Implement the Scan method for the Scanner interface
func (ua *UUIDArray) Scan(src interface{}) error {
	// Perform the necessary conversion from the database format
	// to the UUIDArray type and assign the result to ua
	return (*pgtype.UUIDArray)(&ua.UUIDArray).Scan(src)
}
