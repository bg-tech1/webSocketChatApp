// WebSocket接続を作成
const socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = () => {
    console.log("WebSocket connection established");
};

socket.onerror = (error) => {
    console.error("WebSocket error:", error);
};

// サーバーからのメッセージを受信
socket.onmessage = function (event) {
    // メッセージデータ
    const data = JSON.parse(event.data);
    //上位存在
    const messageContainer = document.getElementById("chatMessagesContainer");
    // メッセージを入れるコンテナ
    const messageDiv = document.createElement("div");
    messageDiv.className = "d-flex align-items-center mb-3";
    messageContainer.appendChild(messageDiv);
    // iconとusername
    const iconDiv = document.createElement("div");
    const icon = document.createElement("i");
    icon.className = "bi bi-chat-left-fill fs-4 m-0";
    iconDiv.appendChild(icon);
    const usernameDiv = document.createElement("div");
    usernameDiv.className = "text-primary m-0";
    usernameDiv.textContent = data.username;
    iconDiv.appendChild(usernameDiv);
    messageDiv.appendChild(iconDiv);
    //メッセージ
    const message = document.createElement("div");
    message.className = "ms-3";
    message.textContent = data.message;
    messageDiv.appendChild(message);
};

// メッセージを送信
document.getElementById("sendButton").addEventListener("click", () => {
    const message = document.getElementById("messageInput");
    const username = new URLSearchParams(window.location.search).get("username");
    const messageBody = {
        message: message.value,
        username: username
    }
    socket.send(JSON.stringify(messageBody));
    message.value = "";
});