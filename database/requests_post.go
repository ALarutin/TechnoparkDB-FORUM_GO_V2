package database

import "data_base/models"

const (
	user   = "user"
	forum  = "forum"
	thread = "thread"
)

func (db *databaseManager) UpdatePost(message string, id int) (post models.Post, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		`SELECT id, author, thread, forum, message, is_edited, parent, created 
				FROM func_update_post($1::text, $2::INT)`,
		message, id)
	err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum,
		&post.Message, &post.IsEdited, &post.Parent, &post.Created)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) GetPostInfo(id int, related []string) (postInfo models.PostInfo, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		`SELECT id, author, thread, forum, message, is_edited, parent, created 
				FROM func_get_post($1::INT)`, id)
	var post models.Post
	err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum,
		&post.Message, &post.IsEdited, &post.Parent, &post.Created)
	if err != nil {
		return
	}
	postInfo.Post = post

	for _, str := range related {
		switch str {
		case user:
			var user models.User
			row := tx.QueryRow(
				`SELECT * FROM func_get_user($1::citext)`, postInfo.Post.Author)
			err = row.Scan(&user.IsNew, &user.ID, &user.Nickname, &user.Email, &user.Fullname, &user.About)
			if err != nil {
				return
			}
			postInfo.Person = &user
		case thread:
			var thread models.Thread
			row := tx.QueryRow(`SELECT * FROM func_get_thread_by_id($1::INT)`, postInfo.Post.Thread)
			err = row.Scan(&thread.IsNew, &thread.ID, &thread.Slug, &thread.Author, &thread.Forum,
				&thread.Title, &thread.Message, &thread.Votes, &thread.Created)
			if err != nil {
				return
			}
			postInfo.Thread = &thread
		case forum:
			var forum models.Forum
			row := tx.QueryRow(`SELECT * FROM func_get_forum($1::citext)`, postInfo.Post.Forum)
			err = row.Scan(&forum.IsNew, &forum.ID, &forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
			if err != nil {
				return
			}
			postInfo.Forum = &forum
		}
	}
	err = tx.Commit()
	return
}
