package controllers

const (
	messageCantFind      = `message": "cant find `
	cantFindUser         = `user with nickname `
	cantFindThread       = `thread with slug or id `
	cantFindForum        = `forum with slug `
	cantFindParentOrUser = `parent or parent in another thread`
	cantFindPost         = `post with id `
	emailUsed            = ` has already taken by another user`
)

const (
	//errorUniqueViolation = `pq: unique_violation`
	errorUniqueViolation = `ERROR: unique_violation (SQLSTATE 23505)`
	errorPqNoDataFound   = `ERROR: no_data_found (SQLSTATE P0002)` //TODO RENAME
	//errorPqNoDataFound       = `pq: no_data_found` //TODO RENAME
	//errorForeignKeyViolation = `pq: foreign_key_violation`
	errorForeignKeyViolation = `ERROR: foreign_key_violation (SQLSTATE 23503)`
)
