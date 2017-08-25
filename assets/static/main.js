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
	this.socket.onmessage = this.socket.onclose = function(ev) { swc.socketHandle(ev); };
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
	case "client-system":
		messageElement.innerText = msg.text;
		break;
	default:
		// TODO
	}

	this.messagesElement.appendChild(messageElement);

};

SimpleWebChat.prototype.messageSubmitHandler = function(ev) {
	ev.preventDefault();
	this.messageInputElement.focus();
	if(this.messageInputElement.value.length == 0)
		return;
	this.socket.send(JSON.stringify({
		type: "message",
		text: this.messageInputElement.value
	}));
	this.messageInputElement.value = "";
};

SimpleWebChat.prototype.socketHandle = function(ev) {
	switch(ev.type) {
	case "message":
		// TODO handle parsing exceptions
		var msg = JSON.parse(ev.data);
		this.addMessage(msg);
		break;
	case "close":
		this.addMessage({
			type: "client-system",
			text: "Connection closed."
		});
		break;
	default:
		// TODO
	}
};

// main
var chat;
window.onload = function() {
	// TODO Non-hardcoded websocket address
	var protocol = "ws:";
	if(window.location.protocol == "https:")
		protocol = "wss:";
	chat = new SimpleWebChat("simple-web-chat", protocol + "//" + window.location.host + "/websocket");
};

