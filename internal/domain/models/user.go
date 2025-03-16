package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserType string

const (
	UserTypeSuperAdmin UserType = "superadmin" // Administrador da plataforma SaaS
	UserTypeAdmin      UserType = "admin"      // Administrador de um restaurante
	UserTypeManager    UserType = "manager"    // Gerente de um restaurante
	UserTypeStaff      UserType = "staff"      // Funcionário de um restaurante
)

type User struct {
	ID           uuid.UUID   `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name         string      `gorm:"size:100;not null" json:"name"`
	Email        string      `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password     string      `gorm:"size:100;not null" json:"-"`
	Type         UserType    `gorm:"size:20;not null;default:'staff'" json:"type"` // superadmin, admin, manager, staff
	RestaurantID *uuid.UUID  `json:"restaurant_id" gorm:"type:uuid"`
	Restaurant   *Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// BeforeSave - Hook para hashear a senha antes de salvar
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password == "" {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifica se a senha corresponde ao hash
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// IsSuperAdmin verifica se o usuário é superadmin
func (u *User) IsSuperAdmin() bool {
	return u.Type == UserTypeSuperAdmin
}

// CanAccessRestaurant verifica se o usuário pode acessar o restaurante especificado
func (u *User) CanAccessRestaurant(restaurantID uuid.UUID) bool {
	// Superadmin pode acessar qualquer restaurante
	if u.IsSuperAdmin() {
		return true
	}

	// Outros usuários só podem acessar seu próprio restaurante
	return u.RestaurantID == nil && *u.RestaurantID == restaurantID
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return nil
}
