# Entrega 2

## 1. Subir dois clusters de banco de dados

1. 1 cluster para User e UserBook OK
1. cluster para Book. OK
1. Cada cluster deve possuir três réplicas. OK
1. O banco de dados deve ser [LevelDB](https://github.com/syndtr/goleveldb) OK
1. Adaptar a API Rest do Raft para seguir o padrão definido OK

## 2. Utilizar Raft para comunicação entre servidor e BD

1. Utilizar pacote [Raft para GO](https://github.com/hashicorp/raft). OK
1. Remover uso do MQTT OK
1. Colocar a chamada http para o servidor raft no handler grpc (onde tinha a chamada pro MQTT) Analisar bib-server 
1. Conferir se é realmente é para usar HTTP na comunicação cad-server/bib-server --> user-database/book-database OK

## 3. Implementar serviço de cache de dados

1. Validade de 5 segundos OK
1. Atualizada a cada modificação/consulta nos clusters. OK

## 4. Atualizar scripts e makefile

1. Implementar script replica.sh para instanciar replicas (0, 1 ou 2) e adicioná-las em um cluster (0 ou 1). OK
1. Atualizar scripts de execução dos servidores se necessário. OK
1. Atualizar makefile para instalar novas dependências necessárias. OK
