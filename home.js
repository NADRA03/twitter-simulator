let hasRendered = false;
const defaultImageUrl = '/assets/user2.png';

export async function render(posts) {
    if (hasRendered) return; 
    hasRendered = true; 

    const postsContainer = document.getElementById("app");
    postsContainer.innerHTML = ''; 
    const homeElement = document.createElement("div");
    homeElement.classList.add("home");

    const user = await getUserDetails();
    const userID = user?.id;
    const userImageUrl = user?.image_url || defaultImageUrl;
    const inputContainer = document.createElement("div");
    inputContainer.classList.add("input-container");
    const userImage = document.createElement("img");
    userImage.src = userImageUrl;
    userImage.classList.add("user-image");
    userImage.onerror = () => {
        userImage.src = defaultImageUrl; 
    };

    const inputField = document.createElement("input");
    inputField.type = "text";
    inputField.placeholder = "What is happening?!";

    const postButton = document.createElement("button");
    postButton.textContent = "Post";
    postButton.classList.add("post-button");
    postButton.disabled = true; 

    inputField.addEventListener("input", () => {
        if (inputField.value.trim() !== "") {
            postButton.disabled = false; 
        } else {
            postButton.disabled = true; 
        }
    });

    postButton.addEventListener("click", () => {
        if (userID) {
            const content = inputField.value.trim();
            const title = "New Post"; 
            
            createPost(userID, title, content);
            
            inputField.value = '';
            postButton.disabled = true; 
        } else {
            console.error("User is not authenticated.");
        }
    });

    const line = document.createElement("div");
    line.classList.add("line");

    const inputFieldContainer = document.createElement("div");
    inputFieldContainer.classList.add("input-field-container");
    inputFieldContainer.appendChild(userImage);
    inputFieldContainer.appendChild(inputField);
    inputContainer.appendChild(inputFieldContainer);
    inputContainer.appendChild(line); 
    inputContainer.appendChild(postButton);
    homeElement.appendChild(inputContainer);

    if (!Array.isArray(posts)) {
        console.error("Expected an array, but got:", posts);
        posts = []; 
    }
    if (!posts || posts.length === 0) {
        const noPostsMessage = document.createElement("p");
        noPostsMessage.textContent = "No posts available";
        homeElement.appendChild(noPostsMessage);
    } else {
        console.log(posts);
        posts.forEach(post => {
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            const formattedDate = formatCreatedAt(post["created-at"]);
            const imageURL = post.image_url || defaultImageUrl; 
            postElement.innerHTML = `
                <p><img src="${imageURL}" onError="this.onerror=null; this.src='${defaultImageUrl}'" class="post-image"/> By ${post.Username} on ${formattedDate}</p>
                <p class="content">${post.content}</p>
            `;
            postElement.addEventListener('click', () => {
                window.location.href = `/post/${post.id}`;
            });
            const postImage = postElement.querySelector(".post-image");
            postImage.addEventListener('click', (e) => {
                e.stopPropagation(); 
                window.location.href = `/a_profile/${post["user-id"]}`;
            });
            homeElement.appendChild(postElement);
        });
    }

    postsContainer.appendChild(homeElement);
}

async function getUserDetails() {
    try {
        const response = await fetch('/User');
        if (!response.ok) {
            throw new Error(`Failed to fetch user details: ${response.statusText}`);
        }

        return await response.json(); 
    } catch (error) {
        console.error("Error fetching user details:", error);
        return null; 
    }
}

async function createPost(userID, title, content) {
    try {
        const response = await fetch('/createPost', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ userID, title, content }),
        });

        if (!response.ok) {
            throw new Error(`Failed to create post: ${response.statusText}`);
        }

        await fetchPosts();
    } catch (error) {
        console.error("Error creating post:", error);
    }
}

function formatCreatedAt(createdAt) {
    const now = new Date();
    const postDate = new Date(createdAt);

    const isToday =
        now.getFullYear() === postDate.getFullYear() &&
        now.getMonth() === postDate.getMonth() &&
        now.getDate() === postDate.getDate();

    if (isToday) {
        const diffMs = now - postDate;
        const diffMinutes = Math.floor(diffMs / 60000); 

        if (diffMinutes < 60) {
            return `${diffMinutes} minute(s) ago`;
        }

        const diffHours = Math.floor(diffMinutes / 60);
        return `${diffHours} hour(s) ago`;
    } else {
        return postDate.toISOString().split('T')[0];
    }
}

async function fetchPosts() {
    try {
        const response = await fetch('/posts');

        if (!response.ok) {
            throw new Error(`Failed to fetch posts: ${response.statusText}`);
        }

        const posts = await response.json(); 
        await render(posts);
    } catch (error) {
        const posts = [];
        await render(posts);
        console.error("Error fetching posts:", error);
        const postsContainer = document.getElementById("postsContainer");
        postsContainer.innerHTML = '<p>Error loading posts.</p>';
    }
}

export async function initialize() {
    await fetchPosts();  
}

document.addEventListener("DOMContentLoaded", initialize);
