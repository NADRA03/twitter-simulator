import { loadPage } from "./loader.js";



async function fetchAllUsers() {
    try {
        console.log("Fetching users...");
        const response = await fetch('/allusers');
        console.log("Response status:", response.status);

        if (!response.ok) throw new Error("Failed to fetch users");

        const users = await response.json();
        console.log("Users retrieved:", users); // Log retrieved users

        // You may still want to enable or disable the create chat button based on user availability
        document.getElementById("createChatButton").disabled = users.length === 0;
    } catch (error) {
        console.error("Error fetching users:", error);
        alert("Error fetching users: " + error.message);
    }
}

// Fetch user chats and display them
async function fetchUserChats() {
    try {
        const response = await fetch('/chats/user');
        console.log("Response status:", response.status); // Log the response status

        // Check if the response is OK (status 200)
        if (!response.ok) {
            const text = await response.text(); // Read the response as text
            console.error(`Failed to fetch chats: ${response.status} - ${text}`);
            throw new Error(`Failed to fetch chats: ${response.statusText}`);
        }

        const text = await response.text(); // Read the response as text first
        console.log("Response text:", text); // Log the raw response text

        let chats;
        try {
            chats = JSON.parse(text); // Attempt to parse the text as JSON
        } catch (error) {
            console.error("Failed to parse JSON:", error);
            chats = null; // Set chats to null if parsing fails
        }

        // Check if chats is null
        if (chats === null) {
            console.log("No chats found.");
            const chatList = document.getElementById("chatList");
            chatList.innerHTML = "<li>No chats found.</li>"; // Display message for no chats
            return; // Exit the function if chats is null
        }

        // Check if chats is an array
        if (!Array.isArray(chats)) {
            console.error("Chats is not an array:", chats);
            // Display message for unexpected response format
            const chatList = document.getElementById("chatList");
            chatList.innerHTML = "<li>Unexpected response format. No chats available.</li>"; 
            return; // Exit the function if chats is not an array
        }

        const chatList = document.getElementById("chatList");
        // Check if the chatList element exists
        if (!chatList) {
            console.error("Chat list element not found.");
            return; // Exit the function if the chatList is not found
        }

        chatList.innerHTML = ""; // Clear existing chat list

        // Check if there are any chats
        if (chats.length === 0) {
            console.log("No chats found.");
            chatList.innerHTML = "<li>No chats found.</li>"; // Display message for no chats
            return; // Exit the function early if no chats
        }

        // Populate the chat list with chats
        chats.forEach(chat => {
            const chatItem = document.createElement("li");
            chatItem.className = "chat-item";

            chatItem.innerHTML = `
            <button class="chat-button" data-chat-id="${chat.chat_id}">
                <img src="${chat.image}" alt="${chat.name}" class="chat-image" 
                      onerror="this.onerror=null; this.src='data:image/svg+xml;utf8,<svg xmlns=&quot;http://www.w3.org/2000/svg&quot; fill=&quot;none&quot; height=&quot;48&quot; width=&quot;48&quot; stroke=&quot;currentColor&quot; stroke-linecap=&quot;round&quot; stroke-linejoin=&quot;round&quot; stroke-width=&quot;2&quot; viewBox=&quot;0 0 24 24&quot;><path d=&quot;M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2&quot;/><circle cx=&quot;9&quot; cy=&quot;7&quot; r=&quot;4&quot;/><path d=&quot;M23 21v-2a4 4 0 0 0-3-3.87&quot;/><path d=&quot;M16 3.13a4 4 0 0 1 0 7.75&quot;/></svg>'" />
            </button>
            <div class="chat-details">
                <h3>${chat.name}</h3>
                <p>${chat.last_message.message_text ? chat.last_message.message_text : 'No messages yet'}</p>
            </div>
            `;

            chatList.appendChild(chatItem);
        });

        // Add event listeners for chat buttons after they are rendered
        const chatButtons = document.querySelectorAll('.chat-button');
        chatButtons.forEach(button => {
            button.onclick = function() {
                const chatId = this.dataset.chatId; // Retrieve the chat ID from data attribute
                navigateToChat(chatId); // Call the navigate function with chat ID
            };
        });
    } catch (error) {
        console.error("Error fetching chats:", error);
        // Optionally, display a generic message on error
        const chatList = document.getElementById("chatList");
        chatList.innerHTML = "<li>Error fetching chats. Please try again later.</li>";
    }
}

// Handle form submission

// Render the HTML
export function render() {
    return `
        <div class="chats">
            <h2>Your Chats</h2>

            <button id="showChatFormButton">Add</button>

            <div id="createChatFormContainer" style="display: none;">
                <button id="closeChatFormButton" style="float: right; cursor: pointer;">X</button>
                <form id="createChatForm" action="/chats/create" method="POST">
                    <label for="chatName">Chat Name:</label>
                    <input type="text" name="name" id="chatName" required><br>

                    <label for="chatBio">Bio:</label>
                    <input type="text" name="bio" id="chatBio" required><br>
                    
                    <label for="chatImage">Image URL:</label>
                    <input type="text" name="image" id="chatImage"><br>
                    
                    <label for="chatType">Type:</label>
                    <select name="type" id="chatType" required>
                        <option value="">Select Chat Type</option>
                        <option value="group">Group</option>
                        <option value="private">Private</option>
                    </select><br>

                    <button type="submit" id="createChatButton">Create Chat</button>
                </form>
            </div>

            <ul id="chatList"></ul>
        </div>
    `;
}


// Function to initialize the page
export function initialize() {

    // Fetch initial data
    fetchAllUsers();
    fetchUserChats();

    // Add event listener to the button to show/hide the form
    document.getElementById('showChatFormButton').onclick = function() {
        const formContainer = document.getElementById('createChatFormContainer');
        formContainer.style.display = formContainer.style.display === 'none' ? 'block' : 'none';
    };

    // Add event listener for the close button
    document.getElementById('closeChatFormButton').onclick = function() {
        document.getElementById('createChatFormContainer').style.display = 'none';
    };
}

function navigateToChat(chatId) {
    loadPage(`chat/${chatId}`);
}


// Call the initialize function to set everything up
document.addEventListener("DOMContentLoaded", initialize);