package database

import "data_base/models"

func (db *databaseManager) GetUser(nickname string) (user models.User, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		`SELECT * FROM func_get_user($1::citext)`,
		nickname)
	err = row.Scan(&user.IsNew, &user.ID, &user.Nickname, &user.Email, &user.Fullname, &user.About)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) CreateUser(user models.User) (users []models.User, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM func_create_user($1::citext, $2::citext, $3::text, $4::text)`,
		user.Nickname, user.Email, user.Fullname, user.About)
	if err != nil {
		return
	}
	defer rows.Close()

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

func (db *databaseManager) UpdateUser(user models.User) (u models.User, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		`SELECT * FROM func_update_user($1::citext, $2::citext, $3::text, $4::text)`,
		user.Nickname, user.Email, user.Fullname, user.About)
	err = row.Scan(&u.IsNew, &u.ID, &u.Nickname, &u.Email, &u.Fullname, &u.About)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
