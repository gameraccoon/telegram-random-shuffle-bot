package database

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Database struct {
	// connection
	conn *sql.DB
}

func sanitizeString(input string) string {
	return strings.Replace(input, "'", "''", -1)
}

func (database *Database) execQuery(query string) {
	_, err := database.conn.Exec(query)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (database *Database) Connect(fileName string) error {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	database.conn = db

	database.execQuery("PRAGMA foreign_keys = ON")

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" global_vars(name TEXT PRIMARY KEY" +
		",integer_value INTEGER" +
		",string_value STRING);")

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" users(id INTEGER NOT NULL PRIMARY KEY" +
		",chat_id INTEGER UNIQUE NOT NULL" +
		")")

	database.execQuery("CREATE UNIQUE INDEX IF NOT EXISTS" +
		" chat_id_index ON users(chat_id)")

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" lists(id INTEGER NOT NULL PRIMARY KEY" +
		",user_id INTEGER" +
		",name STRING" +
		",last_item_id INTEGER NOT NULL" +
		",FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL" +
		")")

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" list_items(id INTEGER NOT NULL PRIMARY KEY" +
		",list_id INTEGER" +
		",text STRING" +
		",FOREIGN KEY(list_id) REFERENCES lists(id) ON DELETE SET NULL" +
		")")

	return nil
}

func (database *Database) Disconnect() {
	database.conn.Close()
	database.conn = nil
}

func (database *Database) IsConnectionOpened() bool {
	return database.conn != nil
}

func (database *Database) createUniqueRecord(table string, values string) int64 {
	var err error
	if len(values) == 0 {
		_, err = database.conn.Exec(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES ", table))
	} else {
		_, err = database.conn.Exec(fmt.Sprintf("INSERT INTO %s VALUES (%s)", table, values))
	}

	if err != nil {
		log.Fatal(err.Error())
		return -1
	}

	rows, err := database.conn.Query(fmt.Sprintf("SELECT id FROM %s ORDER BY id DESC LIMIT 1", table))

	if err != nil {
		log.Fatal(err.Error())
		return -1
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err.Error())
			return -1
		}

		return id
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal("No record created")
	return -1
}

func (database *Database) GetDatabaseVersion() (version string) {
	rows, err := database.conn.Query("SELECT string_value FROM global_vars WHERE name=\"version\"")

	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		// that means it's a new clean database
		version = latestVersion
	}

	return
}

func (database *Database) SetDatabaseVersion(version string) {
	database.execQuery("DELETE FROM global_vars WHERE name='version'")
	database.execQuery(fmt.Sprintf("INSERT INTO global_vars (name, string_value) VALUES ('version', '%s')", sanitizeString(version)))
}

func (database *Database) GetUserId(chatId int64) (userId int64) {
	database.execQuery(fmt.Sprintf("INSERT OR IGNORE INTO users(chat_id) "+
		"VALUES (%d)", chatId))

	rows, err := database.conn.Query(fmt.Sprintf("SELECT id FROM users WHERE chat_id=%d", chatId))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("No user found")
	}

	return
}

func (database *Database) GetUserChatId(userId int64) (chatId int64) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT chat_id FROM users WHERE id=%d", userId))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&chatId)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("No user found")
	}

	return
}

func (database *Database) GetUserLists(userId int64) (ids []int64, texts []string) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT id, name FROM lists WHERE user_id=%d", userId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var text string

		err := rows.Scan(&id, &text)
		if err != nil {
			log.Fatal(err.Error())
		}

		ids = append(ids, id)
		texts = append(texts, text)
	}

	return
}

func (database *Database) GetListItems(listId int64) (ids []int64, texts []string) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT id, text FROM list_items WHERE list_id=%d", listId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var text string

		err := rows.Scan(&id, &text)
		if err != nil {
			log.Fatal(err.Error())
		}

		ids = append(ids, id)
		texts = append(texts, text)
	}

	return
}

func (database *Database) GetLastItem(listId int64) (item_id int64) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT last_item_id FROM lists WHERE id=%d", listId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&item_id)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("No list found")
	}

	return
}

func (database *Database) SetLastItem(listId int64, lastItemId int64) {
	database.execQuery(fmt.Sprintf("UPDATE lists SET last_item_id=%d WHERE id=%d", lastItemId, listId))
}

func (database *Database) CreateList(userId int64, name string) {
	database.execQuery(fmt.Sprintf("INSERT INTO lists (user_id, name, last_item_id) VALUES (%d, '%s', -1)", userId, name))
}

func (database *Database) AddItemsToList(listId int64, items []string) {
	var buffer bytes.Buffer

	for _, item := range items {
		buffer.WriteString(fmt.Sprintf("INSERT INTO list_items (list_id, text) VALUES (%d, '%s');", listId, item))
	}

	database.execQuery(fmt.Sprintf("BEGIN TRANSACTION;" +
		"%s" +
		"COMMIT;", buffer.String()))
}

func (database *Database) RemoveItem(itemId int64) {
	database.execQuery(fmt.Sprintf("DELETE FROM list_items WHERE id=%d", itemId))
}

func (database *Database) DeleteList(listId int64) {
	database.execQuery(fmt.Sprintf("DELETE FROM lists WHERE id=%d", listId))
}
