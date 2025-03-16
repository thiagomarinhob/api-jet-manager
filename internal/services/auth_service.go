// internal/services/auth_service.go
package services

import (
	"errors"
	"fmt"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/domain/repositories"
	"api-jet-manager/internal/infrastructure/auth"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo   repositories.UserRepository
	jwtService *auth.JWTService
}

func NewAuthService(userRepo repositories.UserRepository, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(user *models.User) error {
	// Verificar se já existe um usuário com o mesmo email
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Criar o usuário
	return s.userRepo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	// Buscar usuário pelo email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Verificar a senha
	if err := user.CheckPassword(password); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Gerar o token JWT
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (s *AuthService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *AuthService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *AuthService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

func (s *AuthService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}

// Verifica se já existe um superadmin no sistema
func (s *AuthService) SuperAdminExists() (bool, error) {
	users, err := s.userRepo.FindByType(models.UserTypeSuperAdmin)
	if err != nil {
		return false, err
	}
	return len(users) > 0, nil
}

// Lista usuários por restaurante
func (s *AuthService) ListUsersByRestaurant(restaurantID uuid.UUID) ([]models.User, error) {
	return s.userRepo.FindByRestaurant(restaurantID)
}

// Lista usuários por tipo
func (s *AuthService) ListUsersByType(userType models.UserType) ([]models.User, error) {
	return s.userRepo.FindByType(userType)
}
