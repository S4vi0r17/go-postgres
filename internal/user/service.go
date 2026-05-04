package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service contiene la lógica de negocio.
// Depende del Repository vía struct field (inyección por constructor).
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, dto CreateDTO) (*User, error) {
	// Pre-check de email único. No es 100% atomático (puede haber race),
	// en producción te apoyarías en el UNIQUE constraint de Postgres.
	existing, err := s.repo.FindByEmail(ctx, dto.Email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	u := &User{
		Name:  dto.Name,
		Email: dto.Email,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, dto UpdateDTO) (*User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Merge manual de campos. Cada puntero nil = campo no enviado.
	if dto.Name != nil {
		u.Name = *dto.Name
	}
	if dto.Email != nil {
		u.Email = *dto.Email
	}

	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, u)
}
