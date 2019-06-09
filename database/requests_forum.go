package database

import (
	"data_base/models"
	"time"
)

func (db *databaseManager) CreateForum(forum models.Forum) (f models.Forum, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_create_forum($1::citext, $2::citext, $3::text)`,
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
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_get_forum($1::citext)`, slug)
	err = row.Scan(&forum.IsNew, &forum.ID, &forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) GetThreads(slug string, since time.Time, desc bool, limit int) (threads []models.Thread, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM func_get_threads($1::citext, $2::TIMESTAMP WITH TIME ZONE,
  		$3::BOOLEAN, $4::INT)`, slug, since, desc, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	var thread models.Thread
	for rows.Next() {
		err = rows.Scan(&thread.IsNew, &thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title,
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
	err = tx.Commit()
	return
}

func (db *databaseManager) GetUsers(slug string, since string, desc bool, limit int) (users []models.User, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM func_get_users($1::citext, $2::citext, $3::BOOLEAN, $4::INT)`,
		slug, since, desc, limit)
	if err != nil {
		return
	}
	defer rows.Close()

	var user models.User
	for rows.Next() {
		err = rows.Scan(&user.IsNew, &user.ID, &user.Nickname, &user.Email, &user.Fullname, &user.About)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		err = rows.Err()
		return
	}

	err = tx.Commit()
	return
}
