let socket; 

export function render(chatId) {
    return `
    <div class="chatContainer">
        <button id="infoButton" class="info-button">Info</button>
        <div id="chatMessages"></div>
        <div id="send">
            <input type="text" id="messageInput" placeholder="Type your message..." />
            <button id="sendMessageBtn">Send</button>
        </div>
        <div id="infoBox" class="info-box" style="display: none;">
            <button id="closeInfoBox" class="close-button">X</button>
            <br>
            <p>Information about the chat...</p>
            <button id="addPeopleButton">Invite</button> <!-- Add People button inside infoBox -->
            <div id="addPeopleSection" style="display: none;">
                <p>Add people to the chat:</p>
                <input type="text" id="addPeopleInput" />
                <button style="display: none;" id="confirmAddPeopleBtn">Add</button>
            </div>
            <div id="userSearchResults"></div> <!-- Container for search results -->
        </div>
    </div>`;
}

function initializeWebSocket(chatId) {
    const socketUrl = `ws://localhost:8080/chat/${chatId}`; 
    socket = new WebSocket(socketUrl);

    socket.onopen = () => {
        console.log("Connected to WebSocket for chat ID:", chatId);
    };

    socket.onmessage = (event) => {
        const messageData = JSON.parse(event.data);
        displayMessage(messageData);
    };

    socket.onclose = () => {
        console.log("WebSocket connection closed.");
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
}

function sendMessage(chatId) {
    const messageInput = document.getElementById("messageInput");
    const messageText = messageInput.value;
    if (messageText) {
        const messageData = {
            chatId: chatId,
            userId: 1, 
            messageText: messageText
        };
        socket.send(JSON.stringify(messageData));
        messageInput.value = ""; 
    }
}

function displayMessage(messageData) {
    const chatMessages = document.getElementById("chatMessages");
    const messageElement = document.createElement("div");
    messageElement.textContent = `${messageData.userId}: ${messageData.messageText}`;
    chatMessages.appendChild(messageElement);
}

function searchUsers(searchTerm) {
    if (searchTerm) {
        fetch(`/search-users?term=${encodeURIComponent(searchTerm)}`)
            .then(response => response.json())
            .then(users => displaySearchResults(users))
            .catch(error => console.error('Error fetching users:', error));
    } else {
        // Clear search results if input is empty
        document.getElementById("userSearchResults").innerHTML = "";
    }
}

function displaySearchResults(users) {
    const userSearchResults = document.getElementById("userSearchResults");
    userSearchResults.innerHTML = ""; // Clear previous results

    // Ensure users is an array
    if (!Array.isArray(users)) {
        users = []; // If users is not an array (e.g., null), set it to an empty array
    }

    if (users.length === 0) {
        // Display a message when no users are found
        const noUsersElement = document.createElement("div");
        noUsersElement.id = "noUsersMessage"; // Add an ID for the no users message
        noUsersElement.textContent = "No users found.";
        userSearchResults.appendChild(noUsersElement);
        return; // Exit the function early
    }

    users.forEach((user) => {
        const userElement = document.createElement("div");
        userElement.className = "user";
        userElement.innerHTML = `
            <br>
            <img src="${user.image_url}" style="width: 40px; height: 40px; border-radius: 20px; cursor: pointer;"/> 
            <span style="cursor: pointer;">${user.username}</span>
        `;

        // Click event for inviting the user
        userElement.onclick = () => showInviteConfirmation(user);

        userSearchResults.appendChild(userElement);
    });
}

// Function to show the invite confirmation block
function showInviteConfirmation(user) {
    // Check if there's an existing confirmation block and remove it
    let existingBlock = document.querySelector(".confirmation-block");
    if (existingBlock) {
        document.body.removeChild(existingBlock);
    }

    // Disable all buttons on the page
    const allButtons = document.querySelectorAll("button");
    allButtons.forEach(button => {
        if (!button.classList.contains("confirmation-btn")) {
            button.disabled = true;
        }
    });

    const confirmationBlock = document.createElement("div");
    confirmationBlock.className = "confirmation-block";
    confirmationBlock.innerHTML = `
        <p>Do you want to invite ${user.username} to the chat?</p>
        <button id="yesInviteBtn" class="confirmation-btn">Yes</button>
        <button id="noInviteBtn" class="confirmation-btn">No</button>
    `;
    document.body.appendChild(confirmationBlock);

    // Yes button functionality
    confirmationBlock.querySelector("#yesInviteBtn").onclick = () => {
        inviteUserToChat(user.id, chatId);
        document.body.removeChild(confirmationBlock); // Remove after confirming
        enableAllButtons(); // Re-enable all other buttons
    };

    // No button functionality
    confirmationBlock.querySelector("#noInviteBtn").onclick = () => {
        document.body.removeChild(confirmationBlock); // Remove when clicking No
        enableAllButtons(); // Re-enable all other buttons
    };
}

function enableAllButtons() {
    const allButtons = document.querySelectorAll("button");
    allButtons.forEach(button => {
        button.disabled = false;
    });
}

// Function to handle inviting a user to the chat
function inviteUserToChat(userId, chatId) {
    fetch(`/add-user-to-chat`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
            chat_id: chatId,
            user_id: userId,
            role: 'participant', // or whatever role you want to assign
        }),
    })
    .then(response => {
        if (response.ok) {
            console.log(`User ${userId} invited to chat successfully!`);
        } else {
            console.error('Error inviting user:', response.statusText);
        }
    })
    .catch(error => console.error('Fetch error:', error));
}



