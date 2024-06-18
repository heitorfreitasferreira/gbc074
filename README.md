# rpc-mqtt-library-manager

Repositório com a implementação do trabalho proposto para a matéria de sistemas digitais na FACOM-UFU

A descrição do projeto como copiada [do repositório do professor](https://github.com/paulo-coelho/ds_notes/blob/main/docs/projeto.md) no dia 03/06 está em [aqui](./descricao.md)

## Instalar dependencias

As dependencias necessárias parar compilar o projeto são:

- [protoc](https://developers.google.com/protocol-buffers)
- [go](https://golang.org/)

As dependencias necessárias para rodar o projeto são:

- [mosquitto](https://mosquitto.org/)

Caso tenha o gerenciador de pacotes yay instalado, basta executar o script `install-deps.sh`:

```bash
chmod +x install-deps.sh
./install-deps.sh
```

## Compilar os serviços

```bash
chmod +x compile.sh
./compile.sh
```

Os executáveis estarão na pasta `bin/`

## Rodar os serviços

### Servidores

Para subir os servidores de cadastro e da biblioteca, basta executar o binário diretamente, ou via script (recomendado), passando a porta como argumento:

```bash
chmod +x cad-server.sh
chmod +x bibi-server.sh

./cad-server.sh 42069
./bibi-server.sh 6666
```

É possivel subir varias versões dos servidores, desde que sejam em portas diferentes.

Caso execute o binário diretamente, o argumento deve ser normado como `--port=42069`, e é possível passar o argumento `--host=localhost` para mudar o host do servidor.

### Clientes

Para rodar os clientes, basta executar o binário diretamente, ou via script (recomendado), passando o host e a porta como argumento:

```bash
chmod +x cad-client.sh
chmod +x bibi-client.sh

./cad-client.sh localhost 42069
./bibi-client.sh localhost 6666
```

Isso irá conectar o cliente ao servidor na porta especificada.

Tamém é possível executar o binário diretamente, passando os argumentos `--host=localhost` e `--port=42069`.
