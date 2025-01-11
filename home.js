let hasRendered = false;
const defaultImageUrl = '/assets/user2.png';

export async function render(posts) {
    if (hasRendered) return; // Exit if already rendered
    hasRendered = true; // Set the flag to true to prevent future calls

    const postsContainer = document.getElementById("app");
    postsContainer.innerHTML = ''; // Clear the posts container

    // Create a new element called "home" to wrap everything inside the "app"
    const homeElement = document.createElement("div");
    homeElement.classList.add("home");

    const user = await getUserDetails();
    const userID = user?.id;
    const userImageUrl = user?.image_url || defaultImageUrl;

    // Add "What is happening?" input field, image, and "Post" button
    const inputContainer = document.createElement("div");
    inputContainer.classList.add("input-container");

    // Add user image
    const userImage = document.createElement("img");
    userImage.src = userImageUrl;
    // userImage.alt = "User Image";
    userImage.classList.add("user-image");
    userImage.onerror = () => {
        userImage.src = defaultImageUrl; // Path to the default image
    };

    const inputField = document.createElement("input");
    inputField.type = "text";
    inputField.placeholder = "What is happening?!";

    const postButton = document.createElement("button");
    postButton.textContent = "Post";
    postButton.classList.add("post-button");
    postButton.disabled = true; // Initially disabled

    // Enable button when there is input text
    inputField.addEventListener("input", () => {
        if (inputField.value.trim() !== "") {
            postButton.disabled = false; // Enable button
        } else {
            postButton.disabled = true; // Disable button if input is empty
        }
    });

    // Post button click event
    postButton.addEventListener("click", () => {
        if (userID) {
            const content = inputField.value.trim();
            const title = "New Post"; // Or you can allow the user to enter a title too
            
            // Call CreatePost function here
            createPost(userID, title, content);
            
            // Reset input after post
            inputField.value = '';
            postButton.disabled = true; // Disable button again
        } else {
            console.error("User is not authenticated.");
        }
    });

    // Line element to separate input and button
    const line = document.createElement("div");
    line.classList.add("line");

    // Append elements to the input container
    const inputFieldContainer = document.createElement("div");
    inputFieldContainer.classList.add("input-field-container");
    inputFieldContainer.appendChild(userImage);
    inputFieldContainer.appendChild(inputField);

    inputContainer.appendChild(inputFieldContainer);
    inputContainer.appendChild(line); // Add line here
    inputContainer.appendChild(postButton);

    homeElement.appendChild(inputContainer);

    // Check if posts is an array
    if (!Array.isArray(posts)) {
        console.error("Expected an array, but got:", posts);
        posts = []; // Assign an empty array if it's not an array
    }

    // Check if the posts array is empty or invalid
    if (!posts || posts.length === 0) {
        const noPostsMessage = document.createElement("p");
        noPostsMessage.textContent = "No posts available";
        homeElement.appendChild(noPostsMessage);
    } else {
        console.log(posts);
        // Iterate over the posts and create elements for each
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
                // Redirect to the post's detailed page
                window.location.href = `/post/${post.id}`;
            });
            homeElement.appendChild(postElement);
        });
    }

    // Append the "home" element to the "app" container
    postsContainer.appendChild(homeElement);
}

// Fetch user details (ID and image_url)
async function getUserDetails() {
    try {
        const response = await fetch('/User');
        if (!response.ok) {
            throw new Error(`Failed to fetch user details: ${response.statusText}`);
        }

        return await response.json(); // Return the full user object
    } catch (error) {
        console.error("Error fetching user details:", error);
        return null; // Return null if unable to fetch user
    }
}

// Creates a new post (makes the server request)
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

        // Fetch the posts again to update the list
        await fetchPosts();
    } catch (error) {
        console.error("Error creating post:", error);
    }
}

function formatCreatedAt(createdAt) {
    const now = new Date();
    const postDate = new Date(createdAt);

    // Check if the post date is today
    const isToday =
        now.getFullYear() === postDate.getFullYear() &&
        now.getMonth() === postDate.getMonth() &&
        now.getDate() === postDate.getDate();

    if (isToday) {
        // Calculate how much time has passed
        const diffMs = now - postDate;
        const diffMinutes = Math.floor(diffMs / 60000); // Convert ms to minutes

        if (diffMinutes < 60) {
            return `${diffMinutes} minute(s) ago`;
        }

        const diffHours = Math.floor(diffMinutes / 60);
        return `${diffHours} hour(s) ago`;
    } else {
        // Format as YYYY-MM-DD for non-today dates
        return postDate.toISOString().split('T')[0];
    }
}

// Fetches posts from the server
async function fetchPosts() {
    try {
        const response = await fetch('/posts');

        // If the response is not OK, throw an error
        if (!response.ok) {
            throw new Error(`Failed to fetch posts: ${response.statusText}`);
        }

        const posts = await response.json();  // Parse the JSON response
        await render(posts);
    } catch (error) {
        const posts = [];
        await render(posts);
        console.error("Error fetching posts:", error);
        const postsContainer = document.getElementById("postsContainer");
        postsContainer.innerHTML = '<p>Error loading posts.</p>';
    }
}

// Initialize function to load posts on page load
export async function initialize() {
    await fetchPosts();  // Fetch and display posts on page load
}

// Call initialize when the DOM is ready
document.addEventListener("DOMContentLoaded", initialize);
