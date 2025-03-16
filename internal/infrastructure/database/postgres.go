// internal/infrastructure/database/postgres.go
package database

import (
	"fmt"
	"log"
	"time"

	"api-jet-manager/internal/config"
	"api-jet-manager/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	*gorm.DB
}

func NewPostgresConnection(cfg *config.Config) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.BLUEPRINT_DB_HOST, cfg.BLUEPRINT_DB_PORT, cfg.BLUEPRINT_DB_USERNAME, cfg.BLUEPRINT_DB_PASSWORD, cfg.BLUEPRINT_DB_DATABASE, cfg.DBSSLMode)

	// Tenta conectar com retry
	var db *gorm.DB
	var err error
	maxRetries := 10
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Tentando conectar ao banco de dados (tentativa %d/%d)...", i+1, maxRetries)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			break
		}

		log.Printf("Falha na conexão: %v. Tentando novamente em %v...", err, retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados após %d tentativas: %w", maxRetries, err)
	}

	// Verifica se o banco está acessível
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter conexão SQL: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao fazer ping no banco de dados: %w", err)
	}

	log.Println("Conectado com sucesso ao banco de dados PostgreSQL")
	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func RunMigrations(db *PostgresDB) error {
	log.Println("Running database migrations...")

	// Auto-migrate tabelas
	return db.AutoMigrate(
		&models.User{},
		&models.Table{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Restaurant{},
		&models.FinancialTransaction{},
	)
}
