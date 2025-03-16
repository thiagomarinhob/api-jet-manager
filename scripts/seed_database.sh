#!/bin/bash
# Script para popular o banco de dados com dados de teste

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}==============================================${NC}"
echo -e "${YELLOW}  SEEDER PARA O SISTEMA DE RESTAURANTES SAAS  ${NC}"
echo -e "${YELLOW}==============================================${NC}"

# Verificar se o banco de dados está operacional
echo -e "${YELLOW}Verificando conexão com o banco de dados...${NC}"

# Carregar variáveis do arquivo .env se existir
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Configurações padrão do banco
DB_HOST="localhost"
DB_PORT=${BLUEPRINT_DB_PORT:-"5432"}
DB_USER=${BLUEPRINT_DB_USERNAME:-"docker"}
DB_PASSWORD=${BLUEPRINT_DB_PASSWORD:-"docker"}
DB_NAME=${BLUEPRINT_DB_DATABASE:-"jetmanager"}

# Testar conexão com PostgreSQL
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1" &>/dev/null
if [ $? -ne 0 ]; then
    echo -e "${RED}Não foi possível conectar ao banco de dados. Verifique se o PostgreSQL está em execução e se as credenciais estão corretas.${NC}"
    exit 1
fi

echo -e "${GREEN}Conexão com o banco de dados estabelecida!${NC}"

# Verificar se as tabelas já existem
echo -e "${YELLOW}Verificando se o esquema do banco de dados está criado...${NC}"
TABLES=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")

if [ -z "$TABLES" ]; then
    echo -e "${YELLOW}Banco de dados vazio. Executando migrações primeiro...${NC}"
    
    # Verificar se o arquivo de migração existe
    if [ ! -f "migrations/01_initial_schema.up.sql" ]; then
        echo -e "${RED}Arquivo de migração não encontrado. Verifique se o arquivo migrations/01_initial_schema.up.sql existe.${NC}"
        exit 1
    fi
    
    # Executar migração
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/01_initial_schema.up.sql
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}Erro ao executar as migrações.${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Migrações executadas com sucesso!${NC}"
else
    echo -e "${GREEN}Esquema do banco de dados já existe!${NC}"
fi

# Confirmar com o usuário
echo -e "${YELLOW}ATENÇÃO: Esta operação irá popular o banco de dados com dados fictícios, mantendo apenas o superadmin existente.${NC}"
echo -e "${YELLOW}Deseja continuar? [S/n]${NC}"
read -r response
response=${response:-S}

if [[ ! "$response" =~ ^[Ss]$ ]]; then
    echo -e "${YELLOW}Operação cancelada pelo usuário.${NC}"
    exit 0
fi

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Go não está instalado. Tentando executar com docker...${NC}"
    
    # Verificar se o Docker está instalado
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Docker não está instalado. Por favor, instale Go ou Docker para executar este script.${NC}"
        exit 1
    fi
    
    # Executar com Docker
    echo -e "${YELLOW}Executando o seeder com Docker...${NC}"
    docker run --rm -v "$(pwd):/app" -w /app \
        --network=api-jet-manager \
        -e DB_HOST=localhost \
        -e DB_PORT=$BLUEPRINT_DB_PORT \
        -e DB_USER=$BLUEPRINT_DB_USERNAME \
        -e DB_PASSWORD=$BLUEPRINT_DB_PASSWORD \
        -e DB_NAME=$BLUEPRINT_DB_DATABASE \
        golang:1.23 go run scripts/seed.go
else
    # Executar com Go
    echo -e "${YELLOW}Executando o seeder com Go local...${NC}"
    go run scripts/seed.go
fi

if [ $? -eq 0 ]; then
    echo -e "${GREEN}=======================================================${NC}"
    echo -e "${GREEN}  Banco de dados populado com sucesso!                 ${NC}"
    echo -e "${GREEN}=======================================================${NC}"
    echo -e "${YELLOW}Informações de acesso:${NC}"
    echo -e "${YELLOW}  * SuperAdmin: superadmin@example.com / superadmin   ${NC}"
    echo -e "${YELLOW}  * Admin: admin@X.example.com / admin123             ${NC}"
    echo -e "${YELLOW}    (onde X é o ID do restaurante)                    ${NC}"
    echo -e "${YELLOW}  * Manager: managerY@X.example.com / manager123      ${NC}"
    echo -e "${YELLOW}    (onde X é o ID do restaurante e Y é o número do gerente)${NC}"
    echo -e "${YELLOW}  * Staff: staffY@X.example.com / staff123            ${NC}"
    echo -e "${YELLOW}    (onde X é o ID do restaurante e Y é o número do funcionário)${NC}"
    echo -e "${GREEN}=======================================================${NC}"
else
    echo -e "${RED}Erro ao executar o seeder.${NC}"
    exit 1
fi