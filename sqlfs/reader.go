package sqlfs

import (
	"database/sql"
	"fmt"
	"io"
)

// Reader implements io.ReadCloser
type Reader struct {
	db    *sql.DB
	table string
	buf   []byte
	rows  *sql.Rows
}

// Open returns a reader to read from the given table in db.
func Open(db *sql.DB, table string) (*Reader, error) {
	has, e := HasTable(db, table)
	if !has {
		return nil, fmt.Errorf("Open: table %s doesn't exist", table)
	}
	if e != nil {
		return nil, fmt.Errorf("Open: HasTable failed with %v", e)
	}

	r := &Reader{
		db:    db,
		table: table,
		buf:   nil,
		rows:  nil}

	r.rows, e = r.db.Query(fmt.Sprintf("SELECT block FROM %s ORDER BY id", table))
	if e != nil {
		return nil, fmt.Errorf("Open: failed to query: %v", e)
	}
	return r, nil
}

func (r *Reader) Read(p []byte) (n int, e error) {
	if r.db == nil {
		return 0, fmt.Errorf("Read from a closed reader")
	}
	n = 0
	for n < len(p) {
		m := copy(p[n:], r.buf)
		n += m
		r.buf = r.buf[m:]
		if len(r.buf) <= 0 {
			if r.rows.Next() {
				e = r.rows.Scan(&r.buf)
				if e != nil {
					break
				}
			} else if n < len(p) {
				e = io.EOF
				break
			}
		}
	}
	return n, e
}

func (r *Reader) Close() error {
	if r.rows != nil {
		if e := r.rows.Close(); e != nil {
			return e
		}
		r.rows = nil
	}
	r.db = nil // Mark closed.
	return nil
}
