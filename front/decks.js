let deckNames = {}; // Storage to map name of deck by ID

document.addEventListener('DOMContentLoaded', function() {
    const userID = localStorage.getItem('userID');
    if(!userID) {
        alert('User not logged in');
        window.location.href = '/index.html';
        return;
    }
    
    // Handle deck creation
    document.getElementById('deckForm').addEventListener('submit', function(e) {
        e.preventDefault();

        const deckName = document.getElementById('deckName').value;

        // Convert userID to integer
        const userIdInt = parseInt(userID, 10); // Base 10
        
        // Check if userIdInt is a valid number
        if (isNaN(userIdInt)) {
            alert('Invalid user ID');
            return;
        }

        fetch('http://localhost:8080/decks', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                deck_name: deckName,
                user_id: userIdInt, // Assuming user_id is stored in localStorage
            }),
        })
        .then(response => {
            if(!response.ok) {
                return response.text().then(text => { 
                    console.log('Error response:', text);
                    throw new Error(text);
                });
                
            }
            return response.json();
        })
        .then(data => {
            alert('Deck created successfuly!');
            // Optionally, refresh the table with new deck data
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to create deck: ' + error.message);
        });
        console.log('Deck name:', deckName);
        console.log('User ID:', userID);
    });

    // Handle WinRate updates
    document.getElementById('winrateForm').addEventListener('submit', function(e){
        e.preventDefault();

        const userDeck = document.getElementById('userDeck').value;
        const opponentDeck = document.getElementById('opponentDeck').value;
        const wins = document.getElementById('wins').value;
        const losses = document.getElementById('losses').value;
    
        //const winRate = (wins / (parseInt(wins) + parseInt(losses))) * 100;
        
        // Convert string to id
        const userDeckInt = parseInt(userDeck);
        const opponentDeckInt = parseInt(opponentDeck)
        
        // Convert wins and losses and int
        const winsInt = parseInt(wins);
        const lossesInt = parseInt(losses);

        fetch('http://localhost:8080/matches', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user_deck_id: userDeckInt,
                opponent_deck_id: opponentDeckInt,
                victories: winsInt,
                defeats: lossesInt
            }),
        })
        .then(response => {
            if(!response.ok) {
                return response.text().then(text => {
                    console.log('Error response:', text);
                    throw new Error(text);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log('Response data:', data);
            alert('WinRate updated sucessfully!');
            // Add a new row to the table with updated stats
            const table = document.getElementById('deckStatsTable').getElementsByTagName('tbody')[0];
            const newRow = table.insertRow();

            newRow.insertCell(0).textContent = deckNames[userDeckInt] || userDeckInt;
            newRow.insertCell(1).textContent = deckNames[opponentDeckInt] || opponentDeckInt;
            newRow.insertCell(2).textContent = wins;
            newRow.insertCell(3).textContent = losses;

            // Calculate and display the win rate
            const totalMatches = winsInt + lossesInt;
            const winRate = totalMatches ? (winsInt / totalMatches) * 100 : 0; // Avoid divison for zero
            newRow.insertCell(4).textContent = winRate.toFixed(2) + "%";
        
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
});