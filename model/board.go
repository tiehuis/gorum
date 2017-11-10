package model

import (
	"time"
)

type Board struct {
	Id   int64
	Code string
	Name string
}

func GetBoardById(id int64) (Board, error) {
	key := "GetBoardById-" + string(id)
	cb, found := memcache.Get(key)
	if found {
		return cb.(Board), nil
	}

	r := db.QueryRow(`SELECT * FROM board WHERE id = ?;`, id)

	var b Board
	err := r.Scan(&b.Id, &b.Code, &b.Name)
	if err != nil {
		return Board{}, err
	}

	memcache.SetDefault(key, b)
	return b, nil
}

func GetBoardByCode(code string) (Board, error) {
	key := "GetBoardByCode-" + code
	cb, found := memcache.Get(key)
	if found {
		return cb.(Board), nil
	}

	r := db.QueryRow(`SELECT * FROM board WHERE code = ?;`, code)

	var b Board
	err := r.Scan(&b.Id, &b.Code, &b.Name)
	if err != nil {
		return Board{}, err
	}

	memcache.SetDefault(key, b)
	return b, nil
}

func GetAllBoards() ([]Board, error) {
	key := "GetAllBoards"
	cb, found := memcache.Get(key)
	if found {
		return cb.([]Board), nil
	}

	r, err := db.Query(`SELECT * FROM board;`)
	if err != nil {
		return []Board{}, err
	}
	defer r.Close()

	var bs []Board
	for r.Next() {
		var b Board
		err := r.Scan(&b.Id, &b.Code, &b.Name)
		if err != nil {
			return []Board{}, err
		}

		bs = append(bs, b)
	}

	memcache.SetDefault(key, bs)
	return bs, nil
}

func (b *Board) GetAllPosts() ([]Post, error) {
	key := "GetAllPosts-" + b.Code
	cb, found := memcache.Get(key)
	if found {
		return cb.([]Post), nil
	}

	r, err := db.Query(`SELECT * FROM post WHERE board_parent_id = ? AND thread_parent_id IS NULL;`, b.Id)
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

	// keep database access bounded under heavy load
	memcache.Set(key, ps, time.Second)
	return ps, nil
}
