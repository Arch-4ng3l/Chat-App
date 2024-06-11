var socket = null;
var username = "";
var contactList = [];
var messageList = [];
var contacts = {};
var requestList = [];
var req = false;
var set = false;
var selected = "";
var friendWindow;
var settings;
const sound = new Audio("sound.mp3")
const numStars = 200;



class Message {
    constructor(name, content, side) {
        this.name = name;
        this.content = content;
        this.side = side;
    }
    createHTML() {
        const msg = document.createElement("div");
        const msgNode = document.createTextNode(this.content);
        msg.appendChild(msgNode);
        msg.classList.add("chat-msg");
        msg.classList.add(this.side);
        return msg;
    }
}

class Contact {
    constructor(name) {
        this.name = name;
    }

    createHTML() {
        const contact = document.createElement("div"); 
        const nameNode = document.createTextNode(this.name);
        contact.appendChild(nameNode);
        contact.id = this.name;
        contact.classList.add("contact");
        contact.onclick = selectContact(this.name);
        return contact;
    }
}


class FriendRequest {
    constructor(name) {
        this.name = name;
    }

    createHTML() {
        const contact = document.createElement("div"); 
        const nameNode = document.createTextNode(this.name);
        const button = document.createElement("button");
        button.appendChild(document.createTextNode("accept"));
        button.onclick = acceptFriend(this.name, contact);

        contact.id = this.name;
        contact.appendChild(nameNode);
        contact.appendChild(button);
        contact.classList.add("friend-request");
        return contact;
    }
}

function selectContact(name) {
    return fn = () => {
        if(selected === "") {
            sound.play();
            document.getElementById("chat-container").classList.add("show");
            document.getElementById("contacts").classList.add("show");
            document.getElementById("chat-container").style.display = "block";
        } 
        document.getElementById("msg-name").textContent = name;
        selected = name;
        let msgs = contacts[name];
        if(msgs === undefined) {
            contacts[name] = [];
            msgs = [];
        }
        messageList.innerText = "";
        msgs.forEach(element => {
            messageList.appendChild(element.createHTML());
        });
        messageList.scrollTop = messageList.scrollHeight;
    }
}



function acceptFriend(name, obj) {
    return fn = () => {
        fetch("/api/fr/accept", {
            method: "POST",
            body: JSON.stringify({
                "user1": name,
                "user2": getCookieValue("name"),
            }),
        });
        contactList.appen
        const contact = new Contact(name);
        const div = contact.createHTML();
        contactList.appendChild(div);
        obj.style.display = "none";
    }
}

function login() {
    const username = document.getElementsByName("username")[0].value;
    const password = document.getElementsByName("password")[0].value;

    const acc = {
        "name": username,
        "password": password,
    };

    const accJson = JSON.stringify(acc);
    fetch("/api/login", {
        method: "POST",
        body: accJson,
    })
    .then(response => {
        return response.json();

    })
    .then(data => {
        if(data !== undefined) {
            handleToken(data, username);
        }
    })
}

function signup() {
    const username = document.getElementsByName("username")[0].value;
    const email= document.getElementsByName("email")[0].value;
    const password = document.getElementsByName("password")[0].value;

    let acc = {
        "name": username,
        "email": email,
        "password": password,
    };

    const accJson = JSON.stringify(acc);

    fetch("/api/signup", {
        method: "POST",
        body: accJson,
    })
    .then(response => {
        return response.json();
    })
    .then(data => {
        handleToken(data, username);
    });
}

function createStar() {
    const star = document.createElement("div");
    star.className = "star";
    star.style.left = `${Math.random() * 2500}px`;
    star.style.top = `${Math.random() * 2500}px`;
    document.getElementById("stars").appendChild(star);
}

function createStars() {
    for (let i = 0; i < numStars; i++) {
        createStar();
    }
}

function connectToServerIndex() {
    createStars();
    revTransition();
    let token = getCookieValue("token")
    if(token !== "" && token !== undefined) {
        window.location = "/chat.html";
    }
}

