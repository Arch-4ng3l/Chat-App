let conns;
let acc

function login() {
    const password = document.getElementsByName("password")[0].value;
    acc = {
        "password": password,
    };
    fetchConns();
}

function shutdown() {
    console.log("sh");
    fetch("/api/admin/shutdown", {
        method: "POST",
    })
}

function shutdownReq() {
    document.getElementById("shutdown").style.display = "block";
}

function shutdownCancel() {
    document.getElementById("shutdown").style.display = "none";
}

function fetchConns() {
    const accJson = JSON.stringify(acc);
    fetch("/api/admin/login", {
        method: "POST",
        body: accJson,
    })
    .then(response => {
        if(response.ok) {
            return response.json();
        } else {
            response.json()
            .then(error => {
                console.log(error["error"]);
            });
        }
    })
    .then(data => {
        if(data !== undefined) {
            document.getElementById("login").style.display = "none";
            document.getElementById("admin").style.display = "block";
            conns = data["conns"];
            const doc = document.getElementById("conns")
            doc.innerHTML = "";
            doc.style.display = "block";
            conns.forEach(conn => {
                let div = document.createElement("div");
                div.appendChild(document.createTextNode("Username: " + conn));
                div.classList.add("connection");
                doc.appendChild(div);
            });
        }
    })
}
