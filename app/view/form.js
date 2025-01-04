document.getElementById("submitForm").addEventListener("submit", () => {
    event.preventDefault()
    const reqJson = {
        username: document.getElementById("usernameInput").value,
        pass: document.getElementById("passwordInput").value
    }
    fetch("/register/user/", {
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
            if (data.status_code == 201) {
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