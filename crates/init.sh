#!/bin/bash

echo "Digite a URI do MongoDB:"
read -r MONGO_URI

echo "Digite o nome do banco de dados:"
read -r DATABASE_NAME

echo "Digite o nome da coleção:"
read -r COLLECTION_NAME

# Criar ou sobrescrever o arquivo .env
{
    echo "export MONGO_URI=\"$MONGO_URI\""
    echo "export DATABASE_NAME=\"$DATABASE_NAME\""
    echo "export COLLECTION_NAME=\"$COLLECTION_NAME\""
} > .env

echo ".env configurado com sucesso!"

