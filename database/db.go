package database

type DB struct {
	ConnStr string
}

func NewDB(connstr string) DB {
	return DB{
		ConnStr: connstr,
	}
}

// add db run and return error
func (d *DB) Run() {
}
