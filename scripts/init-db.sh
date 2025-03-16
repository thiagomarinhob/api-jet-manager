#!/bin/bash
set -e

# Script para inicializar o banco de dados

psql -v ON_ERROR_STOP=1 --username "$BLUEPRINT_DB_USERNAME" --dbname "$BLUEPRINT_DB_DATABASE" <<-EOSQL
    CREATE DATABASE jetmanager;
    GRANT ALL PRIVILEGES ON DATABASE jetmanager TO postgres;
EOSQL

echo "Banco de dados 'jetmanager' criado com sucesso!"