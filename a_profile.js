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
        const defaultImageUrl = "https://i.abcnewsfe.com/a/7d5c53c5-f5a7-40a6-9c99-648316b9f7a2/Elon-Musk-1-rt-jm-241016_1729107345064_hpMain.jpg";

        // Function to handle image loading error
        const handleImageError = (e) => {
            e.target.onerror = null; // Prevent infinite loop if the default image also fails
            e.target.src = defaultImageUrl;
            e.target.alt = "Error loading image";
        };

        // Use the default image if no image URL is found
        const userImageUrl = user.image_url || defaultImageUrl;
        const userBigImageUrl = user.big_image_url || defaultImageUrl;

        // Populate the user details in the HTML
        const userDetailsContainer = document.getElementById('userDetails');
        userDetailsContainer.innerHTML = `
            <div>
             <img id="bigImage" src="${userBigImageUrl}"  onError="this.onerror=null; this.src='${defaultImageUrl}'" />
                <h2 id="username">${user.username}</h2>
                <h2 id="name">${user.FirstName}${user.LastName}</h2>
                <img id="profileImage" src="${userImageUrl}" onError="this.onerror=null; this.src='${defaultImageUrl}'" />
            </div>
        `;
    } catch (error) {
        console.error('Error loading user details:', error);
        document.getElementById('userDetails').innerHTML = '<p>Error loading user details.</p>';
    }
}