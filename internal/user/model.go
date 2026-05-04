package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User es la entidad. Las tags GORM definen el schema en Postgres,
// las tags json definen cómo se serializa a JSON.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName le dice a GORM cómo llamar a la tabla.
// Por default usará "users" igual, pero lo ponemos explícito.
func (User) TableName() string { return "users" }

// DTOs ----------------------------------------------------------------

// CreateDTO equivale a CreateUserDto (class-validator).
// Los tags `validate:"..."` los lee go-playground/validator.
type CreateDTO struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Email string `json:"email" validate:"required,email,max=255"`
}

// UpdateDTO usa punteros para distinguir "no enviado" de "enviado vacío".
// Si el cliente no manda el campo, el puntero queda nil y no lo tocamos.
type UpdateDTO struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Email *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
}

// Errores de dominio --------------------------------------------------
// Sentinel errors: se comparan con errors.Is en los handlers.

var (
	ErrNotFound    = errors.New("user not found")
	ErrEmailExists = errors.New("email already exists")
	ErrInvalidID   = errors.New("invalid user id")
)
