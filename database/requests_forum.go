package database

import (
	"data_base/models"
)

func (db *databaseManager) CreateForum(forum models.Forum) (f models.Forum, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_create_forum
	($1::citext, $2::citext, $3::text)`,
		forum.User, forum.Slug, forum.Title)
	err = row.Scan(&f.IsNew, &f.ID, &f.Slug, &f.User, &f.Title, &f.Posts, &f.Threads)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) CreateThread(thread models.Thread) (t models.Thread, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM  func_create_thread
	($1::citext, $2::TIMESTAMP WITH TIME ZONE, $3::citext, $4::text, $5::citext, $6::text)`,
		thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug, thread.Title)
	err = row.Scan(&t.IsNew, &t.ID, &t.Slug, &t.Author, &t.Forum, &t.Title, &t.Message, &t.Votes, &t.Created)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) GetForum(slug string) (forum models.Forum, err error) {
	row := db.dataBase.QueryRow(
		`SELECT id, slug, author, title, posts, threads 
			FROM public.forum WHERE slug = $1`, slug)
	err = row.Scan(&forum.ID, &forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
	return
}

func (db *databaseManager) GetThreads(slug, since, desc, limit string) (threads []models.Thread, err error) {
	str := `SELECT id, (CASE WHEN slug ISNULL THEN '' ELSE slug END), author, forum, title, message, votes, created FROM public.thread
			WHERE forum = $1 `
	if since != "" {
		if desc == "true"{
			str += `AND created <= '` + since + `' `
		} else {
			str += `AND created >= '` + since + `' `
		}
	}
	str += `ORDER BY created `
	if desc == "true" {
		str += `DESC `
	}
	str += `LIMIT $2`

	rows, err := db.dataBase.Query(str, slug, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	var thread models.Thread
	for rows.Next() {
		err = rows.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title,
			&thread.Message, &thread.Votes, &thread.Created)
		if err != nil {
			return
		}
		threads = append(threads, thread)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return
	}
	return
}

func (db *databaseManager) GetUsers(slug, since, desc, limit string) (users []models.User, err error) {
	str := `SELECT * FROM public.person
			WHERE nickname IN (SELECT user_nickname FROM public.forum_users WHERE forum_slug = $1)`
	if since != "" {
		if desc == "true"{
			str += `AND nickname < '` + since + `' `
		} else {
			str += `AND nickname > '` + since + `' `
		}
	}
	str += `ORDER BY nickname `
	if desc == "true" {
		str += `DESC `
	}
	str += `LIMIT $2`

	rows, err := db.dataBase.Query(str, slug, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		err = rows.Scan( &user.ID, &user.Nickname, &user.Email, &user.Fullname, &user.About)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return
	}
	return
}
