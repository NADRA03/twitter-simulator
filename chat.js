let socket;
let chatId;
let globalChatDetails = null;
let globalUsers = [];
let globalMessages = [];
let displayedMessageIds = new Set();
let userLookup = {}; // Create a lookup for user details

// Function to populate the user lookup
function populateUserLookup(users) {
    userLookup = users.reduce((lookup, user) => {
        lookup[user.id] = {
            username: user.username,
            imageUrl: user.image_url ? user.image_url.String : '/assets/user.svg',
            onlineStatus: 'Offline'
        };
        return lookup;
    }, {});
}

// Function to fetch online status for each user
function fetchOnlineStatus() {
    const statusPromises = globalUsers.map(async (user) => {
        try {
            const response = await fetch(`/online-status?userID=${user.id}`);
            const status = await response.text(); // Expecting "Online" or "Offline"
            userLookup[user.id].onlineStatus = status;
        } catch (error) {
            console.error(`Error fetching online status for user ${user.id}:`, error);
        }
    });
    Promise.all(statusPromises);
}

function updateUserList() {
    const userListHTML = globalUsers.map(user => `
        <div class="user-item">
            <img src="${user.image_url.String || '/assets/user.svg'}" alt="${user.username}'s profile picture" class="user-image" />
            <span class="user-name">${user.username}</span>
            <span class="user-status ${userLookup[user.id].onlineStatus === 'Online' ? 'online' : 'offline'}">
                ${userLookup[user.id].onlineStatus}
            </span>
        </div>
    `).join('');

    document.getElementById("userList").innerHTML = userListHTML;
}

// Function to periodically fetch online status and update user list
function startUserStatusUpdate() {
    fetchOnlineStatus();
    updateUserList(); 

    setInterval(async () => {
        fetchOnlineStatus();
        updateUserList();
    }, 15000); 
}


export function render(chatId) {
    const uniqueMessageIds = new Set();
    
    const uniqueMessages = globalMessages.filter(message => {
        if (uniqueMessageIds.has(message.message_id)) {
            return false;
        }
        uniqueMessageIds.add(message.message_id);
        return true;
    });

    // Create a promise array for checking user details for all messages
    const messagePromises = uniqueMessages.map(async (message) => {
        const user = userLookup[message.user_id] || { username: 'Unknown', imageUrl: '/assets/user.svg' };
        
        // Call fetchUserDetails to check if the message should go to the right
        const isCurrentUserMessage = await fetchUserDetails(message.user_id);
        
        // Create the message element
        const messageElement = `
            <div class="message-item ${isCurrentUserMessage ? 'right' : ''}">
                <img src="${user.imageUrl || '/assets/user.svg'}"  class="user-image" />
                <span class="message-username">${user.username}:</span>
                <span class="message-text">${message.message_text.String}</span>
                ${message.image_url ? `<img src="${message.image_url.String}" class="message-image" />` : ''}
            </div>
        `;
        
        return messageElement; // Return the message element
    });

    // Wait for all promises to resolve and then join the results
    Promise.all(messagePromises)
        .then(messageElements => {
            document.getElementById("chatMessages").innerHTML = messageElements.join('');
            document.getElementById("chatMessages").scrollTop = document.getElementById("chatMessages").scrollHeight;
        })
        .catch(error => {
            console.error("Error rendering messages:", error);
        });

    return `
        <div class="chatContainer">
            <button id="infoButton" class="info-button"><img src="/assets/info.svg"  class="user-image" /></button>
            <div id="chatMessages">
                ${uniqueMessages.length > 0 ? '' : '<p>No messages yet.</p>'}
            </div>
            <div id="send">
                <input type="text" id="messageInput" placeholder="Type your message..." />
                <button id="sendMessageBtn">Send</button>
            </div>
            <div id="infoBox" class="info-box" style="display: none;">
                <button id="closeInfoBox" class="close-button">X</button>
                <br>
                <p>Information about the chat...</p>
                <button id="addPeopleButton">Invite</button>
                <div id="addPeopleSection" style="display: none;">
                    <p>Add people to the chat:</p>
                    <input type="text" id="addPeopleInput" />
                    <button  style="display: none;" id="confirmAddPeopleBtn">Add</button>
                </div>
                <div id="userSearchResults"></div>

                <div id="userList" class="user-list">
                    <h3>Participants</h3>
                    ${globalUsers.map(user => `
                        <div class="user-item">
                            <img src="${user.image_url.String || '/assets/user.svg'}" alt="${user.username}'s profile picture" class="user-image" />
                            <span class="user-name">${user.username}</span>
                            <span class="user-status ${userLookup[user.id].onlineStatus === 'Online' ? 'online' : 'offline'}">
                                ${userLookup[user.id].onlineStatus}
                            </span>
                        </div>
                    `).join('')}
                </div>
            </div>
        </div>
    `;
}

function enableAllButtons() {
    const allButtons = document.querySelectorAll("button");
    allButtons.forEach(button => {
        button.disabled = false; // Enable all buttons
    });
}

async function fetchUserDetails(user_id) {
    try {
        const response = await fetch('/User');
        
        if (!response.ok) {
            throw new Error(`Failed to fetch user details: ${response.statusText}`);
        }
        
        const user = await response.json(); // Assuming the response returns user details in JSON format
        
        // Return true if the user's id matches the message_id, false otherwise
        return user.id === user_id;
    } catch (error) {
        console.error(error);
        return false; // Return false if there was an error
    }
}

