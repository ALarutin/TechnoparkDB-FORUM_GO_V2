package database

import (
	"data_base/models"
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

	for i, p := range posts {
		if i != 0 {
			str += `,`
		}
		str += fmt.Sprintf(` ('%v', %v, '%v', '%v', %v, now())`, p.Author, id, forum, p.Message, p.Parent)
	}

	str += ` RETURNING id, author, thread, forum, message, is_edited, parent, created`

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
	row := db.dataBase.QueryRow(
		`SELECT id, (CASE WHEN slug ISNULL THEN '' ELSE slug END), author, forum, title, message, votes, created
			FROM public.thread WHERE id = $1`, threadId)
	err = row.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum,
		&thread.Title, &thread.Message, &thread.Votes, &thread.Created)
	return
}

func (db *databaseManager) GetThreadBySlug(slug string) (thread models.Thread, err error) {
	row := db.dataBase.QueryRow(
		`SELECT id, (CASE WHEN slug ISNULL THEN '' ELSE slug END), author, forum, title, message, votes, created
			FROM public.thread WHERE slug = $1`, slug)
	err = row.Scan(&thread.ID, &thread.Slug, &thread.Author, &thread.Forum,
		&thread.Title, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return
	}
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

func (db *databaseManager) GetThreadIdBySlug(slug string, i int) (id int, err error) {
	row := db.dataBase.QueryRow(`SELECT id FROM public.thread WHERE slug = $1 or id = $2;`, slug, i)
	err = row.Scan(&id)
	if err != nil {
		return
	}
	return
}

func (db *databaseManager) GetPosts(id int, limit, since, sort, desc string) (posts []models.Post, err error) {
	rows, err := db.getRowsForGetPosts(id, limit, since, sort, desc)
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
	return
}

func (db *databaseManager) getRowsForGetPosts(id int, limit, since, sort, desc string) (rows *pgx.Rows, err error) {
	str := `SELECT id, author, thread, forum, message, is_edited, parent, created
			FROM public.post WHERE thread = $1 `

	switch sort {
	case "flat":
		if since != "" {
			if desc == "true" {
				str += `AND id < ` + since + ` `
			} else {
				str += `AND id > ` + since + ` `
			}
		}
		str += `ORDER BY id `
		if desc == "true" {
			str += `DESC `
		}
		str += `Limit ` + limit
	case "tree":
		if since != "" {
			if desc == "true" {
				str += `AND post_path < (SELECT post_path FROM public.post WHERE id = ` + since + `) `
			} else {
				str += `AND post_path > (SELECT post_path FROM public.post WHERE id = ` + since + `) `
			}
		}
		str += `ORDER BY post_path `
		if desc == "true" {
			str += `DESC `
		}
		str += `Limit ` + limit
	case "parent_tree":
		str += `AND rootid IN (SELECT id FROM public.post
				WHERE thread = $1 AND parent = 0 `
		if since != "" {
			if desc == "true" {
				str += `AND id < (SELECT rootid FROM public.post WHERE id = ` + since + `) `
			} else {
				str += `AND id > (SELECT rootid FROM public.post WHERE id = ` + since + `) `
			}
		}
		str += `ORDER BY id `
		if desc == "true" {
			str += `DESC `
		}
		str += `Limit ` + limit + `) `
		str += `ORDER BY rootid `
		if desc == "true" {
			str += `DESC`
		}
		str += `, post_path`
	default:
		if since != "" {
			if desc == "true" {
				str += `AND id < ` + since + ` `
			} else {
				str += `AND id > ` + since + ` `
			}
		}
		str += `ORDER BY created `
		if desc == "true" {
			str += `DESC `
		}
		str += `,id `
		if desc == "true" {
			str += `DESC `
		}
		str += `Limit ` + limit
	}

	rows, err = db.dataBase.Query(str, id)
	if err != nil {
		return
	}

	//switch sort {
	//case "flat":
	//	rows, err = db.dataBase.Query(
	//		`SELECT id, author, thread, forum, message, is_edited, parent, created
	//				FROM func_get_posts_flat($1::INT, $2::INT, $3::INT, $4::BOOLEAN)`,
	//		id, limit, since, desc)
	//	if err != nil {
	//		return
	//	}
	//case "tree":
	//	rows, err = db.dataBase.Query(
	//		`SELECT id, author, thread, forum, message, is_edited, parent, created
	//				FROM func_get_posts_tree($1::INT, $2::INT, $3::INT, $4::BOOLEAN)`,
	//		id, limit, since, desc)
	//	if err != nil {
	//		return
	//	}
	//case "parent_tree":
	//	rows, err = db.dataBase.Query(
	//		`SELECT id, author, thread, forum, message, is_edited, parent, created
	//				FROM func_get_posts_parent_tree($1::INT, $2::INT, $3::INT, $4::BOOLEAN)`,
	//		id, limit, since, desc)
	//	if err != nil {
	//		return
	//	}
	//default:
	//	rows, err = db.dataBase.Query(
	//		`SELECT id, author, thread, forum, message, is_edited, parent, created
	//				FROM func_get_posts($1::INT, $2::INT, $3::INT, $4::BOOLEAN)`,
	//		id, limit, since, desc)
	//	if err != nil {
	//		return
	//	}
	//}
	return
}
