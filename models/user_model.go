package models

import (
	"strings"
	"time"
)

type User struct {
	ID    int       `json:"id" bson:"id" validate:"required"`
	Name  string    `json:"name" bson:"name" validate:"required"`
	DOB   DateField `json:"dob" bson:"dob" validate:"required"`
	Email string    `json:"email" bson:"email" validate:"required"`
}

type DateField struct {
	time.Time
}

func (t *DateField) UnmarshalJSON(b []byte) (err error) {

	date, err := time.Parse("2006-01-02", strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}
	t.Time = date
	return
}
