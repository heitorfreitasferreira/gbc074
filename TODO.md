# Entrega 2

Item | Responsável
--- | ---
1.1 | [ ] @
1.2 | [ ] @
1.3 | [ ] @
1.4 | [ ] @
2.1 | [ ] @
2.2 | [ ] @
2.3 | [ ] @
3.1 | [ ] @
3.2 | [ ] @
4.1 | [ ] @heitorfreitasferreira
4.2 | [ ] @heitorfreitasferreira
4.3 | [ ] @heitorfreitasferreira

## 1. Subir dois clusters de banco de dados

1. 1 cluster para User e UserBook
1. 1 cluster para Book.
1. Cada cluster deve possuir três réplicas.
1. O banco de dados deve ser [LevelDB](https://github.com/syndtr/goleveldb)

## 2. Utilizar Raft para comunicação entre servidor e BD

1. Utilizar pacote [Raft para GO](https://github.com/hashicorp/raft).
1. Remover uso do MQTT
1. Nada muda entre cliente e servidor (manter gRPC).

## 3. Implementar serviço de cache de dados

1. Validade de 5 segundos
1. Atualizada a cada modificação/consulta nos clusters.

## 4. Atualizar scripts e makefile

1. Implementar script replica.sh para instanciar replicas (0, 1 ou 2) e adicioná-las em um cluster (0 ou 1).
1. Atualizar scripts de execução dos servidores se necessário.
1. Atualizar makefile para instalar novas dependências necessárias.
