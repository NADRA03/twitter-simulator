// Renders user details dynamically
export function render(user) {
    const defaultImageUrl = 'https://i.abcnewsfe.com/a/7d5c53c5-f5a7-40a6-9c99-648316b9f7a2/Elon-Musk-1-rt-jm-241016_1729107345064_hpMain.jpg'; 
    const userImageUrl = user.image_url || defaultImageUrl; // Set to user image or default

    return `
        <div>
            <h1>User Details</h1>
            <p><strong>Username:</strong> ${user.username}</p>
            <p><strong>Email:</strong> ${user.email}</p>
            <img src="${userImageUrl}" alt="User Image" style="max-width: 100px;" onerror="this.onerror=null; this.src='${defaultImageUrl}'">
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