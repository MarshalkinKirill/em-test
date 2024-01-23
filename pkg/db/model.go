package db

import "github.com/google/uuid"

type Person struct {
	tableName struct{} `pg:"emsrv.people,alias:t,discard_unknown_columns"`

	PersonID    uuid.UUID `json:"personId" pg:"personId"`
	Name        string    `json:"name" pg:"name" validate:"required,max=25"`
	Surname     string    `json:"surname" pg:"surname" validate:"required,max=25"`
	Patronymic  string    `json:"patronymic" pg:"patronymic" validate:"required,max=25"`
	Age         int       `json:"age,omitempty" pg:"age" validate:"omitempty,gte=0"`
	Gender      string    `json:"gender,omitempty" pg:"gender" validate:"omitempty,oneof=male female other"`
	Nationality string    `json:"nationality,omitempty" pg:"nationality" validate:"omitempty"`
}

type UpdatePerson struct {
	tableName struct{} `pg:"emsrv.people,alias:t,discard_unknown_columns"`

	PersonID    uuid.UUID `json:"personId" pg:"personId" validate:"omitempty"`
	Name        string    `json:"name" pg:"name" validate:"omitempty,max=25"`
	Surname     string    `json:"surname" pg:"surname" validate:"omitempty,max=25"`
	Patronymic  string    `json:"patronymic" pg:"patronymic" validate:"omitempty,max=25"`
	Age         int       `json:"age,omitempty" pg:"age" validate:"omitempty,gte=0"`
	Gender      string    `json:"gender,omitempty" pg:"gender" validate:"omitempty,oneof=Male Female Other"`
	Nationality string    `json:"nationality,omitempty" pg:"nationality" validate:"omitempty"`
}
