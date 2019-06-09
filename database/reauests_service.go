package database

import "data_base/models"

func (db *databaseManager) ClearDatabase() (err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM func_clear_database()`)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *databaseManager) GetDatabase() (database models.Database, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_get_database()`)
	err = row.Scan(&database.Forum, &database.Post, &database.Thread, &database.User)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
