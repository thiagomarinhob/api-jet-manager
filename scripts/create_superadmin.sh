#!/bin/bash
# Script para criar o usuário administrador inicial

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Criando usuário super admin...${NC}"

# Executa o script Go
go run scripts/create_superadmin.go

echo -e "${GREEN}Script finalizado!${NC}"
