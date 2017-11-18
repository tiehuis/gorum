package model

import (
	"database/sql"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/patrickmn/go-cache"
	"github.com/romana/rlog"

	"github.com/tiehuis/gorum/config"
)

const ModelVersion = 1

var db *sql.DB
var dbWMutex *sync.Mutex
var memcache *cache.Cache

func Init() {
	dbh, err := sql.Open("sqlite3", config.DatabaseAddress)
	if err != nil {
		rlog.Critical("Failed to open database:", err)
		os.Exit(1)
	}
	db = dbh

	rlog.Info("Opened database", config.DatabaseAddress)

	CheckExec("PRAGMA journal_mode = WAL;")
	CheckExec("PRAGMA synchronous = FULL;")
	CheckExec("PRAGMA foreign_keys = ON;")

	memcache = cache.New(5*time.Minute, time.Minute)
	rlog.Infof("Enabled memory cache with default expiration of %v", 5*time.Minute)
}

func Migrate() {
	sv := ReadConfig(Version)
	v, err := strconv.Atoi(sv)
	if err != nil {
		rlog.Critical("Found bad database version:", sv)
		os.Exit(1)
	}

	rlog.Info("Found existing database version:", v)

	if v > ModelVersion {
		rlog.Criticalf("Datbase version greater than binary version; Have %s but need %s", v, ModelVersion)
		os.Exit(1)
	}

	switch v {
	case 0:
		rlog.Info("Applying Migration0")
		migration0()
		WriteConfig(Version, "0")

		fallthrough
	case 1:
		WriteConfig(Version, "1")
	}

	if config.UseTestData {
		testdata()
	}
}

func mustPrepare(query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		rlog.Critical("Could not perpare query:", query, err)
		os.Exit(1)
	}
	return stmt
}
