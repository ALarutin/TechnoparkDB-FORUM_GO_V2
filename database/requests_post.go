package database

import (
	"data_base/models"
)

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
	if related[0] == "" {
		str :=
			`SELECT
			p.id, p.author, p.thread, p.forum, p.message, p.is_edited, p.parent, p.created
			FROM public.post as p
			WHERE p.id = $1`
		row := db.dataBase.QueryRow(str, id)
		var post models.Post
		err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum,
			&post.Message, &post.IsEdited, &post.Parent, &post.Created)
		postInfo.Post = post
	} else if related[0] == user {
		if len(related) == 1 {
			str :=
				`SELECT
				p.id, p.author, p.thread, p.forum, p.message, p.is_edited, p.parent, p.created,
				a.id, a.nickname, a.email, a.fullname, a.about
				FROM public.post as p
				INNER JOIN public.person a ON p.author = a.nickname
				WHERE p.id = $1`
			row := db.dataBase.QueryRow(str, id)
			var post models.Post
			var user models.User
			err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum, &post.Message, &post.IsEdited, &post.Parent, &post.Created,
				&user.ID, &user.Nickname, &user.Email, &user.Fullname, &user.About)
			postInfo.Post = post
			postInfo.Person = &user
		} else if len(related) == 2 {
			if related[1] == forum {
				str :=
					`SELECT
				p.id, p.author, p.thread, p.forum, p.message, p.is_edited, p.parent, p.created,
				a.id, a.nickname, a.email, a.fullname, a.about,
				f.id, f.slug, f.author, f.title, f.posts, f.threads
				FROM public.post as p
				INNER JOIN public.person a ON p.author = a.nickname
				INNER JOIN forum f ON p.forum = f.slug
				WHERE p.id = $1`
				row := db.dataBase.QueryRow(str, id)
				var post models.Post
				var user models.User
				var forum models.Forum
				err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum, &post.Message, &post.IsEdited, &post.Parent, &post.Created,
					&user.ID, &user.Nickname, &user.Email, &user.Fullname, &user.About,
					&forum.ID, &forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
				postInfo.Post = post
				postInfo.Person = &user
				postInfo.Forum = &forum
			} else {
				str :=
					`SELECT
				p.id, p.author, p.thread, p.forum, p.message, p.is_edited, p.parent, p.created,
				a.id, a.nickname, a.email, a.fullname, a.about,
				t.id, (CASE WHEN t.slug ISNULL THEN '' ELSE t.slug END), t.author, t.forum, t.title, t.message, t.votes, t.created
				FROM public.post as p
				INNER JOIN public.person a ON p.author = a.nickname
				INNER JOIN thread t on p.thread = t.id
				WHERE p.id = $1`
				row := db.dataBase.QueryRow(str, id)
				var post models.Post
				var user models.User
				var thread models.Thread
				err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum, &post.Message, &post.IsEdited, &post.Parent, &post.Created,
					&user.ID, &user.Nickname, &user.Email, &user.Fullname, &user.About,
					&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title, &thread.Message, &thread.Votes, &thread.Created)
				postInfo.Post = post
				postInfo.Person = &user
				postInfo.Thread = &thread
			}
		} else {
			str :=
				`SELECT
				p.id, p.author, p.thread, p.forum, p.message, p.is_edited, p.parent, p.created , 
				a.id, a.nickname, a.email, a.fullname, a.about,
				t.id, (CASE WHEN t.slug ISNULL THEN '' ELSE t.slug END), t.author, t.forum, t.title, t.message, t.votes, t.created,
				f.id, f.slug, f.author, f.title, f.posts, f.threads
				FROM public.post as p
				INNER JOIN public.person a ON p.author = a.nickname
				INNER JOIN thread t on p.thread = t.id
				INNER JOIN forum f ON p.forum = f.slug
				WHERE p.id = $1`
			row := db.dataBase.QueryRow(str, id)
			var post models.Post
			var user models.User
			var thread models.Thread
			var forum models.Forum
			err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum, &post.Message, &post.IsEdited, &post.Parent, &post.Created,
				&user.ID, &user.Nickname, &user.Email, &user.Fullname, &user.About,
				&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title, &thread.Message, &thread.Votes, &thread.Created,
				&forum.ID, &forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
			postInfo.Post = post
			postInfo.Person = &user
			postInfo.Thread = &thread
			postInfo.Forum = &forum
		}
	} else if related[0] == thread {
		if len(related) == 1 {
			str :=
				`SELECT
				p.id, p.author, p.thread, p.forum, p.message, p.is_edited, p.parent, p.created,
				t.id, (CASE WHEN t.slug ISNULL THEN '' ELSE t.slug END), t.author, t.forum, t.title, t.message, t.votes, t.created
				FROM public.post as p
				INNER JOIN thread t ON p.thread = t.id
				WHERE p.id = $1`
			row := db.dataBase.QueryRow(str, id)
			var post models.Post
			var thread models.Thread
			err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum, &post.Message, &post.IsEdited, &post.Parent, &post.Created,
				&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title, &thread.Message, &thread.Votes, &thread.Created)
			postInfo.Post = post
			postInfo.Thread = &thread
		} else {
			str :=
				`SELECT
				p.id, p.author, p.thread, p.forum, p.message, p.is_edited, p.parent, p.created,
				t.id, (CASE WHEN t.slug ISNULL THEN '' ELSE t.slug END), t.author, t.forum, t.title, t.message, t.votes, t.created,
				f.id, f.slug, f.author, f.title, f.posts, f.threads
				FROM public.post as p
				INNER JOIN thread t ON p.thread = t.id
				INNER JOIN forum f ON p.forum = f.slug
				WHERE p.id = $1`
			row := db.dataBase.QueryRow(str, id)
			var post models.Post
			var thread models.Thread
			var forum models.Forum
			err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum, &post.Message, &post.IsEdited, &post.Parent, &post.Created,
				&thread.ID, &thread.Slug, &thread.Author, &thread.Forum, &thread.Title, &thread.Message, &thread.Votes, &thread.Created,
				&forum.ID, &forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
			postInfo.Post = post
			postInfo.Thread = &thread
			postInfo.Forum = &forum
		}
	} else {
		str := `SELECT
				p.id, p.author, p.thread, p.forum, p.message, p.is_edited, p.parent, p.created,
				f.id, f.slug, f.author, f.title, f.posts, f.threads
				FROM public.post as p
				INNER JOIN forum f ON p.forum = f.slug
				WHERE p.id = $1`
		row := db.dataBase.QueryRow(str, id)
		var post models.Post
		var forum models.Forum
		err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum, &post.Message, &post.IsEdited, &post.Parent, &post.Created,
			&forum.ID, &forum.Slug, &forum.User, &forum.Title, &forum.Posts, &forum.Threads)
		postInfo.Post = post
		postInfo.Forum = &forum
	}
	return
}

func (db *databaseManager) GetPost(id int) (post models.Post, err error) {
	row := db.dataBase.QueryRow(
		`SELECT id, author, thread, forum, message, is_edited, parent, created 
				FROM func_get_post($1::INT)`, id)

	err = row.Scan(&post.ID, &post.Author, &post.Thread, &post.Forum,
		&post.Message, &post.IsEdited, &post.Parent, &post.Created)
	if err != nil {
		return
	}
	return
}
