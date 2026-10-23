const status_ = document.querySelector("#status");
const messages = document.querySelector("#messages");

const connectButton = document.querySelector("#connect");
const disconnectButton = document.querySelector("#disconnect");

const form = document.querySelector("#message-form");
const input = document.querySelector("#message");
const sendButton = document.querySelector("#send");

let socket;

function log(message) {
  const item = document.createElement("li");
  item.textContent = message;
  messages.prepend(item);
}

connectButton.addEventListener("click", () => {
  socket = new WebSocket("ws://localhost:8080/ws");

  log("Tentando conectar...");

  // Os listeners da conexão ficam aqui dentro!
  socket.addEventListener("open", () => {
    status_.textContent = "Conectado";

    connectButton.disabled = true;
    disconnectButton.disabled = false;
    sendButton.disabled = false;

    log("Conexão WebSocket estabelecida!");
  });

  socket.addEventListener("message", (event) => {
    log(`Recebido: ${event.data}`);
  });
  socket.addEventListener("error", () => {
    log("Ocorreu um erro na conexão.");
  });

  socket.addEventListener("close", (event) => {
    status_.textContent = "Desconectado";

    connectButton.disabled = false;
    disconnectButton.disabled = true;
    sendButton.disabled = true;

    log(`Conexão encerrada. Código: ${event.code}`);
  });
});

form.addEventListener("submit", (event) => {
  event.preventDefault();

  if (!socket || socket.readyState !== WebSocket.OPEN) {
    log("Não há conexão aberta.");
    return;
  }

  socket.send(input.value);

  log(`Enviado: ${input.value}`);
  input.value = "";
});

disconnectButton.addEventListener("click", () => {
  if (socket) {
    socket.close(1000, "Desconexão solicitada");
  }
});
