package server

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tenntenn/sqlite"
)

func (ew *EventWatcher) Start() error {
	if err := ew.InitDB(context.Background()); err != nil {
		return err
	}
	ew.InitHandlers()
	if err := ew.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

//DB 追加処理.
func (ew *EventWatcher) AddCondition(ctx context.Context, c *Users) error {
	const sqlStr = `INSERT INTO conditions(username, pass) VALUES (?,?);`
	r, err := ew.db.ExecContext(ctx, sqlStr, c.Username, c.Pass)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

//DB 初期化.
func (ew *EventWatcher) InitDB(ctx context.Context) error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS conditions(
		id	INTEGER PRIMARY KEY,
		Username 	TEXT NOT NULL,
		Pass 	TEXT NOT NULL
	);`

	if _, err := ew.db.ExecContext(ctx, sqlStr); err != nil {
		return err
	}

	return nil
}

//接続処理.
func New(addr string) (*EventWatcher, error) {
	mux := http.NewServeMux()
	db, err := sql.Open(sqlite.DriverName, "user.db")
	if err != nil {
		return nil, err
	}

	return &EventWatcher{
		mux:    mux,
		db:     db,
		server: &http.Server{Addr: addr, Handler: mux},
	}, nil
}

//DB 取得処理.
func (ew *EventWatcher) Conditions(ctx context.Context, limit int) ([]*Users, error) {
	const sqlStr = `SELECT id, username, pass FROM conditions LIMIT ?`
	rows, err := ew.db.QueryContext(ctx, sqlStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() 

	var cs []*Users
	for rows.Next() {
		var c Users
		err := rows.Scan(&c.ID, &c.Username, &c.Pass)
		if err != nil {
			return nil, err
		}
		cs = append(cs, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cs, nil
}
