//ログインを管理
document.getElementById("login").addEventListener("submit", () => {
    event.preventDefault()
    const reqJson = {
        username: document.getElementById("username").value,
        pass: document.getElementById("pass").value
    }
    fetch("/auth/login/", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(reqJson)
    })
        .then(response => {
            return response.json();
        })
        .then((data) => {
            console.log(data);
            if (data.status_code == 200) {
                const queryParams = new URLSearchParams({ username: reqJson.username, access_token: data.access_token });
                window.location.href = "/chat/?" + queryParams.toString();
            } else {
                window.location.href = "/notfound/";
            }
        })
        .catch(error => {
            console.error("Error:", error);
        });
});