#!/bin/bash

START_BOOK_HPORT=11000
START_BOOK_RPORT=12000
START_USER_HPORT=13000
START_USER_RPORT=14000

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

# Garantir que os diretórios para leveldb existam
mkdir -p "$BASE_USER_PATH/$REPLICA_ID"
mkdir -p "$BASE_BOOK_PATH/$REPLICA_ID"


# Instanciar uma réplica inicial:
# ./bin/<database> -id node1 -haddr localhost:11000 -raddr localhost:12000 ~/node
# Adicionar réplicas: 
# ./bin/<database> -id node1 -haddr localhost:11001 -raddr localhost:12001 -join :11000 ~/node1
#Iniciar a réplica apropriada com base no cluster_id
if [ $CLUSTER_ID -eq 0 ]; then
    echo "Iniciando Réplica do Cluster de Usuários $REPLICA_ID em $BASE_USER_PATH/$REPLICA_ID"
    ./bin/user-database -id $REPLICA_ID -haddr localhost:$REPLICA_HPORT -raddr localhost:$REPLICA_RPORT
else
    echo "Iniciando Réplica do Cluster de Livros $REPLICA_ID em $BASE_BOOK_PATH/$REPLICA_ID"
    REPLICA_HPORT=$((START_BOOK_HPORT + REPLICA_ID)) # Porta que réplica atenderá
    REPLICA_RPORT=$((START_BOOK_RPORT + REPLICA_ID)) # Porta para comunicação entre réplicas
    # Garantir que os diretórios para salvar store existam
    # TODO: Verificar se é necessário criar arquivo.
    $REPLICA_DIR="/tmp/grupo9/store/$REPLICA_HPORT"
    mkdir -p "$REPLICA_DIR"

    if [ $REPLICA_ID -eq 0 ]; then
        ./bin/book-database -id $REPLICA_ID -haddr localhost:$REPLICA_HPORT -raddr localhost:$REPLICA_RPORT
    else
        ./bin/book-database -id $REPLICA_ID \
            -haddr localhost:$REPLICA_HPORT \
            -raddr localhost:$REPLICA_RPORT \
            -join :$START_BOOK_RPORT \
            $REPLICA_DIR
    fi
 fi
