package model

func migration0() {
	CheckExec(`CREATE TABLE config(name VARCHAR(64), val TEXT);`)
	CheckExec(`CREATE UNIQUE INDEX config_key_index on config(name);`)

	CheckExec(`CREATE TABLE board(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				code TEXT NOT NULL,
				name TEXT NOT NULL
			);`)
	// We reserve the 0-index as a possible nil value when in go
	CheckExec(`UPDATE SQLITE_SEQUENCE SET seq = 1 WHERE name = 'board';`)

	CheckExec(`CREATE TABLE post(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT NOT NULL,
				thread_parent_id INTEGER REFERENCES post(id),
				board_parent_id INTEGER NOT NULL REFERENCES board(id),
				posted_at INTEGER NOT NULL

				CHECK(
					length(content) <= 2000
				)
			);`)
	CheckExec(`UPDATE SQLITE_SEQUENCE SET seq = 1 WHERE name = 'post';`)
}

func testdata0() {
	CheckExec(`INSERT INTO board(code, name) VALUES
				('g', 'programming'),
				('m', 'math'),
				('z', 'misc')
			;`)

	CheckExec(`INSERT INTO post(content, thread_parent_id, board_parent_id, posted_at) VALUES
				('content 1', NULL, 1, strftime("%s", CURRENT_TIME)),
				('content 2', 1, 1, strftime("%s", CURRENT_TIME)),
				('content 3', 1, 1, strftime("%s", CURRENT_TIME)),
				('content 4', 1, 1, strftime("%s", CURRENT_TIME)),

				('content 1', NULL, 2, strftime("%s", CURRENT_TIME)),
				('content f2', 5, 2, strftime("%s", CURRENT_TIME)),
				('content s3', 5, 2, strftime("%s", CURRENT_TIME))
			;`)
}
