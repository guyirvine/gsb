package main

import (
	"database/sql"
	"fmt"
	"gsb"
	"os"

	log "github.com/sirupsen/logrus"
)

type ExampleDbHandler struct {
	Db *gsb.APRPostgres
}

type ExampleDbMessage struct {
	Label string
}

func (m *ExampleDbMessage) GetPayload() string {
	return m.Label
}

func (h *ExampleDbHandler) GetMessage() gsb.Message {
	return &ExampleDbMessage{}
}

func (h *ExampleHandler) Init() error {
	return nil
}

func (h *ExampleDbHandler) Handle(msg gsb.Message) error {
	var label string
	var err error
	var rows *sql.Rows

	_, err = h.Db.Exec("UPDATE test_tbl SET label = label||'1'")
	if err != nil {
		return err
	}

	err = h.Db.QueryRow("SELECT label FROM test_tbl").Scan(&label)
	if err != nil {
		return err
	}
	fmt.Printf("*******   ExampleHandler.QueryRow.label: %s\n", label)

	rows, err = h.Db.Query("SELECT label FROM test_tbl")
	if err != nil {
		return err
	}

	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		if err := rows.Scan(&label); err != nil {
			return err
		}
		fmt.Printf("*******   ExampleHandler.QueryRows.Next.label: %s\n", label)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	// os.Setenv("GSB_APR_Db", "postgres://admin:pass@0.0.0.0/testdb")
	os.Setenv("GSB_MQ", "beanstalk://")
	// os.Setenv("GSB_MQ", "inmem://")
	log.SetLevel(log.InfoLevel)
	// log.SetLevel(log.DebugLevel)

	host := new(gsb.Host)
	host.Init()
	host.LoadHandler(new(ExampleDbHandler))

	host.Send(&ExampleDbMessage{"ttttt"})
	host.Start()

}
