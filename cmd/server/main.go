// cmd/server/main.go
package main

import (
	"log"
	"time"

	"api-jet-manager/internal/api/routes"
	"api-jet-manager/internal/config"
	"api-jet-manager/internal/infrastructure/database"
)

func main() {
	log.Println("Iniciando o sistema de gerenciamento de restaurante...")

	// Carrega configurações
	log.Println("Carregando configurações...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Falha ao carregar configurações: %v", err)
	}
	log.Printf("Configurações carregadas com sucesso. Host do BD: %s, Modo Gin: %s", cfg.BLUEPRINT_DB_HOST, cfg.GinMode)

	// Aguardar alguns segundos para garantir que o banco de dados esteja pronto
	log.Println("Aguardando o banco de dados ficar disponível...")
	time.Sleep(5 * time.Second)

	// Conecta ao banco de dados
	log.Println("Conectando ao banco de dados...")
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	defer func() {
		log.Println("Fechando conexão com o banco de dados...")
		if err := db.Close(); err != nil {
			log.Printf("Erro ao fechar conexão com o banco de dados: %v", err)
		}
	}()
	log.Println("Conexão com o banco de dados estabelecida com sucesso")

	// Executa migrações
	log.Println("Executando migrações do banco de dados...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Falha ao executar migrações: %v", err)
	}
	log.Println("Migrações executadas com sucesso")

	// Configuração do router e inicialização do servidor
	log.Println("Configurando rotas da API...")
	router := routes.SetupRouter(cfg, db)

	log.Printf("Iniciando servidor na porta %s...", cfg.ServerAddress)
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
