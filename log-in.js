
function showErrorMessage(elementId, message) {
    const errorMessageElement = document.getElementById(elementId);
    console.log('Showing error message:', message); 
    if (errorMessageElement) {
        errorMessageElement.textContent = message;
        errorMessageElement.style.display = "block";
    }
}

export function setupEventListeners() {
    const goToLoginButton = document.getElementById("goToLoginButton");
    const loginForm = document.getElementById("login-form");
    const closeLoginButton = document.getElementById("closeLoginForm");

    const goToSignupButton = document.getElementById("goToSignupButton");
    const signupForm = document.getElementById("signup-form");
    const closeSignupButton = document.getElementById("closeSignupForm");
    const urlParams = new URLSearchParams(window.location.search);
    const error = urlParams.get('error');
    const formType = urlParams.get('form');
    console.log("Error from URL params:", error); 

    if (error) {
        if (formType === "login") {
            if (loginForm) {
                loginForm.style.display = "block";
                signupForm.style.display = "none"; 
                console.log("Displaying login error:", error); 
                showErrorMessage('loginErrorMessage', error); 
            }
        } else if (formType === "signup") {
            if (signupForm) {
                signupForm.style.display = "block";
                loginForm.style.display = "none";
                console.log("Displaying signup error:", error); 
                showErrorMessage('signupErrorMessage', error); 
            }
        }
    }

    if (goToLoginButton) {
        goToLoginButton.addEventListener("click", function() {
            console.log("Go to Login Button clicked!");
            loginForm.style.display = "block";
            signupForm.style.display = "none"; 
        });
    }

    if (closeLoginButton) {
        closeLoginButton.addEventListener("click", function() {
            console.log("Close Login Button clicked!");
            loginForm.style.display = "none";
        });
    }

    if (goToSignupButton) {
        goToSignupButton.addEventListener("click", function() {
            console.log("Go to Signup Button clicked!");
            signupForm.style.display = "block";
            loginForm.style.display = "none"; 
        });
    }


    if (closeSignupButton) {
        closeSignupButton.addEventListener("click", function() {
            console.log("Close Signup Button clicked!");
            signupForm.style.display = "none";
        });
    }
}


export function render() {
    return `
    <div class="login-and-signup">
        <img src="/assets/logo.svg" class="logo" />
        <div class="nav">
            <h1 class="title">Happening now</h1>
            <h1 class="title2">Join today.</h1>
            <button id="goToSignupButton">Create account</button>
            <h1 class="text"><span class="line"></span> or <span class="line"></span></h1>
            <button id="goToLoginButton">Sign in</button>
        </div>
        
        <!-- Login Form -->
        <section id="login-form" style="display: none;">  
            <div class="login">
                <button id="closeLoginForm" style="position: absolute; top: 10px; left: 10px;">✖</button>
                <form class="form" id="loginForm" action="/log-in" method="post">
                    <input type="hidden" name="form" value="login">
                    <label id="username-label" class="label">Username</label>
                    <input type="text" id="loginusername" name="loginusername" placeholder="Enter Your Username" required><br>
                    <label class="label">Password</label>
                    <input type="password" id="loginpassword" name="loginpassword" placeholder="Enter Your Password" required><br>
                      <!-- Display error message for login -->
                      <div id="loginErrorMessage" style="color: red; display: none;"></div>
                    <button type="submit" id="loginButton">Log In</button><br>
                </form>
            </div>
        </section>

        <!-- Signup Form -->
        <section id="signup-form" style="display: none;">
            <div class="signup">
                <button id="closeSignupForm" style="position: absolute; top: 10px; left: 10px;">✖</button>
                <form action="/log-in/create-account" method="post">
                    <input type="hidden" name="form" value="signup">
                    <label id="username-label" class="label">Username</label>
                    <input type="text" id="username" name="username" placeholder="Enter Your Username" required><br>
                    <label class="label">Email</label>
                    <input type="email" id="email" name="email" placeholder="Enter Your Email" required><br>
                    <label class="label">Password</label>
                    <input type="password" id="password" name="password" placeholder="Enter Your Password" required><br>
                    <label class="label">First Name</label>
                    <input type="text" id="firstname" name="firstname" placeholder="First Name" required><br>
                    <label class="label">Last Name</label>
                    <input type="text" id="lastname" name="lastname" placeholder="Last Name" required><br>
                    <label class="label">Age</label>
                    <input type="number" id="age" name="age" placeholder="Age" required><br>
                    <label class="label">Gender</label>
                    <select id="gender" name="gender" required>
                        <option value="male">Male</option>
                        <option value="female">Female</option>
                    </select><br><br>
                    <!-- Display error message for signup -->
                    <div id="signupErrorMessage" style="color: red; display: none;"></div>
                    <button type="submit" id="signupButton" class="signupButton">Sign Up</button><br>
                </form>
            </div>
        </section>
    </div>
    `;
}


export function initialize() {
    setupEventListeners(); 
}

