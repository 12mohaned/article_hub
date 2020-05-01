/*
Class Responbile for Connection of Database and queries
*/

package main

import (
	"database/sql"
	_"github.com/lib/pq"
	"fmt"
	"article_hub/article"
)
var (
	Username string
	firstname string
	lastname string
	email string
	password string 
	Title string
	Content string
	date 	string 
)
const ( 
	host ="localhost"
	port = 5432
	user ="postgres"
	Password = ""
	dbname = "article_hub";
)

type DataBase interface{
	Connection()
	Login()bool 
	Signup()bool
	Write_article()
	getArticles() []Article
}

	/*
log the user in if he/she has an existing account 
*/
func Login(username string , password string)bool{
	postgresconnection := "user="+user+ " " + "password=" +Password + " " + "dbname="+dbname + " " + "sslmode=disable"
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
	  }	
	  sqlStatement := `Select username from users where username= $1 and password = $2`
	  //hashedpassword := hash(password)
	  row,err := db.Query(sqlStatement,username,password)
	  if(err != nil){
		  panic(err)

	  }
	  flag := row.Next()
	  return flag
	}

	/*
Signup function responsible for registering a user into the database
*/
func Signup(username string, firstname string, lastname string, email string, password string)bool{
	var flag bool
	postgresconnection := "user="+user+ " " + "password=" +Password + " " + "dbname="+dbname + " " + "sslmode=disable"
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
		}	
		  
	sqlStatement := `INSERT INTO users (username, firstname, lastname, email,password)
				VALUES ($1, $2, $3, $4,$5)`
	row, err := db.Exec(sqlStatement, username,firstname,lastname,  email,password)	
	
	if err != nil{
		panic(err)
		fmt.Println(row)
	}
	flag = true
	return flag
		}

/*
Add the Article in the DataBase
*/
func Write_article(title string,content string,date string,username string){
	postgresconnection := "user="+user+ " " + "password=" +Password + " " + "dbname="+dbname + " " + "sslmode=disable"
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
		}	
		sqlStatement := `insert into article (title , content, date,username)
			 values ($1,$2,$3,$4)`
		row,err := db.Query(sqlStatement, title,content,date,username)
		if err != nil {
			panic(err)
			fmt.Println(username)
			fmt.Println(row)
		}
		}

/* Return the Article of a user*/
func getArticles(username string)[]Article{
	postgresconnection := "user="+user+ " " + "password=" +Password + " " + "dbname="+dbname + " " + "sslmode=disable"
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
	  }	
	  sqlStatement := `select title, content, date from article where username=$1`
	  rows,err := db.Query(sqlStatement, username)
	 if err != nil{
		panic(err) 
	 }
	 var Articles []Article
	 defer rows.Close()
	 for rows.Next(){
		 var title string
		 var content string
		 var date string
	
		data := rows.Scan(&title, &content, &date)
		article := Article{Title:title, Content:content, date:date, user:username}
		if data != nil{
		}
		Articles = append(Articles,article)
	}
	return Articles
}
/*
Hashing Password in the Database
*/
func hash(s string) string{
	h := sha256.New()
	h.Write([]byte(s))
	hashstring,_ := fmt.Printf("%x", h.Sum(nil))
	return string(hashstring)
}

func main(){
D1 := DataBase{}
}