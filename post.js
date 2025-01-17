export function render(postId) {
    return `
        <div id="post-container">
            <div id="post"></div>
            <div id="reply-form-container"></div> <!-- Add reply form container -->
            <div id="comments-container">
                <div id="comments"></div>
            </div>
        </div>
    `;
}

export async function initialize(postId) {
    // Fetch post data
    console.log(postId);
    const postResponse = await fetch(`/a_post/?id=${postId}`);
    if (!postResponse.ok) {
        document.getElementById('post').innerHTML = `<p>Error loading post.</p>`;
        return;
    }
    const post = await postResponse.json();
    const defaultImageUrl = '/assets/user2.png';
    // Render post
    const postContainer = document.getElementById('post');
    postContainer.innerHTML = `
        <div class="post">
            <img src="${post.image_url}" onerror="this.onerror=null; this.src='${defaultImageUrl}'" class="avatar">
            <div class="post-content">
                <h4>${post.Username}</h4>
                <p>${post.content}</p>
            </div>
        </div>
    `;
    const postImage = postContainer.querySelector(".avatar");
    postImage.addEventListener('click', (e) => {
        e.stopPropagation(); 
        window.location.href = `/a_profile/${post["user-id"]}`;
    });

    const userResponse = await fetch('/User');
        if (!userResponse.ok) {
            throw new Error(`Failed to fetch user details: ${userResponse.statusText}`);
        }
        const user = await userResponse.json();

    // Add reply form for the post
    const replyFormContainer = document.getElementById('reply-form-container');
    replyFormContainer.innerHTML = `
    <form id="reply-form" style="display: flex; align-items: center; gap: 10px;">
        <img src="${user.image_url}" onerror="this.onerror=null; this.src='${defaultImageUrl}'" style="width: 50px; height: 50px; border-radius: 50%;">
        <input id="reply-content" type="text" placeholder="Post your reply" required style="flex-grow: 1; padding: 5px;">
        <button class="post-button" type="submit">Reply</button>
    </form>
    `;

    document.getElementById('reply-form').addEventListener('submit', (event) => {
        event.preventDefault();
        const content = document.getElementById('reply-content').value.trim();
        if (content) {
            submitComment(postId, content);
        }
    });

    // Fetch comments
    const commentsResponse = await fetch(`/postComments/?post_id=${postId}`);
    if (!commentsResponse.ok) {
        document.getElementById('comments').innerHTML = `<p>No comments yet.</p>`;
        return;
    }
    const comments = await commentsResponse.json();

    // Render comments
    const commentsContainer = document.getElementById('comments');
    commentsContainer.innerHTML = comments.length
        ? comments.map(renderComment).join('')
        : `<p>Be the first to comment!</p>`;
}

// Render a single comment
function renderComment(comment) {
    const avatarURL = comment.image_url || '/assets/user2.png';
    const defaultImageUrl = '/assets/user2.png';
    const commentHTML = `
        <div class="comment">
            <img            
                src="${avatarURL}" 
                class="avatar" 
                style="cursor: pointer;" 
                onerror="this.onerror=null; this.src='${defaultImageUrl}'" 
                data-user-id="${comment.user_id}">
            <div class="comment-content">
                <h5>${comment.username}</h5>
                <p>${comment.content}</p>
            </div>
        </div>
    `;

    // Add event listener to the avatar image
    setTimeout(() => {
        const avatarImage = document.querySelector(`img[data-user-id="${comment.user_id}"]`);
        if (avatarImage) {
            avatarImage.addEventListener('click', () => {
                window.location.href = `/a_profile/${comment.user_id}`;
            });
        }
    }, 0);

    return commentHTML;
}

// Submit a new comment
async function submitComment(postId, content) {
    try {
        const userResponse = await fetch('/User');
        if (!userResponse.ok) {
            throw new Error(`Failed to fetch user details: ${userResponse.statusText}`);
        }
        const user = await userResponse.json();

        const commentData = {
            user_id: user.id,
            post_id: parseInt(postId, 10),
            content: content,
            username: user.username,
        };

        console.log('Comment Data:', commentData);

        const response = await fetch('/createComment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(commentData),
        });

        if (response.ok) {
            alert('Comment submitted successfully!');
            location.reload();
        } else {
            const errorText = await response.text();
            console.error('Server response:', errorText);
            alert('Failed to submit comment.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to submit comment.');
    }
}

