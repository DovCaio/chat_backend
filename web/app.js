const userSelect = document.querySelector("#user");
const recipientSelect = document.querySelector("#recipient");
const connectButton = document.querySelector("#connect");
const disconnectButton = document.querySelector("#disconnect");
const status = document.querySelector("#status");
const form = document.querySelector("#message-form");
const input = document.querySelector("#message");
const sendButton = document.querySelector("#send");
const messages = document.querySelector("#messages");

let socket = null;

function logMessage(author, content, ownMessage = false) {
  const item = document.createElement("li");
  item.className = ownMessage ? "own-message" : "other-message";

  const authorElement = document.createElement("strong");
  authorElement.textContent = author;

  const contentElement = document.createElement("p");
  contentElement.textContent = content;

  item.append(authorElement, contentElement);
  messages.appendChild(item);
  messages.scrollTop = messages.scrollHeight;
}

function updateControls(connected) {
  connectButton.disabled = connected;
  disconnectButton.disabled = !connected;
  sendButton.disabled = !connected;

  userSelect.disabled = connected;
}

connectButton.addEventListener("click", () => {
  const user = userSelect.value;

  socket = new WebSocket(
    `ws://${location.host}/ws?user=${encodeURIComponent(user)}`,
  );

  status.textContent = "Conectando...";
  connectButton.disabled = true;

  socket.addEventListener("open", () => {
    status.textContent = `Conectado como ${user}`;
    updateControls(true);
    logMessage("Sistema", "Conexão estabelecida.");
  });

  socket.addEventListener("message", (event) => {
    const message = JSON.parse(event.data);

    logMessage(message.from, message.content);
  });

  socket.addEventListener("error", () => {
    logMessage("Sistema", "Ocorreu um erro na conexão.");
  });

  socket.addEventListener("close", () => {
    status.textContent = "Desconectado";
    updateControls(false);
    logMessage("Sistema", "Conexão encerrada.");
    socket = null;
  });
});

disconnectButton.addEventListener("click", () => {
  socket?.close(1000, "Desconexão solicitada");
});

userSelect.addEventListener("change", () => {
  recipientSelect.value = userSelect.value === "caio" ? "ana" : "caio";
});

recipientSelect.addEventListener("change", () => {
  if (recipientSelect.value === userSelect.value) {
    recipientSelect.value = userSelect.value === "caio" ? "ana" : "caio";
  }
});

form.addEventListener("submit", (event) => {
  event.preventDefault();

  if (!socket || socket.readyState !== WebSocket.OPEN) {
    return;
  }

  const content = input.value.trim();

  if (!content) {
    return;
  }

  const recipient = recipientSelect.value;

  socket.send(
    JSON.stringify({
      to: recipient,
      content,
    }),
  );

  logMessage(`Você → ${recipient}`, content, true);
  input.value = "";
  input.focus();
});
