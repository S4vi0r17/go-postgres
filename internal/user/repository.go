package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository es la capa de acceso a datos.
// No tiene lógica de negocio, solo habla con la DB.
type Repository struct {
	db *gorm.DB
}

// NewRepository es el "constructor". En Go no hay `new` ni DI automática,
// se construyen las dependencias a mano y se inyectan explácitamente.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).First(&u, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	u.ID = uuid.New()
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *Repository) Save(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *Repository) Delete(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Delete(u).Error
}
