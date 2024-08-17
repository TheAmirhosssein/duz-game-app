const socket = new WebSocket('ws://127.0.0.1:8080/ws/');
const urlParams = new URLSearchParams(window.location.search);
const userId = urlParams.get('user_id');
const gameId = urlParams.get('game_id');
const username = urlParams.get('username');
let isUserTurn

socket.addEventListener('open', function (event) {
    console.log('WebSocket conection opened');
}); n

socket.addEventListener('message', function (event) {
    message = JSON.parse(event.data)
    if (message.game_id == gameId) {
        console.log(message)
        if (message.type == "join_game") {
            changeTurn(message.message.user_sign)
        }
    }
});

socket.addEventListener('close', function (event) {
    console.log('WebSocket connection closed');
});

socket.addEventListener('error', function (event) {
    console.error('WebSocket error:', event);
});

function sendMessage(message) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    } else {
        console.log('WebSocket is not open. Ready state is:', socket.readyState);
    }
}


function startGame() {
    if (userId !== null && gameId !== null) {
        var startGameData = {
            type: "join_match",
            game_id: gameId,
            user_id: userId,
        }
        sendMessage(startGameData)
    }
    let hiddenItems = document.getElementsByClassName("hide");
    [...hiddenItems].forEach(element => {
        element.classList.remove("hide");
    });
    document.getElementById("start").classList.add("hide")
    document.getElementById("playerName").innerHTML = username
}

function closeConnection() {
    socket.close();
}

function changeTurn(turn) {
    let turnText
    if (turn == "O") {
        turnText = "نوبت حریف"
    } else {
        turnText = "نوبت شما"
        isUserTurn = true
    }
    document.getElementById("turn").innerHTML = turnText
}
