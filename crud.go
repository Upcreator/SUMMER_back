package main

import (
	"github.com/google/uuid"
)

// Create a new news item
func createNews(news News) error {
	news.ID = uuid.New()
	_, err := db.Exec(`
        INSERT INTO news (id, name, description, type, created_at)
        VALUES ($1, $2, $3, $4, $5)`,
		news.ID, news.Name, news.Description, news.Type, news.CreatedAt)
	return err
}

// Get all news items
func getNews() ([]News, error) {
	rows, err := db.Query("SELECT id, name, description, type, created_at FROM news")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsItems []News
	for rows.Next() {
		var n News
		if err := rows.Scan(&n.ID, &n.Name, &n.Description, &n.Type, &n.CreatedAt); err != nil {
			return nil, err
		}
		newsItems = append(newsItems, n)
	}
	return newsItems, nil
}

// Get a specific news item by id
func getNewsByID(id uuid.UUID) (News, error) {
	var n News
	err := db.QueryRow("SELECT id, name, description, type, created_at FROM news WHERE id = $1", id).
		Scan(&n.ID, &n.Name, &n.Description, &n.Type, &n.CreatedAt)
	return n, err
}

// Update a specific news item
func updateNews(id uuid.UUID, name *string, description *string, newsType *string) error {
	_, err := db.Exec(`
		UPDATE news 
		SET 
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			type = COALESCE($3, type)
		WHERE id = $4`,
		name, description, newsType, id)
	return err
}

// Delete a specific news item
func deleteNews(id uuid.UUID) error {
	_, err := db.Exec("DELETE FROM news WHERE id = $1", id)
	return err
}
