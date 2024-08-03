// 1. Create a WebSocket connection
const socket = new WebSocket('ws://127.0.0.1:8080/ws/');

socket.addEventListener('open', function (event) {
    console.log('WebSocket connection opened');
    socket.send('Hello Server!');
});

socket.addEventListener('message', function (event) {
    console.log('Message from server:', event.data);
});

socket.addEventListener('close', function (event) {
    console.log('WebSocket connection closed');
});

socket.addEventListener('error', function (event) {
    console.error('WebSocket error:', event);
});

// 3. Send messages
function sendMessage(message) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.log('WebSocket is not open. Ready state is:', socket.readyState);
    }
}

const urlParams = new URLSearchParams(window.location.search);
const userId = urlParams.get('user_id');
const gameId = urlParams.get('game_id');

function startGame() {
    if (userId !== null && gameId !== null) {
        var startGameData = {
            type: "join_match",
            game_id: gameId,
            user_id: userId,
        }
        sendMessage(startGameData)
    }
}

// 4. Close the WebSocket connection
function closeConnection() {
    socket.close();
}
