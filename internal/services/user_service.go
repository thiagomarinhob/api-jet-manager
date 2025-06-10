package services

import (
	"errors"
	"fmt"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/domain/repositories"
	"api-jet-manager/internal/infrastructure/auth"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo   repositories.UserRepository
	jwtService *auth.JWTService
}

func NewUserService(userRepo repositories.UserRepository, jwtService *auth.JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *UserService) Register(user *models.User) error {
	// Verificar se já existe um usuário com o mesmo email
	existingUser, err := s.userRepo.FindByEmail(*user.RestaurantID, user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Criar o usuário
	return s.userRepo.Create(user)
}

func (s *UserService) Login(email, password string) (string, *models.User, error) {
	// Buscar usuário pelo email
	user, err := s.userRepo.FindByEmailGlobal(email)
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

func (s *UserService) FindUserByID(restaurantID uuid.UUID, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(restaurantID, id)
}

func (s *UserService) FindUserByEmail(restaurantID uuid.UUID, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(restaurantID, email)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(restaurantID uuid.UUID, id uuid.UUID) error {
	return s.userRepo.Delete(restaurantID, id)
}

func (s *UserService) ListUsers(restaurantID uuid.UUID) ([]models.User, error) {
	return s.userRepo.List(restaurantID)
}

// Verifica se já existe um superadmin no sistema
func (s *UserService) SuperAdminExists() (bool, error) {
	users, err := s.userRepo.FindByTypeGlobal(models.UserTypeSuperAdmin)
	if err != nil {
		return false, err
	}
	return len(users) > 0, nil
}

// Lista usuários por restaurante
func (s *UserService) ListUsersByRestaurant(restaurantID uuid.UUID) ([]models.User, error) {
	return s.userRepo.FindByRestaurant(restaurantID)
}

// Lista usuários por tipo
func (s *UserService) ListUsersByType(restaurantID uuid.UUID, userType models.UserType) ([]models.User, error) {
	return s.userRepo.FindByType(restaurantID, userType)
}
