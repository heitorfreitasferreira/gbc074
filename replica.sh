#!/bin/bash

# Verifique se os parâmetros necessários foram fornecidos
if [ $# -ne 2 ]; then
    echo "Uso: $0 <replica_id> <cluster_id>"
    echo "  replica_id: 0, 1 ou 2"
    echo "  cluster_id: 0 (usuários) ou 1 (livros)"
    exit 1
fi

REPLICA_ID=$1
CLUSTER_ID=$2

# Validar replica_id
if [[ ! $REPLICA_ID =~ ^[0-2]$ ]]; then
    echo "Erro: replica_id deve ser 0, 1 ou 2"
    exit 1
fi

# Validar cluster_id
if [[ ! $CLUSTER_ID =~ ^[0-1]$ ]]; then
    echo "Erro: cluster_id deve ser 0 (usuários) ou 1 (livros)"
    exit 1
fi

BASE_USER_PATH="/tmp/grupo9/user"
BASE_BOOK_PATH="/tmp/grupo9/book"

# Garantir que os diretórios existam
mkdir -p "$BASE_USER_PATH/$REPLICA_ID"
mkdir -p "$BASE_BOOK_PATH/$REPLICA_ID"

# Iniciar a réplica apropriada com base no cluster_id
# if [ $CLUSTER_ID -eq 0 ]; then
#     echo "Iniciando Réplica do Cluster de Usuários $REPLICA_ID em $BASE_USER_PATH/$REPLICA_ID"
#     go run cmd/replica/main.go -cluster 0 -id $REPLICA_ID -path "$BASE_USER_PATH/$REPLICA_ID"
# else
#     echo "Iniciando Réplica do Cluster de Livros $REPLICA_ID em $BASE_BOOK_PATH/$REPLICA_ID"
#     go run cmd/replica/main.go -cluster 1 -id $REPLICA_ID -path "$BASE_BOOK_PATH/$REPLICA_ID"
# fi