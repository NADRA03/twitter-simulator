export function render(Id) {
    return `
        <div id="userDetails">
            <!-- User details will be populated here -->
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
        const defaultImageUrl = "/assets/user.svg";
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

                <div id="followInfo">
            <span id="followingCount"><span class="white-text">${followingCount}</span> Following</span>
            <span id="followersCount"><span class="white-text" id="followersCountText">${followersCount}</span> Followers</span>
                </div>

                <button id="followButton">Follow</button>
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
                
                if (followResponse.ok) {
                    alert("Followed successfully!");
                    // Optionally, refresh follower count after following
                    const updatedFollowersCount = await fetch(`/followersCount?user_id=${Id}`)
                        .then(res => res.json())
                        .catch(() => 0) || 0;
                        document.getElementById("followersCountText").textContent = updatedFollowersCount;
                } else {
                    throw new Error("Could not follow the user");
                }
            } catch (error) {
                console.error("Error following user:", error);
                alert("Error following user");
            }
        });
        
    } catch (error) {
        console.error('Error loading user details:', error);
        document.getElementById('userDetails').innerHTML = '<p>Error loading user details.</p>';
    }
}
