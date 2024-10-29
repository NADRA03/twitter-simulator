export function LoginTemplate(error) {
    let errorElement = "";
    if (error != "") {
        errorElement = `<p id="Error" style="color: red; font-size: small">${error}</p>`;
    }

    const template = `
    <section id="login-form">
     <div class="login">
        <form class="form" id="loginform" action="/login" method="post">
            <input type="hidden" name="form" value="login">
            <h1 class="h1">Log In</h1>
            <label class="lable">Username</label>
            <input type="text" id="loginusername" name="loginusername" placeholder="Enter Your Username" required><br>
            <label class="lable">Password</label>
            <input type="password" id="loginpassword" name="loginpassword" placeholder="Enter Your Password" required><br>
            <br>
            <br>
            <button type="submit" id="loginButton" onclick="window.location.href='./home'">Log In</button><br>
            <label class="lable">Don't Have Account?</label>
            <button type="button" onclick="window.location.href='./signup'">Sign Up</button>
        </form>
    </div>
    </section>
    `
    return template
}


export function SigninTemplate(error) {
    let errorElement = "";
    if (error != "") {
      errorElement = `<p id="Error" style="color: red; font-size: small">${error}</p>`;
    }

    const template = `
    <section id="sigin-form">
     <div class="sigin">
        <form class="form" id="signupform" action="/signup" method="post">
            <input type="hidden" name="form" value="signup">
            <h1 class="h1">Sign Up</h1>
            <label class="lable">Username</label>
            <input type="text" id="username" name="username" placeholder="Enter Your Username" required><br>
            <label class="lable">Email</label>
            <input type="email" id="email"  name="email" placeholder="Enter Your Email" required><br>
            <label class="lable">Password</label>
            <input type="password" id="password" name="password" placeholder="Enter Your Password" required><br>
            <label class="lable">First Name</label>
            <input type="text" id="firstname" name="firstname" placeholder="First Name" required><br>
            <label class="lable">Last Name</label>
            <input type="text" id="lastname" name="lastname" placeholder="Last Name" required><br>
            <label class="lable">Age</label>
            <input type="number" id="age" name="age" placeholder="First Name" required><br>
            <label class="lable">Gender</label>
            <select id="gender" name="gender" required>
                <option value="male">Male</option>
                <option value="female">Female</option>
            </select><br>
            <button type="submit" id="signupButton" class="signupButton" >Sign Up</button><br>
        </form>
    </div>
    </section>
    `
    return template
}

export function FirstPage(error) {
    let errorElement = "";
    if (error != "") {
      errorElement = `<p id="Error" style="color: red; font-size: small">${error}</p>`;
    }

    const template = `
       <div class="background">
        <div class="title">
            <h1 class="t1">Welcome</h1>
            <h1 class="t1">To</h1>
            <h1 class="t1">Real Time Forum </h1>
        </div>
        <div id="buttons" class="buttons">
            <button onclick="window.location.href='/signup'" class="button1" id="btn1">Sign-In</button>
            <button onclick="window.location.href='/login'" class="button1" id="btn2">Log-In</button>
        </div>
    </div>
    `
    return template
}

export function homepage(error) {
    let errorElement = "";
    if (error != "") {
      errorElement = `<p id="Error" style="color: red; font-size: small">${error}</p>`;
    }

    const template = `
      <header>
        <h1 class="site-title">Home Page</h1>
        <div class="links-section">
            <button onclick="window.location.href='home'">Main Page</button>
            <button onclick="window.location.href='createPost'">Create a Post</button>
            <button onclick="window.location.href='/logout'">Logout</button>
        </div>
    </header>

    <main class="main-content">
        <div class="post-container">
            {{range .Posts}}
            <div class="post-card">
                <div class="post-meta">
                    <div class="username"><i class="fas fa-user-circle"></i> {{.Username}}</div>
                    <div class="post-date">Posted on: {{ .CreatedAt.Format "02/01/2006 03:04 PM" }}</div>
                </div>
                <h4 class="post-title">{{.Title}}</h4>
                <p class="post-text">{{.Content}}</p>
                <div class="post-divider"></div>
                <div class="interactions">
                    <div class="category-dropdown">
                        <select id="category-dropdown" name="category">
                            <option value="">View Category</option>
                            {{range  $catName := .CategoriesNames}}
                            <option value="{{$catName}}">{{$catName}}</option>
                            {{end}}
                        </select>
                    </div>
                    <div class="comment-link">
                        <h4><a href="/post/{{.ID}}"><i class="fas fa-comment"></i> Comments</a></h4>
                    </div>
                </div>
            </div>
            {{else}}
            <div class="no-posts-message">
                <p>No posts available.</p>
            </div>
            {{end}}
        </div>
    </main>
    `
    return template
}