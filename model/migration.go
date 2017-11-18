package model

import (
	"fmt"
)

func migration0() {
	CheckExec(`CREATE TABLE config(name VARCHAR(64), val TEXT);`)
	CheckExec(`CREATE UNIQUE INDEX config_key_index on config(name);`)

	CheckExec(`CREATE TABLE board(
				id					INTEGER PRIMARY KEY AUTOINCREMENT,
				code				TEXT NOT NULL,
				name				TEXT NOT NULL
			);`)
	// We reserve the 0-index as a possible nil value when in go
	CheckExec(`UPDATE SQLITE_SEQUENCE SET seq = 1 WHERE name = 'board';`)

	CheckExec(`CREATE TABLE post(
				id					INTEGER PRIMARY KEY AUTOINCREMENT,
				content				TEXT NOT NULL,
				thread_parent_id	INTEGER REFERENCES post(id),
				board_parent_id		INTEGER NOT NULL REFERENCES board(id),
				posted_at			DATETIME DEFAULT current_timestamp,
				updated_at			DATETIME DEFAULT current_timestamp,
				archived_at			DATETIME

				CHECK(
					length(content) <= 2000
				)
			);`)
	CheckExec(`CREATE INDEX post_id_index on post(id)`)
	CheckExec(`CREATE INDEX post_thread_parent_id_index on post(thread_parent_id)`)
	CheckExec(`UPDATE SQLITE_SEQUENCE SET seq = 1 WHERE name = 'post';`)

	CheckExec(`INSERT INTO board(code, name) VALUES
				('g', 'programming'),
				('m', 'math'),
				('z', 'misc')
			;`)
}

func testdata() {
	CheckExec(`INSERT INTO post(content, thread_parent_id, board_parent_id) VALUES
				('content 1', NULL, 1),
				('content 2', 1, 1),
				('content 3', 1, 1),
				('content 4', 1, 1),

				('content 1', NULL, 2),
				('content f2', 5, 2),
				('content s3', 5, 2)
			;`)

	for i := 0; i < 200; i += 1 {
		CheckExec(fmt.Sprintf(`INSERT INTO post(content, thread_parent_id, board_parent_id) VALUES
					('content %d', 1, 1)`, i))
	}
}
