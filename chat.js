let socket; 
let chatId;
let globalChatDetails = null;
let globalUsers = [];
let globalMessages = [];


export function render(chatId) {
    console.log("Rendering users:", globalUsers);  // Log globalUsers to confirm data
    return `
        <div class="chatContainer">
            <button id="infoButton" class="info-button">Info</button>

           <!-- Chat messages section -->
            <div id="chatMessages">
                ${globalMessages && globalMessages.length > 0 ? globalMessages.map(message => `
                    <div class="message-item">
                        <span class="message-username">${message.username}:</span>
                        <span class="message-text">${message.message_text}</span>
                        ${message.image_url ? `<img src="${message.image_url}" alt="Attached image" class="message-image" />` : ''}
                    </div>
                `).join('') : '<p>No messages yet.</p>'}
            </div>

            <div id="send">
                <input type="text" id="messageInput" placeholder="Type your message..." />
                <input type="text" id="imageInput" style="display: none;" placeholder="Type your message..." />
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
                    <button style="display: none;" id="confirmAddPeopleBtn">Add</button>
                </div>
                <div id="userSearchResults"></div>

                <!-- User list section -->
                <div id="userList" class="user-list">
                <h3>Participants</h3>
                ${(() => {
                    let userHtml = ''; // Initialize an empty string for the HTML
                    globalUsers.forEach(user => {
                        userHtml += `
                            <div class="user-item">
                                <img src="${user.image_url.String || 'placeholder-image-url.jpg'}" alt="${user.username}'s profile picture" class="user-image" />
                                <span class="user-name">${user.username}</span>
                            </div>
                        `;
                    });
                    return userHtml; // Return the complete HTML string
                })()}
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

function displayMessage(messageData) {
    const chatMessages = document.getElementById("chatMessages");
    const messageElement = document.createElement("div");
    
    messageElement.innerHTML = `
        <span class="message-username">${messageData.user_id}:</span>
        <span class="message-text">${messageData.message_text}</span>
        ${messageData.image_url ? `<img src="${messageData.image_url}" alt="Attached image" class="message-image" />` : ''}
    `;
    
    chatMessages.appendChild(messageElement);
}

function initializeWebSocket(chatId) {
    const socketUrl = `ws://localhost:8088/ws?chat_id=${chatId}`; 
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

function fetchChatDetails(chatId) {
    return fetch(`/chat_details?chat_id=${chatId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // Set global variables with the response data
            globalChatDetails = data.chat_details;
            globalUsers = data.users;
            globalMessages = data.messages;
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
    render(chatId);
    
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
