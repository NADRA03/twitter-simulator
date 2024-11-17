export function setupEventListeners() {
    const postForm = document.querySelector('form[action="/createPost"]');

    if (postForm) {
        postForm.addEventListener('submit', function(event) {
            event.preventDefault(); 

            // Validate form data
            const isValid = validateForm(); 
            if (!isValid) {
                return; 
            }

            const formData = new FormData(postForm); 

            console.log('Creating post with data:', Object.fromEntries(formData.entries()));

            // Perform the POST request to submit the form
            fetch('/createPost', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    return response.json(); 
                } else {
                    throw new Error('Network response was not ok.');
                }
            })
            .then(data => {
                console.log('Post created:', data);
                window.location.href = '/';
            })
            .catch(error => {
                console.error('There was a problem with the fetch operation:', error);
            });
        });
    }
}

export function render() {
    return`
       <!-- Form and Sidebar Container -->
                    <div class="form-and-sidebar-container">
                        <!-- Sidebar for categories -->
                        <form action="/createPost" method="POST" onsubmit="return validateForm()">
                            <div class="messtrial">
                                <div class="sidebar">
                                    <label for="category"><h3>Categories</h3></label>
                                    <div class="checkbox-group">
                                        {{range .Categories}}
                                        <div class="checkbox-item">
                                            <input type="checkbox" id="{{.ID}}" name="category" value="{{.ID}}">
                                            <label for="{{.ID}}">{{.Name}}</label>
                                        </div>
                                        {{end}}
                                    </div>
                                </div>
                    
                                <!-- Main content area -->
                                <div class="main-creating">
                                    <h3>Creating post ...</h3>
                                    
                                        <div class="form-group">
                                            <label for="title">Title</label>
                                            <input type="text" id="title" name="title" required><br>
                                            <span id="title-error" class="error-message hidden"></span><br>
                                        </div>
                    
                                        <div class="form-group">
                                            <label for="content">Content</label>
                                            <textarea id="content" name="content" rows="10" required></textarea>
                                            <br>
                                            <span id="content-error" class="error-message hidden"></span><br>
                                        </div>
                    
                                        <button type="submit">
                                            <i class="fas fa-paper-plane"></i> Send
                                        </button>
                                    </form>
                                </div>
                            </div>
                        
                        
                    </div>
                </div>
    `
}

export function initializeCreatePost() {
    setupEventListeners(); 
}