// Renders user details dynamically
export function render(user) {
    const defaultImageUrl = '/assets/user.svg'; 
    const userImageUrl = user.image_url || defaultImageUrl; // Set to user image or default
    const defaultCoverUrl = "/assets/wallpaper.png";
    const userBigImageUrl = user.big_image_url || defaultCoverUrl;
    return `
        <div id="userDetails">
            <div>
            <img id="bigImage" src="${userBigImageUrl}" onError="this.onerror=null; this.src='${defaultCoverUrl}'" />
            <h2 id="username"> ${user.username}</h2>
            <h2 id="name">${user.FirstName} ${user.LastName}</h2>
            <img id="profileImage" src="${userImageUrl}" alt="User Image" onerror="this.onerror=null; this.src='${defaultImageUrl}'">
            </div>
        </div>
    `;
}

// Fetches user details from the server
async function fetchUserDetails() {
    try {
        const response = await fetch('/User');
        if (!response.ok) {
            throw new Error(`Failed to fetch user details: ${response.statusText}`);
        }
        
        const user = await response.json();
        console.log(user);
        const app = document.getElementById('app');
        app.innerHTML = render(user);
    } catch (error) {
        console.error('Error fetching user details:', error);
        document.getElementById('app').innerHTML = '<p>Error loading user details.</p>';
    }
}

// Initialize function to load user details on page load
export async function initialize() {
    await fetchUserDetails();  // Fetch and display user details on page load
}