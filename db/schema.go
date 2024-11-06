package db

const schema string = `
CREATE TABLE IF NOT EXISTS sharepoint_notifications (
	id						VARCHAR(60)    PRIMARY KEY,
	name					VARCHAR,
	description				VARCHAR,
	created_on				INTEGER,
	modified_on				INTEGER,
	message_ids				JSON,
	created_by				VARCHAR(100),
	modified_by				VARCHAR(100),
	expires_on				INTEGER,
	has_attachments			BOOLEAN
);
`
