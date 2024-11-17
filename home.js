export function render() {
    return `
        
        <main class="main-content">
        <div class="post-container">
            {{range .Posts}}
            <div class="post-card">
                <div class="post-meta">
                    <div class="username"><i class="fas fa-user-circle"></i> {{.Username}}</div>
                    <div class="post-date">Posted on: {{ .CreatedAt.Format "02/01/2006 03:04 PM" }}</div>
                </div>
                <h4 class="post-title">{{.Title}}</h4>
                <p class="post-text">{{.Content}}</p>
                <div class="post-divider"></div>
                <div class="interactions">
                    <div class="category-dropdown">
                        <select id="category-dropdown" name="category">
                            <option value="">View Category</option>
                            {{range  $catName := .CategoriesNames}}
                            <option value="{{$catName}}">{{$catName}}</option>
                            {{end}}
                        </select>
                    </div>
                    <div class="comment-link">
                        <h4><a href="/post/{{.ID}}"><i class="fas fa-comment"></i> Comments</a></h4>
                    </div>
                </div>
            </div>
            {{else}}
            <div class="no-posts-message">
                <p>No posts available.</p>
            </div>
            {{end}}
        </div>
    </main>
    `;
}