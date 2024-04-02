package gsb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type APRPostgres struct {
	aprURL *url.URL
	Conn   *sqlx.DB
	TX     *sql.Tx
	ctx    context.Context
}

func (a *APRPostgres) connect() error {
	var err error

	password, _ := a.aprURL.User.Password()
	dbName := a.aprURL.Path[1:]
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", a.aprURL.User.Username(), dbName, password, a.aprURL.Host)
	log.Debugf("APRPostgres.connect.connectionString: %s", connectionString)

	a.Conn, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		return err
	}

	// defer a.DB.Close()

	return nil
}

func (a *APRPostgres) Reset() error {
	return a.connect()
}

func (a *APRPostgres) Init(aprURL *url.URL) error {
	a.aprURL = aprURL
	return a.connect()
}

func (a *APRPostgres) Exec(sql string) (sql.Result, error) {
	log.Debugf("APRPostgres.Exec.sql, %s", sql)
	return a.Conn.ExecContext(a.ctx, sql)
}

func (a *APRPostgres) QueryRow(sql string) *sql.Row {
	log.Debugf("APRPostgres.QueryRow.sql, %s", sql)
	return a.Conn.QueryRowContext(a.ctx, sql)
}

func (a *APRPostgres) Query(sql string) (*sql.Rows, error) {
	log.Debugf("APRPostgres.QueryRow.sql, %s", sql)
	return a.Conn.QueryContext(a.ctx, sql)
}

func (a *APRPostgres) Begin() error {
	log.Debug("APRPostgres.Begin")
	var err error

	a.ctx = context.Background()
	a.TX, err = a.Conn.BeginTx(a.ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	return err
}
func (a *APRPostgres) Rollback() error {
	log.Debug("APRPostgres.Rollback")
	return a.TX.Rollback()
}
func (a *APRPostgres) Commit() error {
	log.Debug("APRPostgres.Commit")
	return a.TX.Commit()
}