export function initialize(chatId) {
    const chatContent = render(chatId);
    document.getElementById("sidebar2").innerHTML = chatContent;

    initializeWebSocket(chatId);

    const sendMessageBtn = document.getElementById("sendMessageBtn");
    if (sendMessageBtn) {
        sendMessageBtn.onclick = () => {
            sendMessage(chatId); 
        };
    }

    const messageInput = document.getElementById("addPeopleInput");
    messageInput.addEventListener("input", () => {
        const searchTerm = messageInput.value;
        searchUsers(searchTerm);
    });

    const infoButton = document.getElementById("infoButton");
    const infoBox = document.getElementById("infoBox");
    const closeInfoBox = document.getElementById("closeInfoBox");
    const addPeopleButton = document.getElementById("addPeopleButton");
    const addPeopleSection = document.getElementById("addPeopleSection");
    const confirmAddPeopleBtn = document.getElementById("confirmAddPeopleBtn");

    // Show the info box when the "Info" button is clicked
    infoButton.onclick = () => {
        infoBox.style.display = "block";
    };

    // Hide the info box when the "X" button is clicked
    closeInfoBox.onclick = () => {
        infoBox.style.display = "none";
        addPeopleSection.style.display = "none"; // Hide the add people section when closing infoBox
    
        // Clear search results when closing the info box
        document.getElementById("userSearchResults").innerHTML = "";
    };

    // Show the add people section when the "Add People" button is clicked
    addPeopleButton.onclick = () => {
        addPeopleSection.style.display = "block"; // Show the add people section
    };

    // Add functionality to confirm adding people
    confirmAddPeopleBtn.onclick = () => {
        const emailOrUsername = document.getElementById("addPeopleInput").value;
        if (emailOrUsername) {
            console.log(`Adding: ${emailOrUsername}`); // Add your logic to handle adding people here
            // Optionally clear the input
            document.getElementById("addPeopleInput").value = "";
        }
    };
}

document.addEventListener("DOMContentLoaded", () => {
    // Extract `chatId` from the URL path
    const pathSegments = window.location.pathname.split('/');
    const chatId = pathSegments[pathSegments.length - 1]; // Get the last segment as `chatId`

    if (chatId) {
        initialize(chatId);  // Initialize with the extracted `chatId`
    } else {
        console.error("chatId not found in the URL.");
    }
});
