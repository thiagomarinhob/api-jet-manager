#!/bin/bash
# Script para verificar a estrutura do projeto

echo "Verificando estrutura de diretórios do projeto..."

# Verifica se os diretórios principais existem
echo "Verificando diretórios principais:"
directories=(
  "cmd/server"
  "internal/api/handlers"
  "internal/api/middlewares"
  "internal/api/routes"
  "internal/config"
  "internal/domain/models"
  "internal/domain/repositories"
  "internal/infrastructure/auth"
  "internal/infrastructure/database"
  "internal/infrastructure/repositories"
  "internal/services"
  "migrations"
  "scripts"
)

for dir in "${directories[@]}"; do
  if [ -d "$dir" ]; then
    echo "✓ Diretório $dir existe"
  else
    echo "✗ Diretório $dir não existe - criando..."
    mkdir -p "$dir"
  fi
done

# Verifica se os arquivos principais existem
echo -e "\nVerificando arquivos principais:"
files=(
  "cmd/server/main.go"
  "docker-compose.yml"
  "Dockerfile"
  "go.mod"
  "go.sum"
  ".env.example"
  "scripts/init-db.sh"
)

for file in "${files[@]}"; do
  if [ -f "$file" ]; then
    echo "✓ Arquivo $file existe"
  else
    echo "✗ Arquivo $file não existe"
  fi
done

# Verifica permissões de execução dos scripts
echo -e "\nVerificando permissões de execução:"
scripts=(
  "scripts/init-db.sh"
  "scripts/build.sh"
  "scripts/verify.sh"
)

for script in "${scripts[@]}"; do
  if [ -f "$script" ]; then
    if [ -x "$script" ]; then
      echo "✓ Script $script tem permissão de execução"
    else
      echo "✗ Script $script não tem permissão de execução - configurando..."
      chmod +x "$script"
    fi
  fi
done

echo -e "\nVerificação concluída!"