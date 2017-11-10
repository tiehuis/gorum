package model

const (
	Version = "version"
)

func ReadConfig(key string) string {
	row := db.QueryRow(`SELECT val FROM config WHERE name=?;`, key)
	var val string
	if err := row.Scan(&val); err == nil {
		return val
	}
	return "0"
}

func WriteConfig(key, val string) {
	var old string
	if db.QueryRow(`SELECT val FROM config WHERE name=?;`, key).Scan(&old) == nil {
		if old != val {
			CheckExec(`UPDATE config SET val=? WHERE name=?;`, val, key)
		}
	} else {
		CheckExec(`INSERT INTO config(name, val) values(?, ?);`, key, val)
	}
}
