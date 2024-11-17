package twitter


import(
	"database/sql"
	"fmt"
)

type Category struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
}

func GetAllCategories(db *sql.DB) ([]Category, error) {
	query := `
	 SELECT * FROM Categories;
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Errorf("error quering the categories : %v", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning category: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %v", err)
	}

	return categories, nil
}

func GetPostsByCategory(db *sql.DB, categoryID int) ([]Posts, error) {
	query := `
	SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
	FROM Posts p
	JOIN users u ON p.user_id = u.id
	JOIN PostCategories pc ON p.id = pc.post_id
    WHERE pc.category_id = ?
    ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error quering posts by categories : %v",)
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return posts, nil
}

func GetCategoriesOfPosts(db *sql.DB, postsID int) ([]int, error) {
	query := `
	SELECT category_id
	FROM PostCategories
	WHERE post_id = ?
	`
	rows, err := db.Query(query, postsID)
	if err != nil {
		return nil, fmt.Errorf("error querying categories of post: %v", err)
	}
	defer rows.Close()
	var categoryIDs []int
	for rows.Next() {
		var categoryID int
		if err := rows.Scan(&categoryID); err != nil {
			return nil, fmt.Errorf("error scanning category ID: %v", err)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return categoryIDs, nil
}

func GetCategoryNamesByIDs(db *sql.DB, categoryIDs []int) ([]string, error) {
	// SELECT id,name FROM Categories;
	var categoryNames []string
	for _, id := range categoryIDs {
		query := `
		SELECT name FROM Categories WHERE id = ?
	`
		var categoryName string
		err := db.QueryRow(query, id).Scan(&categoryName)
		if err != nil {
			return nil, fmt.Errorf("error querying category name: %v", err)
		}
		categoryNames = append(categoryNames, categoryName)
	}
	return categoryNames, nil
}
