import * as temp from "./template.js";

export function renderLoginTemplate(error = "") {
    const templete = temp.LoginTemplate(error);
    const style = `<link href="/static/login.css" rel="stylesheet">`;
    document.getElementById("content").innerHTML = style + templete;
    renderLoginEvent();
}

export function renderSignInTemplate(error = "") {
    const templete = temp.SigninTemplate(error);
    const style = `<link href="/static/signin.css" rel="stylesheet">`;
    document.getElementById("content").innerHTML = style + templete;
    renderSignInEvent();
    
}

export function renderFirstPage(error = "") {
    const templete = temp.FirstPage(error);
    const style = `<link href="/static/home.css" rel="stylesheet">`;
    document.getElementById("content").innerHTML = style + templete;
    renderFirstPageEvent();
}

export function renderHomePage(error = "") {
    const template = temp.homepage(error);
    const style = `<link href="/static/homepage.css" rel="stylesheet">`;
    document.getElementById("content").innerHTML = style + template;
    renderHomePageEvent();
}

function renderHomePageEvent() {
    // Add any event listeners here if needed
    // For example, if you want to handle category selection or any other interactions
}

function renderLoginEvent() {
    const loginform = document.getElementById("loginform");
    const loginsection = document.getElementById("login-form");
    const loginbutton = document.getElementById("loginButton");
    loginbutton.addEventListener('click', () => {
        loginform.addEventListener('submit', async (event) => {
            event.preventDefault();

            const formData = new FormData(loginform);
            const data = Object.fromEntries(formData);

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams(data),
                });

                if (response.ok) {
                    alert('logIn Successful!');
                    window.location.href = '/home';
                } else {
                    alert(`logIn failed: username invalid`);
                }
            } catch (error) {
                alert('An error occurred: ');
            }
        });
    });
}


function renderSignInEvent() {
    const signupform = document.getElementById("signupform");

    signupform.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new FormData(signupform);
        const data = Object.fromEntries(formData);

        try {
            const response = await fetch('/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: new URLSearchParams(data),
            });

            if (response.ok) {
                alert('Registration successful!');
                window.location.href = '/home';
            } else {
                const errorText = await response.text();
                alert(`Registration failed: ${errorText}`);
            }
        } catch (error) {
            alert('An error occurred: ' + error.message);
        }
    });
}

function renderFirstPageEvent() {
    const loginbutton = document.getElementById("btn2");
    const signinbutton = document.getElementById("btn1");
    const buttons = document.getElementById("buttons");

    loginbutton.addEventListener('click', () => {
        buttons.addEventListener('submit', async (event) => {
            event.preventDefault();

            const formData = new FormData(buttons);
            const data = Object.fromEntries(formData);

            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams(data),
                });

                if (response.ok) {
                    alert('logIn Successful!');
                    window.location.href = '/home';
                } else {
                    alert(`logIn failed: username invalid`);
                }
            } catch (error) {
                alert('An error occurred: ');
            }
            
        });
    });


    signinbutton.addEventListener('click', () => {
        buttons.addEventListener('submit', async (event) => {
            event.preventDefault();

            try {
                const response = await fetch('/signup', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams(data),

                });

                if (response.ok) {
                    alert('Registration successful!');
                    window.location.href = '/home';
                } else {
                    alert(`Registration failed:username invalid `);
                }
            } catch (error) {
                alert('An error occurred: ');
            }
            
        });
    });

}