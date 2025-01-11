import { loadPage } from "./loader.js";
const defaultImageUrl = '/assets/user2.png';


export function render(Id) {
    return `
        <div id="userDetails">
            <!-- User details and posts will be populated here -->
            <div id="postsContainer"></div>
        </div>
    `;
}

export async function initialize(Id) {
    try {
        // Fetch the user details from the server
        const response = await fetch(`/user_profile/${Id}`);
        
        if (!response.ok) {
            throw new Error('User not found');
        }

        // Parse the JSON response
        const user = await response.json();

        // Default image URL if no image is provided
        const defaultImageUrl = "/assets/user2.png";
        const defaultCoverUrl = "/assets/wallpaper.png";

        // Use the default image if no image URL is found
        const userImageUrl = user.image_url || defaultImageUrl;
        const userBigImageUrl = user.big_image_url || defaultCoverUrl;

        // Fetch follower and following counts, defaulting to 0 if undefined
        const followersCount = await fetch(`/followersCount?user_id=${Id}`)
            .then(res => res.json())
            .then(data => data.followersCount) // Get the count from the JSON response
            .catch(() => 0) || 0;

        const followingCount = await fetch(`/followingCount?user_id=${Id}`)
            .then(res => res.json())
            .then(data => data.followingCount) // Assuming followingCount returns a similar JSON response
            .catch(() => 0) || 0;

        // Populate the user details in the HTML
        const userDetailsContainer = document.getElementById('userDetails');
        userDetailsContainer.innerHTML = `
            <div>
                <img id="bigImage" src="${userBigImageUrl}" onError="this.onerror=null; this.src='${defaultCoverUrl}'" />
                <h2 id="username">${user.username}</h2>
                <h2 id="name">${user.FirstName} ${user.LastName}</h2>
                <img id="profileImage" src="${userImageUrl}" onError="this.onerror=null; this.src='${defaultImageUrl}'" />
                <button id="directButton" class="direct-btn">
                    <img id="direct" src="/assets/direct.svg" alt="Direct message" />
                </button>
                <div id="followInfo">
                    <span id="followingCount"><span class="white-text">${followingCount}</span> Following</span>
                    <span id="followersCount"><span class="white-text" id="followersCountText">${followersCount}</span> Followers</span>
                </div>
                <button id="followButton">Follow</button>
            </div>
            <div id="postsContainer">
                <!-- Posts will be dynamically added here -->
            </div>
        `;

        // Add event listener for the Follow button
        document.getElementById('followButton').addEventListener('click', async () => {
            try {
                // Send followed_id in the URL as a query parameter
                const followResponse = await fetch(`/follow?followed_id=${Id}`, { 
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });
                
                if (followResponse.ok && followResponse.status === 201) {
                    // Optionally, refresh follower count after following
                const updatedFollowersCount = await fetch(`/followersCount?user_id=${Id}`)
                    .then(res => res.json())
                    .catch(() => ({ followersCount: 0 })); // Default to 0 if fetch fails
                
                // Extract the `followersCount` property
                const count = updatedFollowersCount.followersCount || 0;
                
                // Update the text content of the element
                document.getElementById("followersCountText").textContent = count;
                } else {
                document.body.dataset.status = '401'; 
                loadPage(`401`);
                }
            } catch (error) {
                console.error("Error following user:", error);
            }
        });

        document.getElementById('directButton').addEventListener('click', async () => {
            try {
                const response = await fetch(`/chats/direct?id=${Id}`, {
                    method: 'POST', // Use POST for creating or directing to a chat
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });
        
                if (!response.ok) {
                    const errorText = await response.text();
                    throw new Error(errorText);
                }
        
                const result = await response.json();
                console.log("Chat redirected/created:", result);
        
                if (result.redirectUrl) {
                    window.location.href = result.redirectUrl; // Redirect to the chat page
                } else {
                    alert("Chat was created successfully!");
                }
            } catch (error) {
                console.error("Error handling chat creation or redirection:", error);
                alert("Error handling chat creation or redirection.");
            }
        });
        

        // Fetch posts for the user and display them
        await fetchPosts(Id);

    } catch (error) {
        console.error('Error loading user details:', error);
        document.getElementById('userDetails').innerHTML = '<p>Error loading user details.</p>';
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

// Fetches posts from the server and displays them
async function fetchPosts(Id) {
    try {
        const response = await fetch(`/posts/?user_id=${Id}`);
        if (!response.ok) {
            throw new Error(`Failed to fetch posts: ${response.statusText}`);
        }

        const posts = await response.json();  // Parse the JSON response

        const postsContainer = document.getElementById("postsContainer");
        postsContainer.innerHTML = ''; // Clear previous posts

        if (!Array.isArray(posts)) {
            console.error("Expected an array, but got:", posts);
            posts = [];
        }

        if (posts.length === 0) {
            const noPostsMessage = document.createElement("p");
            noPostsMessage.textContent = "No posts available";
            postsContainer.appendChild(noPostsMessage);
        } else {
            posts.forEach(post => {
                const postElement = document.createElement("div");
                postElement.classList.add("post");

                const formattedDate = formatCreatedAt(post["created-at"]);

                const imageURL = post.image_url || defaultImageUrl;
                postElement.innerHTML = `
                    <p><img src="${imageURL}" onerror="this.onerror=null; this.src='${defaultImageUrl}'" class="post-image"/> By ${post.Username} on ${formattedDate}</p>
                    <p class="content">${post.content}</p>
                `;
                postElement.addEventListener('click', () => {
                    // Redirect to the post's detailed page
                    window.location.href = `/post/${post.id}`;
                });
                const postImage = postElement.querySelector(".post-image");
                postImage.addEventListener('click', (e) => {
                    e.stopPropagation(); // Prevent triggering the post click event
                    window.location.href = `/a_profile/${post["user-id"]}`;
                });
                postsContainer.appendChild(postElement);
            });
        }

    } catch (error) {
        console.error("Error fetching posts:", error);
        const postsContainer = document.getElementById("postsContainer");
        postsContainer.innerHTML = '<p class="no-posts-message">No posts yet.</p>';
    }
}