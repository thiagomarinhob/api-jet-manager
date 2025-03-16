package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// UserType representa o tipo do usuário
type UserType string

const (
	UserTypeSuperAdmin UserType = "superadmin"
	UserTypeAdmin      UserType = "admin"
	UserTypeManager    UserType = "manager"
	UserTypeStaff      UserType = "staff"
)

// User representa o modelo de usuário
type User struct {
	ID           string   `gorm:"primaryKey"`
	Name         string   `gorm:"size:100;not null"`
	Email        string   `gorm:"size:100;uniqueIndex;not null"`
	Password     string   `gorm:"size:100;not null"`
	Type         UserType `gorm:"size:20;not null;default:'staff'"`
	RestaurantID *string  `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	// Obtem os dados do banco de dados
	dbHost := "localhost"
	dbPort := getEnv("BLUEPRINT_DB_PORT", "5432")
	dbUser := getEnv("BLUEPRINT_DB_USERNAME", "docker")
	dbPassword := getEnv("BLUEPRINT_DB_PASSWORD", "docker")
	dbName := getEnv("BLUEPRINT_DB_DATABASE", "jetmanager")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Constrói a string de conexão
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Conecta ao banco de dados
	log.Println("Conectando ao banco de dados...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	log.Println("Conexão estabelecida com sucesso!")

	// Verifica se a tabela users existe
	if !db.Migrator().HasTable(&User{}) {
		log.Println("Tabela 'users' não existe. Criando...")
		if err := db.AutoMigrate(&User{}); err != nil {
			log.Fatalf("Erro ao criar tabela users: %v", err)
		}
	}

	// Verifica se já existe algum superadmin
	var count int64
	db.Model(&User{}).Where("type = ?", UserTypeSuperAdmin).Count(&count)
	if count > 0 {
		log.Println("Já existe pelo menos um usuário superadmin no sistema.")
		log.Println("Se você esqueceu a senha, use o processo de recuperação de senha ou crie outro usuário.")
		return
	}

	// Dados do novo superadmin
	adminName := "Super Admin"
	adminEmail := "superadmin@example.com"
	adminPassword := "superadmin123" // Você deve alterar isso!

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erro ao gerar hash da senha: %v", err)
	}

	// Cria o usuário superadmin
	admin := User{
		ID:           uuid.New().String(),
		Name:         adminName,
		Email:        adminEmail,
		Password:     string(hashedPassword),
		Type:         UserTypeSuperAdmin,
		RestaurantID: nil,
	}

	// Insere no banco de dados
	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("Erro ao criar usuário superadmin: %v", err)
	}

	log.Println("===================================================")
	log.Println("Usuário superadmin criado com sucesso!")
	log.Println("Nome: ", adminName)
	log.Println("Email: ", adminEmail)
	log.Println("Senha: ", adminPassword, " (Altere-a após o primeiro login!)")
	log.Println("===================================================")
}

// getEnv obtém variável de ambiente ou retorna valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
