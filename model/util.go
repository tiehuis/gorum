package model

import "log"

func CheckExec(q string, args ...interface{}) {
	_, err := db.Exec(q, args...)
	if err != nil {
		log.Fatalf("failure on query: %s\n%s", q, err)
	}
}
