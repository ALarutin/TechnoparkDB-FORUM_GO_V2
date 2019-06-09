package database

import (
	"data_base/models"
	"data_base/presentation/logger"
	"fmt"
	"github.com/jackc/pgx"
)

func (db *databaseManager) CreatePost(posts []models.Post, id int, forum string) (outPs []models.Post, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	str := `INSERT INTO public.post(author, thread, forum, message, parent, created) VALUES`

	for i, p := range posts{
		if i != 0 {
			str += `,`
		}
		str += fmt.Sprintf(` ('%v', %v, '%v', '%v', %v, now())`, p.Author, id, forum,  p.Message, p.Parent)
	}

	str += ` RETURNING id, author, thread, forum, message, is_edited, parent, created`

	logger.Error.Print(str)
	rows, err := tx.Query(str)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.ID, &p.Author, &p.Thread, &p.Forum,
			&p.Message, &p.IsEdited, &p.Parent, &p.Created)
		if err != nil {
			return
		}
		outPs = append(outPs, p)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) GetThreadById(threadId int) (thread models.Thread, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_get_thread_by_id($1::INT)`,threadId)
	err = row.Scan(&thread.IsNew, &thread.ID, &thread.Slug, &thread.Author, &thread.Forum,
		&thread.Title, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) GetThreadBySlug(slug string) (thread models.Thread, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_get_thread_by_slug($1::citext)`, slug)
	err = row.Scan(&thread.IsNew, &thread.ID, &thread.Slug, &thread.Author, &thread.Forum,
		&thread.Title, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) UpdateThread(message string, title string, slug string, threadId int) (thread models.Thread, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_update_thread($1::text, $2::text, $3::citext, $4::INT)`,
		message, title, slug, threadId)
	err = row.Scan(&thread.IsNew, &thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title,
		&thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) CreateOrUpdateVote(vote models.Vote, slug string, threadId int) (thread models.Thread, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_create_or_update_vote($1::citext, $2::citext, $3::INT, $4::INT)`,
		vote.Nickname, slug, threadId, vote.Voice)
	err = row.Scan(&thread.IsNew, &thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title,
		&thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}


func (db *databaseManager) GetPosts(slug string, id int, limit int, since int, sort string, desc bool) (posts []models.Post, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := db.getRowsForGetPosts(tx, id, limit, since, slug, sort, desc)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum,
			&post.Message, &post.IsEdited, &post.Parent, &post.Created)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return
	}

	err = tx.Commit()
	return
}

func(db *databaseManager) getRowsForGetPosts(tx *pgx.Tx, id, limit, since int, slug, sort string, desc bool) (rows *pgx.Rows, err error){


	switch sort {
	case "flat":
		rows, err = tx.Query(
			`SELECT id, author, thread, forum, message, is_edited, parent, created 
					FROM func_get_posts_flat($1::citext, $2::INT, $3::INT, $4::INT, $5::BOOLEAN)`,
			slug, id, limit, since, desc)
		if err != nil {
			return
		}
	case "tree":
		rows, err = tx.Query(
			`SELECT id, author, thread, forum, message, is_edited, parent, created
					FROM func_get_posts_tree($1::citext, $2::INT, $3::INT, $4::INT, $5::BOOLEAN)`,
			slug, id, limit, since, desc)
		if err != nil {
			return
		}
	case "parent_tree":
		rows, err = tx.Query(
			`SELECT id, author, thread, forum, message, is_edited, parent, created
					FROM func_get_posts_parent_tree($1::citext, $2::INT, $3::INT, $4::INT, $5::BOOLEAN)`,
			slug, id, limit, since, desc)
		if err != nil {
			return
		}
	default:
		rows, err = tx.Query(
			`SELECT id, author, thread, forum, message, is_edited, parent, created
					FROM func_get_posts($1::citext, $2::INT, $3::INT, $4::INT, $5::BOOLEAN)`,
			slug, id, limit, since, desc)
		if err != nil {
			return
		}
	}
	return
}