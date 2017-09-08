package controller

import (
	"fmt"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"io/ioutil"
	"mtest/common/errors"
)

const (
	connect1 = "root:12345678@tcp(127.0.0.1:3306)/"
	connect2 = "root:111@tcp(127.0.0.1:3306)/"
)


func TestGorp() {
	initDb()
	//getMySQL()
}
func getMySQL() string {
	b, err := ioutil.ReadFile("dev/mtest_users.sql") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return string(b) // convert content to a 'string'
}

func initDb() *gorp.DbMap {
	db, err := createAndOpen("gorpdb")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "MyISAM", Encoding: "utf8"}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return dbmap
}

func createAndOpen(name string) (*sql.DB, error) {
	var dbConnect string

	db, err := sql.Open("mysql", connect2)
	if err != nil {
		db, err = sql.Open("mysql", connect1)
		if err != nil {
			return nil, errors.Wrapf(err, "Error connecting to DB")
		}
		fmt.Println("Connected to [1]:" + connect1)
		dbConnect = connect1
		// panic(err)
	} else {
		fmt.Println("Connected to [2]:" + connect2)
		dbConnect = connect2
	}
	// defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		// panic(err)
	}
	// db.Close()

	db, err = sql.Open("mysql", dbConnect+name)
	if err != nil {
		// panic(err)
	}
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return nil, err
	}
	_, err = db.Exec(getMySQL())
	if err != nil {
		fmt.Println(err)
	}
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
	// defer db.Close()
	return db, err
}

type Post struct {
	// db tag lets you specify the column name if it differs from the struct field
	Id      int64 `db:"post_id"`
	Created int64
	Title   string `db:",size:50"`               // Column size set to 50
	Body    string `db:"article_body,size:1024"` // Set both column name and size
}
