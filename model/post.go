package model

import (
	"log"
	"time"
)

type Post struct {
	Id             int64
	Image          string
	Content        string
	ThreadParentId int64
	BoardParentId  int64
	PostedAt       time.Time
}

type PostW struct {
	ThreadParentId int64
	Content        string
}

func CreatePost(p PostW) (int64, error) {
	t, err := GetPostById(p.ThreadParentId)
	if err != nil {
		return 0, err
	}

	_, err = db.Exec(`INSERT INTO post(content, thread_parent_id, board_parent_id) VALUES
					 (?, ?, ?);`, p.Content, p.ThreadParentId, t.BoardParentId)
	if err != nil {
		return 0, err
	}

	r := db.QueryRow(`SELECT last_insert_rowid();`)

	var id int64
	err = r.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

type ThreadW struct {
	BoardParentId int64
	Content       string
}

func CreateThread(t ThreadW) (int64, error) {
	_, err := db.Exec(`INSERT INTO post(content, board_parent_id) VALUES
					 (?, ?);`, t.Content, t.BoardParentId)
	if err != nil {
		return 0, err
	}

	r := db.QueryRow(`SELECT last_insert_rowid();`)

	var id int64
	err = r.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetPostById(id int64) (Post, error) {
	r := db.QueryRow(`SELECT * FROM post WHERE id = ?;`, id)

	var p Post
	// Convert possible NULL thread_parent_id fields to 0, these are unused	by the database.
	var tid interface{}

	err := r.Scan(&p.Id, &p.Content, &tid, &p.BoardParentId, &p.PostedAt)
	if err != nil {
		return Post{}, err
	}

	if tid == nil {
		p.ThreadParentId = 0
	} else {
		p.ThreadParentId = tid.(int64)
	}

	return p, nil
}

func (p *Post) GetParentThread() ([]Post, error) {
	var qId int64
	if p.ThreadParentId == 0 {
		qId = p.Id
	} else {
		qId = p.ThreadParentId
	}

	r, err := db.Query(`SELECT * FROM post WHERE thread_parent_id = ? OR id = ?;`, qId, qId)
	if err != nil {
		return []Post{}, err
	}
	defer r.Close()

	var ps []Post
	for r.Next() {
		var p Post
		// Convert possible NULL thread_parent_id fields to 0, these are unused	by the database.
		var tid interface{}

		err := r.Scan(&p.Id, &p.Content, &tid, &p.BoardParentId, &p.PostedAt)
		if err != nil {
			return []Post{}, err
		}

		if tid == nil {
			p.ThreadParentId = 0
		} else {
			p.ThreadParentId = tid.(int64)
		}

		ps = append(ps, p)
	}

	if len(ps) == 0 {
		log.Panicln("a thread must have one post (itself)")
	}

	return ps, nil
}
