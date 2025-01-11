// Function to show error messages dynamically
function showErrorMessage(elementId, message) {
    const errorMessageElement = document.getElementById(elementId);
    console.log('Showing error message:', message); // Log the error message
    if (errorMessageElement) {
        errorMessageElement.textContent = message;
        errorMessageElement.style.display = "block";
    }
}

// Function to handle all event listeners related to login and signup forms
export function setupEventListeners() {
    const goToLoginButton = document.getElementById("goToLoginButton");
    const loginForm = document.getElementById("login-form");
    const closeLoginButton = document.getElementById("closeLoginForm");

    const goToSignupButton = document.getElementById("goToSignupButton");
    const signupForm = document.getElementById("signup-form");
    const closeSignupButton = document.getElementById("closeSignupForm");

    // Check if there's an error message in the URL query params and display it
    const urlParams = new URLSearchParams(window.location.search);
    const error = urlParams.get('error');
    const formType = urlParams.get('form'); // Identify which form caused the error
    console.log("Error from URL params:", error); // Check if error is captured

    if (error) {
        // Display the corresponding form based on the formType parameter
        if (formType === "login") {
            if (loginForm) {
                loginForm.style.display = "block";
                signupForm.style.display = "none"; // Hide signup form
                console.log("Displaying login error:", error); // Confirming if it's showing login error
                showErrorMessage('loginErrorMessage', error); // Show error for login
            }
        } else if (formType === "signup") {
            if (signupForm) {
                signupForm.style.display = "block";
                loginForm.style.display = "none"; // Hide login form
                console.log("Displaying signup error:", error); // Confirming if it's showing signup error
                showErrorMessage('signupErrorMessage', error); // Show error for signup
            }
        }
    }

    // Event listener to show login form and hide signup form
    if (goToLoginButton) {
        goToLoginButton.addEventListener("click", function() {
            console.log("Go to Login Button clicked!");
            loginForm.style.display = "block";
            signupForm.style.display = "none"; // Hide signup form
        });
    }

    // Event listener to close login form
    if (closeLoginButton) {
        closeLoginButton.addEventListener("click", function() {
            console.log("Close Login Button clicked!");
            loginForm.style.display = "none";
        });
    }

    // Event listener to show signup form and hide login form
    if (goToSignupButton) {
        goToSignupButton.addEventListener("click", function() {
            console.log("Go to Signup Button clicked!");
            signupForm.style.display = "block";
            loginForm.style.display = "none"; // Hide login form
        });
    }

    // Event listener to close signup form
    if (closeSignupButton) {
        closeSignupButton.addEventListener("click", function() {
            console.log("Close Signup Button clicked!");
            signupForm.style.display = "none";
        });
    }
}

// Render function to create the HTML structure
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

// Function to initialize the login module
export function initialize() {
    setupEventListeners(); // Set up event listeners after rendering the module
}

