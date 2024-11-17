export function setupEventListeners() {
    const commentForm = document.querySelector('.adding-com form'); 

    if (commentForm) {
        commentForm.addEventListener('submit', function(event) {
            event.preventDefault(); 

            const formData = new FormData(commentForm); 
            const postId = formData.get('post_id'); 
            const content = formData.get('content'); 

            console.log(`Adding comment to post ID: ${postId} with content: ${content}`);

            // Perform the POST request to submit the comment
            fetch('/addComment', {
                method: 'POST',
                body: formData,
            })
            .then(response => {
                if (response.ok) {
                    return response.json(); // Assuming the server responds with JSON
                } else {
                    throw new Error('Network response was not ok.');
                }
            })
            .then(data => {
                console.log('Comment added:', data);
            })
            .catch(error => {
                console.error('There was a problem with the fetch operation:', error);
            });
        });
    }
}

export function render() {
    return `
      <div class="post-comments-content">
                        <div class="com-post">
                
                    <div class="com-post-com">
                        <div class="mainPostCardMeta">
                            <div class="mainPostCardUsername"><i class="fas fa-user-circle"></i> {{.Post.Username}}, post ID :{{.Post.ID}}</div>
                            <div class="post-date"> Posted on : {{.Post.CreatedAt.Format "03:04 PM 02/01/2006 "}} </div>
                        </div>
                        <h2 class="mainPostCardTitle">{{.Post.Title}}</h2>
                        <p class="mainPostCardText">{{.Post.Content}}</p>
                        <div class="interact">
                            <h4><i class="fas fa-comment"></i>comment<div class="comment"></h4>   
                            </div>
                      
                        </div>
                    </div>
                   <!-- beginning of populating comments -->
                  
<div class="bluetrial">
    
    {{range $outerIndex, $outerElement := .Comments}}
    <div class="com-post-main">
        <div class="mainPostCardMeta">
            <div class="mainPostCardUsername"><i class="fas fa-user-circle"></i>{{.ComUsername}}</div>
            <div class="post-date"> commented on :{{ .CreatedAt.Format "02/01/2006 03:04 PM" }} </div>
        </div>
        <p class="comment-content-log">{{.Content}}</p>
    </div>
    {{else}}
    <p>No comments yet.</p>
    {{end}}
</div>
                    <!-- Add Comment Form -->
                    <div class="adding-com">
                        <form action="/addComment" method="POST"">
                            <input type="text" name="post_id" id="post_id" value={{.Post.ID}} hidden>
                            <label for="content">Comment:</label>
                            <textarea id="content" name="content" rows="1" placeholder="Leave a comment"> </textarea>
                            <button type="submit">
                                <i class="fas fa-paper-plane"></i> Send
                            </button>
                            
                        </form>
            
                    </div>
        
        </div>
    `
}

export function initialize() {
    setupEventListeners(); 
}