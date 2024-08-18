const socket = new WebSocket('ws://127.0.0.1:8080/ws/');
const urlParams = new URLSearchParams(window.location.search);
const userId = urlParams.get('user_id');
const gameId = urlParams.get('game_id');
const username = urlParams.get('username');
let isUserTurn = false;
let userSign = "";
let counter = 0;
let multiplication = 0;

socket.addEventListener('open', function (event) {
    console.log('WebSocket conection opened');
});

socket.addEventListener('message', function (event) {
    message = JSON.parse(event.data)
    if (message.game_id === gameId) {
        console.log(message)
        if (message.type === "join_game") {
            userSign = message.message.user_sign
            if (userSign === "X") {
                isUserTurn = true
            }
            changeTurn()
        }
        if (message.type === "move") {
            document.getElementById(message.message.square).classList.add(`${message.message.sign}-sign`)
            if (message.user_id !== userId) {
                isUserTurn = true
                changeTurn()
            }
        }
        if (message.type === "error" && message.user_id === userId) {
            alert(message.message.error)
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

function changeTurn() {
    let turnText
    if (!isUserTurn) {
        turnText = "نوبت حریف"
    } else {
        turnText = "نوبت شما"
    }
    document.getElementById("turn").innerHTML = turnText
}


function move(event) {
    if (!isUserTurn) {
        alert("نوبت شما نیست")
        return null
    }

    if (counter === 3) {
        alert("شما نمی توانید حرکتی انجام بدهید")
        return null
    }
    console.log(counter);

    const classList = event.target.classList;
    const emptySquare = "empty-square"
    const signClass = `${userSign}-sign`

    if (classList.contains(emptySquare)) {
        isUserTurn = false
        var snd = new Audio("/static/Voice/add.mp3");
        snd.play();
        event.target.classList.add(signClass);
        event.target.classList.remove(emptySquare);
        counter++;
        var moveData = {
            type: "move",
            game_id: gameId,
            user_id: userId,
            square: event.target.id,
        }
        sendMessage(moveData)
        changeTurn()
    }

}

function handleDoubleClick(event) {
    var snd = new Audio("/static/Voice/delete.mp3");
    snd.play();

    if (event.target.getAttribute('src') === "/static/img/2.png") {
        multiplication--;
        event.target.setAttribute('src', "/static/img/3.png");
    }
}

const squares = document.querySelectorAll('.square');

squares.forEach(function (square) {
    square.addEventListener('click', move);
});