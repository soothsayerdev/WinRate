let deckNames = {}; // Storage to map name of deck by name

document.addEventListener('DOMContentLoaded', function () {
    const userID = localStorage.getItem('userID');
    if (!userID) {
        alert('User not logged in');
        window.location.href = 'frontend/index.html';
        return;
    }

    // Carregar os decks disponíveis
    fetch('http://localhost:8080/decks')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            data.forEach(deck => {
                deckNames[deck.DeckName] = deck.ID; // Mapeia pelo nome do deck
            });
            console.log("Deck names loaded: ", deckNames);
            loadMatches(); // Carregar as partidas após carregar os decks
        })
        .catch(error => {
            console.error('Error fetching decks:', error);
        });

    // Carregar a tabela de partidas do usuário
    async function loadMatches() {
        try {
            const response = await fetch('/matches'); // Endpoint para buscar as partidas
            if (!response.ok) {
                throw new Error(`Erro ao carregar as partidas: ${response.status}`);
            }
            const matches = await response.json();

            const tableBody = document.getElementById('deckStatsTable').getElementsByTagName('tbody')[0];
            tableBody.innerHTML = ''; // Limpar a tabela antes de renderizar

            matches.forEach(match => {
                const row = `
                    <tr>
                        <td>${match.user_deck_name}</td>
                        <td>${match.opponent_deck_name}</td>
                        <td>${match.victories}</td>
                        <td>${match.defeats}</td>
                        <td>${((match.victories / (match.victories + match.defeats)) * 100).toFixed(2)}%</td>
                    </tr>
                `;
                tableBody.innerHTML += row;
            });
        } catch (error) {
            console.error('Erro ao carregar as partidas:', error);
        }
    }

    // Função para criar uma nova partida
    async function createMatch() {
        const userDeckName = document.getElementById('userDeckName').value.trim();
        const opponentDeckName = document.getElementById('opponentDeckName').value.trim();
        const victories = parseInt(document.getElementById('victories').value.trim(), 10);
        const defeats = parseInt(document.getElementById('defeats').value.trim(), 10);

        if (!userDeckName || !opponentDeckName || isNaN(victories) || isNaN(defeats)) {
            alert('Por favor, preencha todos os campos corretamente.');
            return;
        }

        // Verificar se os nomes dos decks são válidos
        if (!deckNames[userDeckName] || !deckNames[opponentDeckName]) {
            alert('Um ou mais nomes de deck não são válidos.');
            return;
        }

        const payload = {
            user_deck_name: userDeckName,
            opponent_deck_name: opponentDeckName,
            victories,
            defeats,
        };

        try {
            const response = await fetch('/create-match', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });

            if (response.ok) {
                alert('Match criada com sucesso!');
                loadMatches(); // Recarregar a tabela de partidas
            } else {
                const errorText = await response.text();
                console.error('Erro ao criar a partida:', errorText);
                alert('Erro ao criar a partida. Tente novamente.');
            }
        } catch (error) {
            console.error('Erro ao criar a partida:', error);
        }
    }

    // Associar a função ao botão de criação de partidas
    document.getElementById('createMatchButton').addEventListener('click', createMatch);

    // Criar novo deck
    document.getElementById('deckForm').addEventListener('submit', function (e) {
        e.preventDefault();

        const deckName = document.getElementById('deckName').value.trim();
        const userIdInt = parseInt(userID, 10);

        if (!deckName) {
            alert('O nome do deck não pode estar vazio.');
            return;
        }

        fetch('http://localhost:8080/decks', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                deck_name: deckName,
                user_id: userIdInt,
            }),
        })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(text => {
                        console.log('Error response:', text);
                        throw new Error(text);
                    });
                }
                return response.json();
            })
            .then(data => {
                alert('Deck criado com sucesso!');
                deckNames[data.deck_name] = data.deckId; // Adiciona ao mapeamento
            })
            .catch(error => {
                console.error('Erro ao criar o deck:', error);
                alert('Erro ao criar o deck. Tente novamente.');
            });
    });
});
