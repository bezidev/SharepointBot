package db

type SharepointNotification struct {
	ID             string `db:"id"`
	Name           string `db:"name"`
	Description    string `db:"description"`
	CreatedOn      int    `db:"created_on"`
	ModifiedOn     int    `db:"modified_on"`
	CreatedBy      string `db:"created_by"`
	ModifiedBy     string `db:"modified_by"`
	MessageIDs     string `db:"message_ids"`
	ExpiresOn      int    `db:"expires_on"`
	HasAttachments bool   `db:"has_attachments"`
}

func (db *sqlImpl) GetSharepointNotification(id string) (notification SharepointNotification, err error) {
	err = db.db.Get(&notification, "SELECT * FROM sharepoint_notifications WHERE id=$1", id)
	return notification, err
}

func (db *sqlImpl) GetSharepointNotifications() (notification []SharepointNotification, err error) {
	err = db.db.Select(&notification, "SELECT * FROM sharepoint_notifications ORDER BY modified_on ASC")
	return notification, err
}

func (db *sqlImpl) InsertSharepointNotification(notification SharepointNotification) (err error) {
	_, err = db.db.NamedExec(
		`INSERT INTO sharepoint_notifications
	(id,
	 name,
	 description,
	 created_on,
	 modified_on,
	 created_by,
	 modified_by,
	 message_ids,
	 expires_on,
	 has_attachments) 
VALUES (:id,
		:name,
		:description,
		:created_on,
		:modified_on,
		:created_by,
		:modified_by,
		:message_ids,
		:expires_on,
		:has_attachments)
`, notification)
	return err
}

func (db *sqlImpl) UpdateSharepointNotification(notification SharepointNotification) error {
	_, err := db.db.NamedExec(
		`UPDATE sharepoint_notifications SET
			name=:name,
			description=:description,
			modified_on=:modified_on,
			modified_by=:modified_by,
			message_ids=:message_ids,
			expires_on=:expires_on,
			has_attachments=:has_attachments
WHERE id=:id`,
		notification)
	return err
}

func (db *sqlImpl) DeleteSharepointNotification(id string) error {
	_, err := db.db.Exec(`DELETE FROM sharepoint_notifications WHERE id=$1`, id)
	return err
}
