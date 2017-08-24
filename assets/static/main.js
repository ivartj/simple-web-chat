function SimpleWebChat(elementId, websocketAddress) {
	// TODO: check for exception
	this.element = document.getElementById(elementId);
	this.messagesElement = this.element.getElementsByClassName("swc-messages")[0];
	this.messageFormElement = this.element.getElementsByClassName("swc-message-form")[0];
	this.messageInputElement = this.element.getElementsByClassName("swc-message-input")[0];

	// TODO clean up this
	var swc = this;
	this.messageFormElement.addEventListener("submit", function(ev) { swc.messageSubmitHandler(ev); }, { capture: true });

	this.socket = new WebSocket(websocketAddress);
	this.socket.onmessage = function(ev) { swc.messageReceiveHandler(ev); };
}

SimpleWebChat.prototype.addMessage = function(msg) {
	var messageElement = document.createElement("div");
	messageElement.classList.add("swc-message");
	if(msg.color)
		messageElement.style.color = msg.color;

	switch(msg.type) {
	case "message":
		messageElement.innerText = msg.user + ": " + msg.text;
		break;
	case "join":
		messageElement.innerText = "* " + msg.user + " joined";
		break;
	case "leave":
		messageElement.innerText = "* " + msg.user + " left";
		break;
	default:
		// TODO
	}

	this.messagesElement.appendChild(messageElement);

};

SimpleWebChat.prototype.messageSubmitHandler = function(ev) {
	ev.preventDefault();
	if(this.messageInputElement.value.length == 0)
		return;
	this.socket.send(JSON.stringify({
		type: "message",
		text: this.messageInputElement.value
	}));
	this.messageInputElement.value = "";
};

SimpleWebChat.prototype.messageReceiveHandler = function(ev) {
	// TODO handle exceptions
	var msg = JSON.parse(ev.data);
	this.addMessage(msg);
};

// main
var chat;
window.onload = function() {
	// TODO Non-hardcoded websocket address
	chat = new SimpleWebChat("simple-web-chat", "ws://" + window.location.host + "/websocket");
};

