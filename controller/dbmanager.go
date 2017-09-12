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
// Set and test connection
// Execute mysql dump
// Setup structs

func TestGorp() *gorp.DbMap {
	return initDb()
	//getMySQL()
}
func getMySQL() string {
	b, err := ioutil.ReadFile("dev/mtest_users.sql") // just pass the file name
	if err != nil {
		fmt.Println("ERROR Getting SQL file: " + err.Error())
	}
	return string(b) // convert content to a 'string'
}

func initDb() *gorp.DbMap {
	db, err := createAndOpen("gorpdb")
	if err != nil {
		fmt.Printf("ERROR create and open db: %s \n", err)
		//return nil
	}
	//checkErr(err, "sql.Open failed")
	fmt.Println("Getting DBmap")
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "MyISAM", Encoding: "utf8"}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	err = dbmap.Insert(&Post{Body: "test", Title: "title"})
	if err != nil {
		fmt.Println("ERROR: try insert: ", err.Error())
	}
	dbmap.Exec("DROP TABLE IF EXISTS users2")
	table := dbmap.AddTableWithName(UserAuth{}, "users2")
	dbmap.CreateTablesIfNotExists()
	if table == nil {
		fmt.Println("Empty table pointer")
	}
	err = dbmap.Insert(&UserAuth{Name: "TestName", LastAccess: "yesterday", UserLogin: "Login"})
	if err != nil {
		fmt.Println("ERROR: try insert: ", err.Error())
	}

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
		fmt.Println("Trying to connect to next connection")
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
	fmt.Println("Creating database if not exists")
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		fmt.Printf("Err creating db: %s", err)
		err = nil
		// panic(err)
	}
	db.Close()
	fmt.Println("Update connect -> directly use selected database " + dbConnect)
	db, err = sql.Open("mysql", dbConnect+name)
	if err != nil {
		fmt.Printf("ERROR openning new connection to database: %s", err)
		// panic(err)
	} /*
		if _, err := db.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
			return nil, errors.Wrapf(err, "Error setting fk")
		}*/
	fmt.Println("Trying to set mysql dump")
	_, err = db.Exec(getMySQL())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Initialization finished")
	//db.Exec("SET FOREIGN_KEY_CHECKS=1")
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
