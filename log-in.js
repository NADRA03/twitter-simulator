export function setupEventListeners() {
    const goToLoginButton = document.getElementById("goToLoginButton");
    const loginForm = document.getElementById("login-form");
    const closeLoginButton = document.getElementById("closeLoginForm");

    const goToSignupButton = document.getElementById("goToSignupButton");
    const signupForm = document.getElementById("signup-form");
    const closeSignupButton = document.getElementById("closeSignupForm");

    if (goToLoginButton) {
        goToLoginButton.addEventListener("click", function() {
            console.log("Go to Login Button clicked!");
            loginForm.style.display = "block";
            signupForm.style.display = "none"; // Hide signup form
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
            loginForm.style.display = "none"; // Hide login form
        });
    }

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
        <svg class="logo" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Layer_1" width="24px" height="24px" viewBox="0 0 24 24" style="enable-background:new 0 0 24 24;" xml:space="preserve"><path d="M14.095479,10.316482L22.286354,1h-1.940718l-7.115352,8.087682L7.551414,1H1l8.589488,12.231093L1,23h1.940717  l7.509372-8.542861L16.448587,23H23L14.095479,10.316482z M11.436522,13.338465l-0.871624-1.218704l-6.924311-9.68815h2.981339  l5.58978,7.82155l0.867949,1.218704l7.26506,10.166271h-2.981339L11.436522,13.338465z"/></svg>
        <div class="nav">
            <button id="goToLoginButton">Go to Login</button>
            </br>
            </br>
            <button id="goToSignupButton">Go to Signup</button>
        </div>
        
        <!-- Login Form -->
        <section id="login-form" style="display: none;">  
            <div class="login">
                <button id="closeLoginForm" style="position: absolute; top: 10px; right: 10px;">✖</button>
                <form class="form" id="loginForm" action="/log-in" method="post">
                    <input type="hidden" name="form" value="login">
                    <h1 class="h1">Log In</h1>
                    <label class="label">Username</label>
                    <input type="text" id="loginusername" name="loginusername" placeholder="Enter Your Username" required><br>
                    <label class="label">Password</label>
                    <input type="password" id="loginpassword" name="loginpassword" placeholder="Enter Your Password" required><br>
                    <button type="submit" id="loginButton">Log In</button><br>
                    <label class="label">Don't Have an Account?</label>
                    <button type="button" onclick="window.location.href='./sign-up'">Sign Up</button>
                </form>
            </div>
        </section>

        <!-- Signup Form -->
        <section id="signup-form" style="display: none;">
            <div class="signup">
                <button id="closeSignupForm" style="position: absolute; top: 10px; right: 10px;">✖</button>
                <form action="/log-in/create-account" method="post">
                    <input type="hidden" name="form" value="signup">
                    <h1 class="h1">Sign Up</h1>
                    <label class="label">Username</label>
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
                    </select><br>
                    <button type="submit" id="signupButton" class="signupButton">Sign Up</button><br>
                </form>
            </div>
        </section>
    `;
}

// Function to initialize the login module
export function initialize() {
    setupEventListeners(); // Set up event listeners after rendering the module
}