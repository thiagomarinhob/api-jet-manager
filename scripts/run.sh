#!/bin/bash
# Script para construir e executar o projeto com Docker

# Verifica a estrutura do projeto
./scripts/verify.sh

# Define permissões para os scripts
chmod +x scripts/*.sh

# Para e remove containers existentes
echo "Parando e removendo containers existentes..."
docker compose down -v

# Constrói e inicia os containers
echo "Construindo e iniciando os containers..."
docker compose up --build

# Outras opções úteis:
# docker-compose up --build -d  # Para executar em background
# docker-compose logs -f        # Para acompanhar os logs quando em background