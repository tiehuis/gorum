package model

import (
	"database/sql"
	"log"
	"time"

	"github.com/tiehuis/gorum/config"
	"github.com/tiehuis/gorum/util"
)

type Post struct {
	Id             int64
	Content        string
	ThreadParentId NullInt
	BoardParentId  int64
	PostedAt       time.Time
	UpdatedAt      time.Time
	ArchivedAt     NullTime
}

type PostW struct {
	ThreadParentId int64
	Content        string
}

func PrepareQueries() {
	// post.CreatePost
	qInsertNewPost = mustPrepare(`INSERT INTO post(content, thread_parent_id, board_parent_id) VALUES (?, ?, ?);`)
	qUpdateOldPost = mustPrepare(`UPDATE post SET updated_at = CURRENT_TIMESTAMP WHERE id = ?;`)

	// post.CreateThread
	qInsertNewThread = mustPrepare(`INSERT INTO post(content, board_parent_id) VALUES (?, ?);`)
	qGetExpiringThreadId = mustPrepare(`SELECT id FROM post LIMIT 1 OFFSET ?;`)
	qUpdatePostTime = mustPrepare(`UPDATE post SET updated_at = CURRENT_TIMESTAMP WHERE id = ?;`)
}

var qInsertNewPost *sql.Stmt
var qUpdateOldPost *sql.Stmt

func CreatePost(p PostW) (int64, error) {
	t, err := GetPostById(p.ThreadParentId)
	if err != nil {
		return 0, err
	}

	// Pre-format to improve read performance
	p.Content = util.FormatPost(p.Content)

	// All database writes require obtaining a mutex since sqlite only allows
	// a single-writer.
	//
	// TODO: Can we expire this automatically and cancel the query with a rollback?
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Stmt(qInsertNewPost).Exec(p.Content, p.ThreadParentId, t.BoardParentId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// TODO: Check the bump limit at this point in time and don't update if exceeded
	_, err = tx.Stmt(qUpdateOldPost).Exec(p.ThreadParentId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

type ThreadW struct {
	BoardParentId int64
	Content       string
}

var qInsertNewThread *sql.Stmt
var qGetExpiringThreadId *sql.Stmt
var qUpdatePostTime *sql.Stmt

func CreateThread(t ThreadW) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// Pre-format to improve read performance
	t.Content = util.FormatPost(t.Content)

	result, err := tx.Stmt(qInsertNewThread).Exec(t.Content, t.BoardParentId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	r := tx.Stmt(qGetExpiringThreadId).QueryRow(config.BoardThreadLimit)

	var id int64
	err = r.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Stmt(qUpdatePostTime).Exec(id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	nid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return nid, nil
}

func GetPostById(id int64) (Post, error) {
	r := db.QueryRow(`SELECT * FROM post WHERE id = ?;`, id)

	p, err := scanPost(r)
	if err != nil {
		return Post{}, err
	}

	return p, nil
}

func (p *Post) GetParentThreadCount() (int64, error) {
	var qId int64
	if p.ThreadParentId.Valid {
		qId = p.ThreadParentId.Int64
	} else {
		qId = p.Id
	}

	r := db.QueryRow(`SELECT COUNT(*) FROM post WHERE thread_parent_id = ? OR id = ?;`, qId, qId)

	var c int64
	err := r.Scan(&c)
	if err != nil {
		return 0, err
	}

	return c, nil
}

func (p *Post) GetParentThread() ([]Post, error) {
	var qId int64
	if p.ThreadParentId.Valid {
		qId = p.ThreadParentId.Int64
	} else {
		qId = p.Id
	}

	key := "GetParentThread-" + string(qId)
	cb, found := memcache.Get(key)
	if found {
		return cb.([]Post), nil
	}

	r, err := db.Query(`SELECT * FROM post WHERE thread_parent_id = ? OR id = ?;`, qId, qId)
	if err != nil {
		return []Post{}, err
	}
	defer r.Close()

	var ps []Post
	for r.Next() {
		p, err := scanPost(r)
		if err != nil {
			return []Post{}, err
		}

		ps = append(ps, p)
	}

	if len(ps) == 0 {
		log.Panicln("a thread must have one post (itself)")
	}

	memcache.Set(key, ps, 10*time.Second)
	return ps, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanPost(r scanner) (Post, error) {
	var p Post

	// Scanning all the post information takes ages?
	err := r.Scan(&p.Id, &p.Content, &p.ThreadParentId, &p.BoardParentId, &p.PostedAt, &p.UpdatedAt, &p.ArchivedAt)
	if err != nil {
		return Post{}, err
	}

	return p, nil

}
