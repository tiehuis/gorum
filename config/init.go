package config

import (
	"flag"

	"github.com/romana/rlog"
)

var (
	DatabaseAddress  string
	HttpAddress      string
	PostTimeout      int64
	UseTestData      bool
	PostRateLimit    int64
	BoardThreadLimit int64
)

func Init() {
	flag.StringVar(&DatabaseAddress, "sqlite-db", "file::memory:?mode=memory&cache=shared", "Database name")
	flag.StringVar(&HttpAddress, "listen-address", ":8080", "Server address")
	flag.Int64Var(&PostTimeout, "post-timeout", 60, "Post timeout (seconds)")
	flag.BoolVar(&UseTestData, "testdata", false, "Seed database with test posts")
	flag.Int64Var(&PostRateLimit, "rate-limit", 60, "Post rate limit")
	flag.Int64Var(&BoardThreadLimit, "board-thread-limit", 15, "Board thread limit")
	flag.Parse()

	rlog.Debugf(`
DatabaseAddress  string: %v
HttpAddress      string: %v
PostTimeout      int64: %v
UseTestData      bool: %v
PostRateLimit    int64: %v
BoardThreadLimit int64: %v
`,
		DatabaseAddress,
		HttpAddress,
		PostTimeout,
		UseTestData,
		PostRateLimit,
		BoardThreadLimit)
}
