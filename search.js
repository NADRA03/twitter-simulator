
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
            <img src="${user.image_url}" style=" border-radius: 20px; cursor: pointer;" data-user-id="${user.id}" onError="this.onerror=null; this.src='/assets/user2.png'" /> 
            <span style="cursor: pointer;" data-user-id="${user.id}">${user.username}</span>
        `;

        const userImage = userElement.querySelector('img');
        const userName = userElement.querySelector('span');

        [userImage, userName].forEach(element => {
            element.addEventListener('click', () => {
                window.location.href = `/a_profile/${user.id}`;
            });
        });

        userSearchResults.appendChild(userElement);
    });
}

export async function initialize() {
    const searchInput = document.getElementById("searchInput");
    searchInput.addEventListener("input", () => {
        const searchTerm = searchInput.value;
        searchUsers(searchTerm);
    });

    const peopleTab = document.getElementById("peopleTab");
    const latestTab = document.getElementById("latestTab");
    const topTab = document.getElementById("topTab");
    const mediaTab = document.getElementById("mediaTab");

    function setActiveTab(activeTab) {
        [peopleTab, latestTab, topTab, mediaTab].forEach(tab => {
            tab.classList.remove("active");
        });
        activeTab.classList.add("active");

        if (activeTab === peopleTab) {
            searchUsers(searchInput.value);
        } else {
            document.getElementById("userSearchResults").innerHTML = `<p>No results for ${activeTab.innerText} yet.</p>`;
        }
    }

    // Event listeners for tabs
    peopleTab.addEventListener("click", () => setActiveTab(peopleTab));
    latestTab.addEventListener("click", () => setActiveTab(latestTab));
    topTab.addEventListener("click", () => setActiveTab(topTab));
    mediaTab.addEventListener("click", () => setActiveTab(mediaTab));
}
