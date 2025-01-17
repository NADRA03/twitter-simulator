import { loadPage } from "./loader.js";
const defaultImageUrl = '/assets/user2.png';


async function fetchAllUsers() {
    try {
        console.log("Fetching users...");
        const response = await fetch('/allusers');
        console.log("Response status:", response.status);

        if (!response.ok) throw new Error("Failed to fetch users");

        const users = await response.json();
        console.log("Users retrieved:", users); 

        document.getElementById("createChatButton").disabled = users.length === 0;
    } catch (error) {
        document.body.dataset.status = '401'; 
        loadPage(`401`);
    }
}

// Fetch user chats and display them
async function fetchUserChats() {
    try {
        const response = await fetch('/chats/user');
        console.log("Response status:", response.status); 

        if (!response.ok) {
            const text = await response.text(); 
            console.error(`Failed to fetch chats: ${response.status} - ${text}`);
            throw new Error(`Failed to fetch chats: ${response.statusText}`);
        }

        const text = await response.text(); 
        console.log("Response text:", text); 

        let chats;
        try {
            chats = JSON.parse(text); 
        } catch (error) {
            console.error("Failed to parse JSON:", error);
            chats = null;
        }

        if (chats === null) {
            console.log("No chats found.");
            const chatList = document.getElementById("chatList");
            chatList.innerHTML = "<li>No chats found.</li>"; 
            return; 
        }

        if (!Array.isArray(chats)) {
            console.error("Chats is not an array:", chats);
            const chatList = document.getElementById("chatList");
            chatList.innerHTML = "<li>Unexpected response format. No chats available.</li>"; 
            return; 
        }

        const chatList = document.getElementById("chatList");
        if (!chatList) {
            console.error("Chat list element not found.");
            return; 
        }

        chatList.innerHTML = ""; 

        if (chats.length === 0) {
            console.log("No chats found.");
            chatList.innerHTML = "<li>No chats found.</li>"; 
            return; 
        }

        chats.forEach(chat => {
            const chatItem = document.createElement("li");
            chatItem.className = "chat-item";

            chatItem.innerHTML = `
            <button class="chat-button" data-chat-id="${chat.chat_id}">
                <img src="${chat.image}" alt="chat group" class="chat-image" onError="this.onerror=null; this.src='${defaultImageUrl}'"
                />
            </button>
            <div class="chat-details">
                <h3>${chat.name}</h3>
                <p style="color:green;">${chat.last_message.message_text ? chat.last_message.message_text : 'No messages yet'}</p>
            </div>
            `;

            chatList.appendChild(chatItem);
        });

        const chatButtons = document.querySelectorAll('.chat-button');
        chatButtons.forEach(button => {
            button.onclick = function() {
                const chatId = this.dataset.chatId; 
                navigateToChat(chatId); 
            };
        });
    } catch (error) {
        console.error("Error fetching chats:", error);
        const chatList = document.getElementById("chatList");
        chatList.innerHTML = "<li>Error fetching chats. Please try again later.</li>";
    }
}


export function render() {
    return `
        <div class="chats">

            <button id="showChatFormButton">Add</button>

            <div id="createChatFormContainer" style="display: none;">
                <button id="closeChatFormButton" style="float: right; cursor: pointer; font-weight: 10;">X</button>
                <form id="createChatForm" action="/chats/create" method="POST">
                    <label for="chatName">Chat Name:</label>
                    <input type="text" name="name" id="chatName" required><br>

                    <label for="chatBio">Bio:</label>
                    <input type="text" name="bio" id="chatBio" required><br>
                    
                    <label for="chatImage">Image URL:</label>
                    <input type="text" name="image" id="chatImage"><br>
                    
                    <select name="type" id="chatType" required>
                    <option value="group" selected>Group</option>
                    </select>
                    <button type="submit" id="createChatButton">Create Chat</button>
                </form>
            </div>

            <ul id="chatList"></ul>
        </div>
    `;
}


export function initialize() {

    fetchAllUsers();
    fetchUserChats();

    document.getElementById('showChatFormButton').onclick = function() {
        const formContainer = document.getElementById('createChatFormContainer');
        formContainer.style.display = formContainer.style.display === 'none' ? 'block' : 'none';
    };

    document.getElementById('closeChatFormButton').onclick = function() {
        document.getElementById('createChatFormContainer').style.display = 'none';
    };
}

function navigateToChat(chatId) {
    loadPage(`chat/${chatId}`);
}


document.addEventListener("DOMContentLoaded", initialize);