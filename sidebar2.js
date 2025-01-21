const defaultImageUrl = '/assets/user2.png';


export function render() {
    return `
        <div id="user-list">
        </div>
    `;
}

async function fetchAndRenderUsers() {
    try {
        const response = await fetch('/allusers');
        if (!response.ok) {
            throw new Error(`Failed to fetch users: ${response.status}`);
        }

        const users = await response.json();

        const userList = document.getElementById('user-list');
        userList.innerHTML = users.length > 0
            ? users.map(user => `
                <li class="user-item" data-id="${user.id}">
                  <img 
    src="${user.image_url || './assets/user2.png'}" 
    onerror="this.onerror=null; this.src='/assets/user2.png';"  
    class="user-avatar ${user.status === 'online' ? 'online' : ''}" 
/>
                    <div class="user-info">
                        <span class="username">${user.username}</span>
                    </div>
                </li>
            `).join('')
            : '<li>No users available.</li>';

        // Add click event listeners for each user
        const userItems = document.querySelectorAll('.user-item');
        userItems.forEach(item => {
            item.addEventListener('click', () => {
                const userId = item.getAttribute('data-id');
                redirectToChat(userId);
            });
        });

    } catch (error) {
        console.error('Error fetching users:', error);
        document.getElementById('user-list').innerHTML = '<li>Error loading users.</li>';
    }
}

// Function to redirect or create a chat
async function redirectToChat(userId) {
    try {
        const response = await fetch(`/chats/direct?id=${userId}`, {
            method: 'POST',
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
            window.location.href = result.redirectUrl; 
        } else {
            alert("Chat was created successfully!");
        }
    } catch (error) {
        console.error("Error handling chat creation or redirection:", error);
        alert("Error handling chat creation or redirection.");
    }
}

export async function initialize() {
await fetchAndRenderUsers()
}

