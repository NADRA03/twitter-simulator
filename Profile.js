import { loadPage } from "./loader.js";
const defaultImageUrl = '/assets/user2.png';
// Renders user details dynamically without "Follow" and "Direct" buttons
export function render(user, followersCount, followingCount) {
    const defaultImageUrl = '/assets/user2.png';
    const defaultCoverUrl = '/assets/wallpaper.png';
    const userImageUrl = user.image_url || defaultImageUrl;
    const userBigImageUrl = user.big_image_url || defaultCoverUrl;

    return `
        <div id="userDetails">
            <div>
                <img id="bigImage" src="${userBigImageUrl}" onError="this.onerror=null; this.src='${defaultCoverUrl}'" />
                <h2 id="username">${user.username}</h2>
                <h2 id="name">${user.FirstName} ${user.LastName}</h2>
                <img id="profileImage" src="${userImageUrl}" onError="this.onerror=null; this.src='${defaultImageUrl}'" />
            </div>
            <div id="followInfo">
                    <span id="followingCount"><span class="white-text">${followingCount}</span> Following</span>
                    <span id="followersCount"><span class="white-text" id="followersCountText">${followersCount}</span> Followers</span>
            </div>
            <div id="postsContainer">
                <!-- Posts will be dynamically added here -->
            </div>
        </div>
    `;
}

// Initialize and load user details with posts
export async function initialize() {
    try {
        // Fetch user details
        const response = await fetch('/User');
        if (!response.ok) {
            throw new Error(`Failed to fetch user details: ${response.statusText}`);
        }

        const user = await response.json();
        const followersCount = await fetch(`/followersCount?user_id=${user.id}`)
        .then(res => res.json())
        .then(data => data.followersCount) // Get the count from the JSON response
        .catch(() => 0) || 0;

        const followingCount = await fetch(`/followingCount?user_id=${user.id}`)
        .then(res => res.json())
        .then(data => data.followingCount) // Assuming followingCount returns a similar JSON response
        .catch(() => 0) || 0;

        const app = document.getElementById('app');
        app.innerHTML = render(user, followersCount, followingCount);

        // Fetch and display posts for the user
        await fetchPosts(user.id);
    } catch (error) {
        console.error('Error initializing user details:', error);
        document.body.dataset.status = '401'; 
        loadPage(`401`);
    }
}

// Fetch posts for a user and render them
async function fetchPosts(userId) {
    try {
        const response = await fetch(`/posts/?user_id=${userId}`);
        if (!response.ok) {
            throw new Error(`Failed to fetch posts: ${response.statusText}`);
        }

        const posts = await response.json();
        const postsContainer = document.getElementById('postsContainer');
        postsContainer.innerHTML = ''; // Clear existing posts

        if (!Array.isArray(posts) || posts.length === 0) {
            postsContainer.innerHTML = '<p>No posts available</p>';
            return;
        }

        posts.forEach(post => {
            const postElement = document.createElement('div');
            postElement.classList.add('post');

            const formattedDate = formatCreatedAt(post['created-at']);
            const imageURL = post.image_url || '/assets/user.svg';

            postElement.innerHTML = `
                <p><img src="${imageURL}" onError="this.onerror=null; this.src='${defaultImageUrl}'" class="post-image"/> By ${post.Username} on ${formattedDate}</p>
                <p class="content">${post.content}</p>
            `;

            postElement.addEventListener('click', () => {
                // Redirect to the post's detailed page
                window.location.href = `/post/${post.id}`;
            });
            postsContainer.appendChild(postElement);
        });
    } catch (error) {
        console.error('Error fetching posts:', error);
        const postsContainer = document.getElementById('postsContainer');
        document.body.dataset.status = '401'; 
        loadPage(`401`);
    }
}

// Format post date for display
function formatCreatedAt(createdAt) {
    const now = new Date();
    const postDate = new Date(createdAt);

    if (
        now.getFullYear() === postDate.getFullYear() &&
        now.getMonth() === postDate.getMonth() &&
        now.getDate() === postDate.getDate()
    ) {
        const diffMs = now - postDate;
        const diffMinutes = Math.floor(diffMs / 60000);

        if (diffMinutes < 60) {
            return `${diffMinutes} minute(s) ago`;
        }
        return `${Math.floor(diffMinutes / 60)} hour(s) ago`;
    }
    return postDate.toISOString().split('T')[0];
}
