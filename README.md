# WinRate

Autor: Diogo Lourenço Andrade
Descrição: WinRate é um sistema de acompanhamento de partidas de jogos de cartas, que permite registrar e analisar as vitórias e derrotas de cada jogador, além de monitorar o desempenho de diferentes decks em diversas partidas. O sistema fornece uma API RESTful para interagir com os dados, com funcionalidades de visualização de partidas, cálculos de taxa de vitória e gestão de decks.

Funcionalidades
Acompanhamento de partidas: Registra informações sobre cada partida, como os decks usados, vitórias e derrotas.
Cálculo de taxa de vitória: Exibe a taxa de vitória de cada jogador e deck.
API RESTful: Interface para integração com outros sistemas e aplicações, permitindo a consulta, inserção e modificação de dados de partidas e decks.
Suporte para múltiplos jogadores e decks: O sistema permite que vários jogadores registrem e acompanhem suas partidas com diferentes decks.
Tecnologias Utilizadas
Go (Golang): Linguagem principal para o desenvolvimento do backend.
MySQL: Banco de dados relacional para armazenar informações de partidas e decks.
Gorilla Mux: Pacote para roteamento HTTP em Go.
JSON: Formato de dados utilizado para comunicação entre o cliente e o servidor.
Estrutura do Projeto
/cmd: Contém o código principal da aplicação.
/models: Contém os modelos de dados, como Match e Deck.
/handlers: Contém a lógica das rotas e manipulação de requisições.
/db: Conexões e consultas ao banco de dados MySQL.
/migrations: Scripts para criação e atualização do banco de dados.
Endpoints
GET /matches/{id}
Retorna todas as partidas de um jogador específico, com detalhes sobre o deck usado, vitórias, derrotas e data da criação.

Parâmetro:

id (int): ID do jogador.
Resposta:

Status: 200 OK
Corpo:
json
Copiar código
[
{
"matchsID": 1,
"user_deck_name": "Deck A",
"opponent_deck_name": "Deck B",
"victories": 3,
"defeats": 2,
"created_at": "2024-12-03 10:30:00"
}
]
POST /matches
Cria uma nova partida com as informações fornecidas.

Corpo:

json
Copiar código
{
"user_deck_name": "Deck A",
"opponent_deck_name": "Deck B",
"victories": 3,
"defeats": 2
}
Resposta:

Status: 201 Created
Corpo:
json
Copiar código
{
"matchsID": 1,
"user_deck_name": "Deck A",
"opponent_deck_name": "Deck B",
"victories": 3,
"defeats": 2,
"created_at": "2024-12-03 10:30:00"
}
Como Rodar o Projeto
Pré-requisitos
Go (Golang): Certifique-se de ter o Go instalado na sua máquina.
MySQL: Instale e configure o MySQL para armazenar os dados.
Configuração do Banco de Dados
Crie um banco de dados MySQL com o nome WinRate.
Aplique as migrações necessárias para criar as tabelas de partidas e decks.
Rodando o Projeto
Clone o repositório:

bash
Copiar código
git clone https://github.com/diogolouro/WinRate.git
cd WinRate
Instale as dependências:

bash
Copiar código
go mod tidy
Execute a aplicação:

bash
Copiar código
go run main.go
A aplicação estará rodando na URL http://localhost:8080.

Contribuições
Sinta-se à vontade para contribuir com melhorias ou novos recursos. Para isso, siga os seguintes passos:

Fork este repositório.
Crie uma branch para a sua feature (git checkout -b minha-feature).
Commit suas alterações (git commit -am 'Adiciona nova feature').
Push para a branch (git push origin minha-feature).
Abra um Pull Request.
Licença
Este projeto está licenciado sob a Licença MIT - consulte o arquivo LICENSE para mais detalhes.
