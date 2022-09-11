package sq

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SaveMessage struct {
	messageId int
	tag       string
	value     string
}

type DBsq struct {
	db *sql.DB
}

func InitDB(filepath string) DBsq {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return DBsq{
		db: db,
	}
}

func (d *DBsq) Close() {
	d.db.Close()
}

func (d *DBsq) CreateTable() {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS savemessage(
		Id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		tag TEXT,
		value TEXT,
		messageId INTEGER
	);
	`

	_, err := d.db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func (d *DBsq) storeItem(items []SaveMessage) error {
	sql_additem := `
	INSERT OR REPLACE INTO savemessage(
		tag,
		value,
		messageId
	) values(?, ?, ?)
	`

	stmt, err := d.db.Prepare(sql_additem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.tag, item.value, item.messageId)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func (d *DBsq) getMessage(tag string) ([]SaveMessage, error) {
	sql_readall := `
	SELECT tag, value, messageId FROM savemessage WHERE tag = ?
	`

	rows, err := d.db.Query(sql_readall, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SaveMessage
	for rows.Next() {
		item := SaveMessage{}
		err2 := rows.Scan(&item.tag, &item.value, &item.messageId)
		if err2 != nil {
			return nil, err2
		}
		result = append(result, item)
	}
	return result, nil
}

func (d *DBsq) ManageMessage(command string, tagStr string, msgStr string, messageId int) (string, error) {

	if command == "сохрани" {
		err := d.setMessage(tagStr, msgStr, messageId)
		if err != nil {
			return "", err

		}
	}
	if command == "покажи" {
		MessageAray, err := d.getMessage(tagStr)
		if err != nil {
			return "", err
		}
		msg := formatSaveMessage(MessageAray)
		return msg, nil
	}
	return "", nil
}

func (d *DBsq) setMessage(tagStr string, msgStr string, messageId int) error {
	s := SaveMessage{
		messageId: messageId,
		tag:       tagStr,
		value:     msgStr,
	}
	var a []SaveMessage
	a = append(a, s)
	err := d.storeItem(a)
	if err != nil {
		return err
	}
	return nil
}

func formatSaveMessage(sm []SaveMessage) string {
    m:=""
	for _,a := range sm {
       m= m+a.value+"\n" 
	}
	return m
}