async function displayMessage(messageData) {
    // Prevent displaying the same message multiple times
    if (displayedMessageIds.has(messageData.message_id)) {
        return;
    }
    displayedMessageIds.add(messageData.message_id);

    const chatMessages = document.getElementById("chatMessages");
    const user = userLookup[messageData.user_id] || { username: 'Unknown', imageUrl: '/assets/user.svg' };

    const messageElement = document.createElement("div");
    messageElement.className = "message-item";
    messageElement.innerHTML = `
        <img src="${user.imageUrl || '/assets/user.svg'}"  class="user-image" />
        <span class="message-username">${user.username}:</span>
        <span class="message-text">${messageData.message_text}</span>
        ${messageData.image_url ? `<img src="${messageData.image_url}" alt="Attached image" class="message-image" />` : ''}
    `;

    // Call fetchUserDetails to check if the message should go to the right
    const isCurrentUserMessage = await fetchUserDetails(messageData.user_id);
    console.log("isCurrentUserMessage:", isCurrentUserMessage); 

    // If the user is the same as the logged-in user, move the message to the right
    if (isCurrentUserMessage) {
        messageElement.classList.add('right');
    }

    chatMessages.appendChild(messageElement);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}

function initializeWebSocket(chatId) {
    const socketUrl = `ws://localhost:8088/ws?chat_id=${chatId}`;
    socket = new WebSocket(socketUrl);

    socket.onopen = () => {
        console.log("Connected to WebSocket for chat ID:", chatId);
    };

    socket.onmessage = (event) => {
        const messageData = JSON.parse(event.data);
        if (!messageData.message_id) {
            console.warn("Message without ID received:", messageData);
            return;
        }
        displayMessage(messageData);
    };

    socket.onclose = () => {
        console.log("WebSocket connection closed.");
        clearDisplayedMessages(); // Clear displayed messages on close to prepare for reconnection
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
}

// Define the clearDisplayedMessages function to reset the displayed messages
function clearDisplayedMessages() {
    displayedMessageIds.clear();
}

function sendMessage(chatId) {
    try {
        if (socket.readyState === WebSocket.OPEN) {
            const messageInput = document.getElementById("messageInput");
            const messageText = messageInput.value;
            const imageInput = document.getElementById("imageInput");
            const imageUrl = imageInput ? imageInput.value : "";

            if (messageText) {
                const messageData = {
                    chatId: chatId,
                    messageText: messageText,
                    imageUrl: imageUrl
                };

                console.log("Preparing to send message:", messageData);

                try {
                    socket.send(JSON.stringify({
                        action: "add_message",
                        data: messageData
                    }));
                    console.log("Message successfully sent:", messageData);

                    // Clear inputs after sending
                    messageInput.value = "";
                    if (imageInput) imageInput.value = "";
                } catch (sendError) {
                    console.error("Failed to send message:", sendError);
                    alert("There was an error sending your message. Please try again.");
                }
            } else {
                console.warn("Message text is empty, nothing to send.");
            }
        } else {
            console.error("WebSocket is not open. Cannot send message.");
            alert("Connection to the server is closed. Please check your connection.");
        }
    } catch (error) {
        console.error("An error occurred in sendMessage:", error);
    }
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
        const userElement = document.createElement("div");     ////////////////////
        userElement.className = "user";
        userElement.innerHTML = `
            <br>
            <img src="${user.image_url || '/assets/user.svg'}" style="width: 40px; height: 40px; border-radius: 20px; cursor: pointer;"   onerror="this.onerror=null; this.src='./assets/user.svg';"/> 
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
        inviteUserToChat(user.id); // Call inviteUserToChat without passing chatId
        document.body.removeChild(confirmationBlock); // Remove after confirming
        enableAllButtons(); // Re-enable all other buttons
    };

    // No button functionality
    confirmationBlock.querySelector("#noInviteBtn").onclick = () => {
        document.body.removeChild(confirmationBlock); // Remove when clicking No
        enableAllButtons(); // Re-enable all other buttons
    };
}

// Function to handle inviting a user to the chat
function inviteUserToChat(userId) {
    fetch(`/chats/addUser`, {
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

async function fetchChatDetails(chatId) {
    return fetch(`/chat_details?chat_id=${chatId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            globalChatDetails = data.chat_details;
            globalUsers = data.users;
            globalMessages = data.messages;
            populateUserLookup(globalUsers); // Populate user lookup here
            fetchOnlineStatus(); 
        })
        .catch(error => {
            console.error('Error fetching chat details:', error);
        });
}

export function initialize(chatIdParam) {
    chatId = chatIdParam;
    console.log("Chat ID:", chatId); // Add this line
    const chatContent = render(chatId);
    document.getElementById("sidebar2").innerHTML = chatContent;
  
    initializeWebSocket(chatId);
    
    fetchChatDetails(chatId).then(() => {
        console.log("Chat Details:", globalChatDetails);
        console.log("Users:", globalUsers);
        console.log("Messages:", globalMessages);
    });

    startUserStatusUpdate();
    
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
    const chatIdParam = pathSegments[pathSegments.length - 1]; // Get the last segment as `chatId`
    if (chatIdParam) {
        initialize(chatIdParam);  // Initialize with the extracted `chatId`
    } else {
        console.error("chatId not found in the URL.");
    }
});
