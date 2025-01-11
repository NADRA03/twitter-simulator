// Render function to initialize the basic HTML structure for the search page
export function render() {
    return `
        <div class="search-container">
            <input type="text" id="searchInput" placeholder="Search" />
            
            <!-- Search Type Buttons -->
            <div class="search-type-buttons">
                <button id="peopleTab" class="active">People</button>
                <button id="latestTab">Latest</button>
                <button id="topTab">Top</button>
                <button id="mediaTab">Media</button>
            </div>

            <!-- Search Results Container -->
            <div id="userSearchResults"></div>
        </div>
    `;
}

// Function to search for users based on input
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

// Function to display search results for "People"
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
            <img src="${user.image_url}" style=" border-radius: 20px; cursor: pointer;" data-user-id="${user.id}" onError="this.onerror=null; this.src='/assets/user2.png'" /> 
            <span style="cursor: pointer;" data-user-id="${user.id}">${user.username}</span>
        `;

        // Add event listeners for both image and name
        const userImage = userElement.querySelector('img');
        const userName = userElement.querySelector('span');

        // Redirect to the user's profile page when image or name is clicked
        [userImage, userName].forEach(element => {
            element.addEventListener('click', () => {
                window.location.href = `/a_profile/${user.id}`;
            });
        });

        userSearchResults.appendChild(userElement);
    });
}

// Initialize function to set up event listeners
export async function initialize() {
    const searchInput = document.getElementById("searchInput");
    searchInput.addEventListener("input", () => {
        const searchTerm = searchInput.value;
        searchUsers(searchTerm);
    });

    // Tab buttons
    const peopleTab = document.getElementById("peopleTab");
    const latestTab = document.getElementById("latestTab");
    const topTab = document.getElementById("topTab");
    const mediaTab = document.getElementById("mediaTab");

    // Set "People" as the active tab initially
    function setActiveTab(activeTab) {
        [peopleTab, latestTab, topTab, mediaTab].forEach(tab => {
            tab.classList.remove("active");
        });
        activeTab.classList.add("active");

        // For now, only display "People" results
        if (activeTab === peopleTab) {
            searchUsers(searchInput.value);
        } else {
            // Clear results for other tabs
            document.getElementById("userSearchResults").innerHTML = `<p>No results for ${activeTab.innerText} yet.</p>`;
        }
    }

    // Event listeners for tabs
    peopleTab.addEventListener("click", () => setActiveTab(peopleTab));
    latestTab.addEventListener("click", () => setActiveTab(latestTab));
    topTab.addEventListener("click", () => setActiveTab(topTab));
    mediaTab.addEventListener("click", () => setActiveTab(mediaTab));
}
