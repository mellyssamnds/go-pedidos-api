package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/repository"
	"golang.org/x/crypto/bcrypt"
)

type ClienteService struct {
	clienteRepo repository.ClienteRepository
}

func NewClienteService(clienteRepo repository.ClienteRepository) *ClienteService {
	return &ClienteService{clienteRepo: clienteRepo}
}

func (s *ClienteService) CreateCliente(ctx context.Context, name, email, password string) (model.Cliente, error) {
	existingCliente, err := s.clienteRepo.FindByEmail(ctx, email)

	if err != nil && !errors.Is(err, model.ErrClienteNaoEncontrado) {
		return model.Cliente{}, err
	}

	if existingCliente.ID != uuid.Nil {
		return model.Cliente{}, model.ErrEmailCadastrado
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.Cliente{}, err
	}

	cliente := model.Cliente{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := s.clienteRepo.Save(ctx, cliente); err != nil {
		return model.Cliente{}, err
	}

	return cliente, nil
}
