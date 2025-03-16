package services

import (
	"time"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/domain/repositories"

	"github.com/google/uuid"
)

type RestaurantService struct {
	restaurantRepo repositories.RestaurantRepository
}

func NewRestaurantService(restaurantRepo repositories.RestaurantRepository) *RestaurantService {
	return &RestaurantService{
		restaurantRepo: restaurantRepo,
	}
}

func (s *RestaurantService) Create(restaurant *models.Restaurant) error {
	// Definir valores padrão
	if restaurant.Status == "" {
		restaurant.Status = models.SubscriptionStatusTrial
	}

	// Se for um teste gratuito, definir a data de término
	if restaurant.Status == models.SubscriptionStatusTrial {
		trialEnd := time.Now().AddDate(0, 1, 0) // 1 mês de teste
		restaurant.TrialEndsAt = &trialEnd
	}

	return s.restaurantRepo.Create(restaurant)
}

func (s *RestaurantService) GetByID(id uuid.UUID) (*models.Restaurant, error) {
	return s.restaurantRepo.FindByID(id)
}

func (s *RestaurantService) Update(restaurant *models.Restaurant) error {
	return s.restaurantRepo.Update(restaurant)
}

func (s *RestaurantService) Delete(id uuid.UUID) error {
	return s.restaurantRepo.Delete(id)
}

func (s *RestaurantService) List() ([]models.Restaurant, error) {
	return s.restaurantRepo.List()
}

func (s *RestaurantService) GetByStatus(status models.SubscriptionStatus) ([]models.Restaurant, error) {
	return s.restaurantRepo.FindByStatus(status)
}

func (s *RestaurantService) SearchByName(name string) ([]models.Restaurant, error) {
	return s.restaurantRepo.FindByName(name)
}

func (s *RestaurantService) UpdateStatus(id uuid.UUID, status models.SubscriptionStatus) error {
	return s.restaurantRepo.UpdateStatus(id, status)
}

// Atualiza restaurantes com testes expirados para status inativo
func (s *RestaurantService) UpdateExpiredTrials() error {
	restaurants, err := s.restaurantRepo.FindByStatus(models.SubscriptionStatusTrial)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, restaurant := range restaurants {
		if restaurant.TrialEndsAt != nil && now.After(*restaurant.TrialEndsAt) {
			s.restaurantRepo.UpdateStatus(restaurant.ID, models.SubscriptionStatusInactive)
		}
	}

	return nil
}

// Verifica se um restaurante está ativo ou em período de teste
func (s *RestaurantService) IsRestaurantActive(id uuid.UUID) (bool, error) {
	restaurant, err := s.restaurantRepo.FindByID(id)
	if err != nil {
		return false, err
	}

	if restaurant.Status == models.SubscriptionStatusActive {
		return true, nil
	}

	if restaurant.Status == models.SubscriptionStatusTrial {
		now := time.Now()
		if restaurant.TrialEndsAt != nil && now.Before(*restaurant.TrialEndsAt) {
			return true, nil
		}
	}

	return false, nil
}
