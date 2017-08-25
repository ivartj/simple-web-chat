function SimpleWebChat(elementId, websocketAddress) {

	var swc = this;

	swc.element = document.getElementById(elementId);
	swc.messagesElement = swc.element.getElementsByClassName("swc-messages")[0];
	swc.messageFormElement = swc.element.getElementsByClassName("swc-message-form")[0];
	swc.messageInputElement = swc.element.getElementsByClassName("swc-message-input")[0];

	swc.messageFormElement.addEventListener("submit", function(ev) { swc.messageSubmitHandler(ev); }, { capture: true });
	swc.socket = new WebSocket(websocketAddress);
	swc.socket.onmessage = swc.socket.onclose = function(ev) { swc.socketHandle(ev); };
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
		messageElement.innerText = "Error: Unexpected 'type' field value '" + msg.type + "' in JSON message.";
		break;
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
		var msg = JSON.parse(ev.data);
		this.addMessage(msg);
		break;
	case "close":
		this.addMessage({
			type: "client-system",
			text: "Connection closed: " + ev.reason
		});
		break;
	default:
		this.addMessage({
			type: "client-system",
			text: "Error: Unexpected 'type' field value '" + ev.type + "' in JSON message received over WebSocket."
		})
	}
};

// main
var chat;
window.onload = function() {
	var protocol = "ws:";
	if(window.location.protocol == "https:")
		protocol = "wss:";
	chat = new SimpleWebChat("simple-web-chat", protocol + "//" + window.location.host + "/websocket");
};

