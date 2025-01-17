let socket;
let chatId;
let globalChatDetails = null;
let globalUsers = [];
let globalMessages = [];
let displayedMessageIds = new Set();
let userLookup = {}; 
let loadingMessages = false; ; 


function populateUserLookup(users) {
    userLookup = users.reduce((lookup, user) => {
        lookup[user.id] = {
            username: user.username,
            imageUrl: user.image_url ? user.image_url.String : '/assets/user2.png',
            onlineStatus: 'Offline'
        };
        return lookup;
    }, {});
}


function fetchOnlineStatus() {
    const statusPromises = globalUsers.map(async (user) => {
        try {
            const response = await fetch(`/online-status?userID=${user.id}`);
            const status = await response.text(); 
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
            <img src="${user.image_url.String || '/assets/user2.png'}" alt="${user.username}'s profile picture" class="user-image" />
            <span class="user-name">${user.username}</span>
            <span class="user-status ${userLookup[user.id].onlineStatus === 'Online' ? 'online' : 'offline'}">
                ${userLookup[user.id].onlineStatus}
            </span>
        </div>
    `).join('');

    document.getElementById("userList").innerHTML = userListHTML;
}


function startUserStatusUpdate() {
    fetchOnlineStatus();
    updateUserList(); 
    setInterval(async () => {
        fetchOnlineStatus();
        updateUserList();
    }, 15000); 
}


export function render(chatId) {
    console.log("Rendering chat with ID:", chatId);
    console.log("Global messages:", globalMessages);
    console.log("User list:", globalUsers);
    if (!globalMessages) {
        globalMessages = [];
    }

    const uniqueMessageIds = new Set();
    const uniqueMessages = globalMessages.filter(message => {
        if (uniqueMessageIds.has(message.message_id.Int64)) {
            return false;
        }
        uniqueMessageIds.add(message.message_id.Int64);
        return true;
    });

    const messagePromises = uniqueMessages.map(async (message) => {
        const userId = message.user_id.Valid ? message.user_id.Int64 : null;
        const user = userId && userLookup[userId] || { username: 'Unknown', imageUrl: '/assets/user2.png' };
        const isCurrentUserMessage = await fetchUserDetails(userId);
        const messageText = message.message_text.Valid ? message.message_text.String : 'Invalid message text';
        const imageUrl = message.image_url.Valid ? message.image_url.String : null;
        const messageElement = `
            <div class="message-item ${isCurrentUserMessage ? 'right' : ''}">
                <img src="${user.imageUrl || '/assets/user2.png'}" class="user-image" />
                <span class="message-username">${user.username}:</span>
                <span class="message-text">${messageText}</span>
                ${imageUrl ? `<img src="${imageUrl}" class="message-image" />` : ''}
            </div>
        `;
        
        return messageElement;
    });


    Promise.all(messagePromises)
        .then(messageElements => {
            const chatMessagesElement = document.getElementById("chatMessages");
            chatMessagesElement.innerHTML = messageElements.join('');
            // Scroll to the bottom to ensure latest message is visible
            // chatMessagesElement.scrollTop = chatMessagesElement.scrollHeight;
        })
        .catch(error => {
            console.error("Error rendering messages:", error);
        });

    return `
        <div class="chatContainer">
            <button id="infoButton" class="info-button"><img src="/assets/info.svg" class="user-image" /></button>
            <div id="chatMessages">
                ${uniqueMessages.length === 0 ? '<p>No messages yet.</p>' : ''}
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
                    <button style="display: none;" id="confirmAddPeopleBtn">Add</button>
                </div>
                <div id="userSearchResults"></div>

                <div id="userList" class="user-list">
                    <h3>Participants</h3>
                    ${globalUsers.map(user => `
                        <div class="user-item">
                            <img src="${user.image_url.String || '/assets/user2.png'}" alt="${user.username}'s profile picture" class="user-image" />
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
        button.disabled = false; 
    });
}


async function fetchUserDetails(user_id) {
    try {
        const response = await fetch('/User');
        if (!response.ok) {
            throw new Error(`Failed to fetch user details: ${response.statusText}`);
        } 
        const user = await response.json(); 
        return user.id === user_id;
    } catch (error) {
        console.error(error);
        return false; 
    }
}


async function displayMessage(messageData) {
    if (displayedMessageIds.has(messageData.message_id)) {
        return;
    }
    displayedMessageIds.add(messageData.message_id);

    const chatMessages = document.getElementById("chatMessages");
    const user = userLookup[messageData.user_id] || { username: 'Unknown', imageUrl: '/assets/user2.png' };

    const messageElement = document.createElement("div");
    messageElement.className = "message-item";
    messageElement.innerHTML = `
        <img src="${user.imageUrl || '/assets/user2.png'}"  class="user-image" />
        <span class="message-username">${user.username}:</span>
        <span class="message-text">${messageData.message_text}</span>
        ${messageData.image_url ? `<img src="${messageData.image_url}" alt="Attached image" class="message-image" />` : ''}
    `;
    const isCurrentUserMessage = await fetchUserDetails(messageData.user_id);
    console.log("isCurrentUserMessage:", isCurrentUserMessage); 

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
        clearDisplayedMessages(); 
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
}


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
        document.getElementById("userSearchResults").innerHTML = "";
    }
}


function displaySearchResults(users) {
    const userSearchResults = document.getElementById("userSearchResults");
    userSearchResults.innerHTML = "";


    if (!Array.isArray(users)) {
        users = []; 
    }

    if (users.length === 0) {
        const noUsersElement = document.createElement("div");
        noUsersElement.id = "noUsersMessage"; 
        noUsersElement.textContent = "No users found.";
        userSearchResults.appendChild(noUsersElement);
        return; 
    }

    users.forEach((user) => {
        const userElement = document.createElement("div");     
        userElement.className = "user";
        userElement.innerHTML = `
            <br>
            <div class="searched">
            <img src="${user.image_url || '/assets/user2.png'}" style="width: 40px; height: 40px; border-radius: 20px; cursor: pointer;"   onerror="this.onerror=null; this.src='/assets/user2.png';"/> 
            <p class="search-text" style="cursor: pointer;">${user.username}</p>
            <div>
        `;
        userElement.onclick = () => showInviteConfirmation(user);

        userSearchResults.appendChild(userElement);
    });
}


function showInviteConfirmation(user) {
    let existingBlock = document.querySelector(".confirmation-block");
    if (existingBlock) {
        document.body.removeChild(existingBlock);
    }

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

    confirmationBlock.querySelector("#yesInviteBtn").onclick = () => {
        inviteUserToChat(user.id); 
        
        document.body.removeChild(confirmationBlock); 
        enableAllButtons(); 
        location.reload()
    };

    confirmationBlock.querySelector("#noInviteBtn").onclick = () => {
        document.body.removeChild(confirmationBlock); 
        enableAllButtons(); 
    };
}


function inviteUserToChat(userId) {
    fetch(`/chats/addUser`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
            chat_id: chatId,
            user_id: userId,
            role: 'participant', 
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
            populateUserLookup(globalUsers); 
            fetchOnlineStatus(); 
        })
        .catch(error => {
            console.error('Error fetching chat details:', error);
        });
}


 export async function initialize(chatIdParam) {
    chatId = chatIdParam;
    console.log("Chat ID:", chatId); 
    
    initializeWebSocket(chatId);
    await fetchChatDetails(chatId);
    
    const chatContent = render(chatId);
    document.getElementById("sidebar2").innerHTML = chatContent;
    
    console.log("Chat Details:", globalChatDetails);
    console.log("Users:", globalUsers);
    console.log("Messages:", globalMessages);

    startUserStatusUpdate();
    
    const sendMessageBtn = document.getElementById("sendMessageBtn");
    if (sendMessageBtn) {
        sendMessageBtn.onclick = () => {
            sendMessage(chatId); 
        };
    }

    const messageInput = document.getElementById("addPeopleInput");
    document.getElementById("chatMessages").addEventListener("scroll", handleScroll);
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


    infoButton.onclick = () => {
        infoBox.style.display = "block";
    };

    closeInfoBox.onclick = () => {
        infoBox.style.display = "none";
        addPeopleSection.style.display = "none"; 
    
        document.getElementById("userSearchResults").innerHTML = "";
    };

    addPeopleButton.onclick = () => {
        addPeopleSection.style.display = "block";
    };

    confirmAddPeopleBtn.onclick = () => {
        const emailOrUsername = document.getElementById("addPeopleInput").value;
        if (emailOrUsername) {
            console.log(`Adding: ${emailOrUsername}`); 
            document.getElementById("addPeopleInput").value = "";
        }
    };
}

function handleScroll() {
    const chatMessages = document.getElementById("chatMessages");
    if (chatMessages.scrollTop === 0 && !loadingMessages) {
        loadMoreMessages();
    }
}


async function loadMoreMessages() {
    loadingMessages = true;

    const lastMessageId = globalMessages.length > 0 ? globalMessages[0].message_id.Int64 : null;

    fetch(`/loadMoreMessages?chatId=${chatId}&lastMessageId=${lastMessageId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(newMessages => {

            if (newMessages.messages && newMessages.messages.length > 0) {
                globalMessages = [...newMessages.messages, ...globalMessages];
                render(chatId);
            }
        })
        .catch(error => {
            console.error("Error loading more messages:", error);
        })
        .finally(() => {
            loadingMessages = false;
        });
}


document.addEventListener("DOMContentLoaded", () => {
    const pathSegments = window.location.pathname.split('/');
    const chatIdParam = pathSegments[pathSegments.length - 1]; 
    if (chatIdParam) {
       initialize(chatIdParam); 
    } else {
        console.error("chatId not found in the URL.");
    }
});