function messageHandle(event) {
    let jsonObj = JSON.parse(event.data);

    if(jsonObj["sender"] !== undefined) {
        const content = jsonObj["msg"];
        let [chat, site] = jsonObj["sender"] === username ? 
            [jsonObj["receiver"], "right"] 
            :[jsonObj["sender"], "left"];

        if(contacts[chat] === undefined) {
            contacts[chat] = [];
        }
        const msg = new Message(chat, content, site);
        contacts[chat].push(msg);
        if(chat === selected && selected !== "") {
            addMessage(msg);
        }
        return;
    }
    const contactName = jsonObj["user1"] !== username ? jsonObj["user1"] : jsonObj["user2"];
    let contact;
    if(jsonObj["accepted"] === true) {
        contact = new Contact(contactName);
        const div = contact.createHTML();
        contactList.appendChild(div);
    } else {
        contact = new FriendRequest(contactName);
        const div = contact.createHTML();
        requestList.appendChild(div);
    }
    return
}

function addMessage(msg) {
    const div = msg.createHTML();
    messageList.appendChild(div);
    messageList.scrollTop = messageList.scrollHeight;
}

function openRequestWindow() {
    document.getElementById("friend-text").innerHTML = ""
    requestList.classList.toggle("show");
    req = !req;
    if(req && set) {
        set = false;
        settings.classList.toggle("show");
    }
}

function openSettings() {
    settings.classList.toggle("show");
    set = !set;
    if(req && set) {
        req = false;
        requestList.classList.toggle("show");
    }
}

function logout(){
    setCookie("token", "");
    setCookie("name", "");
    username = "";
    if(socket !== null) {
        socket.close();
        socket = null;
    }
    if(window.location == "chat.html") {
        document.getElementById("chat-container").classList.remove("show");
    }
    transition("/");
}


function connectToServer() {
    createStars();
    revTransition();
    contactList = document.getElementById("contacts");
    messageList= document.getElementById("messages");
    requestList= document.getElementById("requests");
    settings = document.getElementById("settings");
    const enter = (event) => {
        if (event.key === 'Enter') {
            sendMessage();
            document.getElementById("content").blur();
        }
    }

    document.addEventListener('keyup', enter);
    document.getElementById("content").addEventListener('keyup', enter);
    if(getCookieValue("token") != "") {
        const str = "ws://" + window.location.hostname +":3000/api/ws";
        username = getCookieValue("name");
        try {
            socket = new WebSocket(str);
            socket.onmessage = messageHandle;
        } 
        catch (e) {
            logout()
        }
        document.getElementById("username").appendChild(document.createTextNode("Welcome " + username));

    } else {
        logout()
    }
}



function addFriend() {
    const input = document.getElementById("friendname");
    const username = getCookieValue("name")
    if(input.value == username) {
        input.value = "";
        return
    }
    const friend = JSON.stringify({
        "user1": username,
        "user2": input.value,
    });
    fetch("/api/fr", {
        method: "POST",
        body: friend,
    })
    .then(response => {
        input.value = "";
        if(response.ok) {
            document.getElementById("friend-text").innerHTML = "Anfrage geschickt"
        } else {
            document.getElementById("friend-text").innerHTML = "Anfrage konnte nicht gesendet werden"
        }
    });
}

function sendMessage() {
    if(selected === "") {
        return;
    }
    const content = document.getElementById("content").value;
    if (content.replace(/\s+/g, "") === "") {
        return
    };
    socket.send(JSON.stringify({
        "receiver": selected,
        "msg": content,
    }));
    if(contacts[selected] === undefined) {
        contacts[selected] = [];
    }
    const msg = new Message(selected, content, "right");
    contacts[selected].push(msg);
    if(selected !== "") {
         addMessage(msg)
    }
    document.getElementById("content").value = "";
}


function transition(location) {
    document.getElementById("loading").classList.add("transition");
    document.getElementById("loading").style.display = "block";
    setTimeout(function () {
        window.location = location;
    }, 1900);
}

function revTransition() {
    document.getElementById("loading-rev").classList.add("rev-trans");
    document.getElementById("loading-rev").style.display = "block";
    setTimeout(function () {
        document.getElementById("loading-rev").style.display = "none";
        document.getElementById("loading-rev").style.zIndex= "-100";

    }, 1900);
}

function handleToken(data, username) {
    if(data["error"] !== undefined) {
        let err = data["error"]
        document.getElementById("errors").innerHTML = err
        return
    }
    document.getElementById("errors").innerHTML = ""
    const token = data.token;
    setCookie("token", token);
    setCookie("name", username);
    transition("/chat.html");
}


function setCookie(key, value) {
    document.cookie =  key + "=" + value + "; SameSite=None; Secure";
}

function getCookieValue(name) {
    let token = document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)')?.pop() || ''
    return (token === undefined ? "" : token)
}

