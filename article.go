package main

import (
	"html/template"
	"net/http"
	"regexp"
	"fmt"
	"crypto/sha256"
	"database/sql"
	_"github.com/lib/pq"
	"github.com/gorilla/sessions"
	"time"
)
/* DataBase connections and paramter's
*/
const ( 
	host ="localhost"
	port = 5432
	user ="postgres"
	Password = ""
	dbname = "article_hub";
)
/* Cookie Set-up and Information*/
var (
	key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)	
)

type Article struct {
	Title string
	Content string 
	date 	string 
	user string
}

type Users struct {
	UserName string
	FirstName string
	LastName string
	Email string
	password string 
}
/*
Responsible for Home Page Requests
*/
func HomeHandler(Response http.ResponseWriter, Request *http.Request){

session, _ := store.Get(Request, "cookie-name")
auth, ok := session.Values["authenticated"].(bool)
if  !ok || !auth {
	http.Error(Response, "You Should Log in to Access this page", http.StatusForbidden)
	return
}

Request.ParseForm()
Title := Request.FormValue("Title")
Content := Request.FormValue("Content")
flag := ArticleValidation(Title,Content)

if flag{
	date := time.Now()
	DateFormat := date.Format("01-02-2006 15:04:05 Monday")
	name := session.Values["username"]
	data,_ := fmt.Printf("%x",name)
	write_article(Title,Content,DateFormat,string(data))
}

template,_ := template.ParseFiles("Home.html")
template.Execute(Response,nil)
}

/*
Route to the User profile which contain the Articles he wrote
*/
func ProfileHandler(Response http.ResponseWriter, Request *http.Request){
	session, _ := store.Get(Request, "cookie-name")
	auth, ok := session.Values["authenticated"].(bool)
	// var Articles []Article
	if  !ok || !auth {
		http.Error(Response, "You Should Log in to Access this page", http.StatusForbidden)
		return
	}
	name := session.Values["username"]
	data,_ := fmt.Printf("%x",name)
	getArticles(string(data))
	template,_ := template.ParseFiles("profile.html")
	template.Execute(Response,nil)
}
/*
Page where user Sign up
*/
func SignupHandler(Response http.ResponseWriter, Request *http.Request){
	Request.ParseForm()
	var flag bool
	var template_name string
	template_name = "Signup.html"
	name := Request.FormValue("username")
	FirstName := Request.FormValue("FirstName")
	LastName  := Request.FormValue("LastName")
	Email     := Request.FormValue("Email")
	password  := Request.FormValue("Password")

	user := Users {UserName :name ,FirstName:FirstName ,LastName:LastName, Email:Email, password : password}
	verified := Signupvalidation(user.UserName,user.FirstName,user.LastName,user.Email,user.password)
	if(verified == "valid"){
		
		//hashedpassword := hash(user.password)
		flag = singup(user.UserName, user.FirstName, user.LastName, user.Email, user.password)
		if flag {
			template_name = "Home.html"
			session, _ := store.Get(Request, "cookie-name")
			session.Values["authenticated"] = true
			session.Values["username"] = name
			session.Save(Request, Response)
		}
	}
	template,_ := template.ParseFiles(template_name)
	template.Execute(Response,flag)
	}

	/*
	Responsible for Logging in Page Requests
	*/
func LoginHandler(Response http.ResponseWriter,Request *http.Request){

	var template_name string
	var flag bool
	template_name = "login.html"
	name := Request.FormValue("username")
	password := Request.FormValue("Password")
		if loginvalidation(name,password) != false {
		flag := login(name,password)

	if flag{
		template_name = "Home.html"
		session, _ := store.Get(Request, "cookie-name")
		session.Values["authenticated"] = true
		session.Values["username"] = name
		session.Save(Request, Response)
	}
		}
	template,_ := template.ParseFiles(template_name)
	template.Execute(Response,flag)
}


func LogoutHandler(Response http.ResponseWriter,Request *http.Request){

	session, _ := store.Get(Request, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(Request, Response)
	template,_ := template.ParseFiles("Logout.html")
	template.Execute(Response,nil)
}

func ArticleValidation(Title string, Content string)bool{
	TitleValidation,_ := regexp.MatchString("[a-zA-Z0-9]{1,}",Title)
	ContentValidation,_ := regexp.MatchString("[a-zA-Z0-9]{1,}",Content)
	
	if TitleValidation != true{
		return false
	}

	if ContentValidation != true{
		return false
	}
	return true
}
/*
Validate the User input when trying to sign-up
*/
func Signupvalidation(username string, FirstName string, LastName string, Email string, password string)string {
	usernamevalidation, _ := regexp.MatchString("[a-zA-Z0-9]{3,20}",username)
	FirstNamevalidation, _ := regexp.MatchString("[a-zA-Z0-9]{3,20}",FirstName)
	LastNamevalidation, _ := regexp.MatchString("[a-zA-Z0-9]{3,20}",LastName)
	EmailValidation, _ := regexp.MatchString("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+[a-zA-Z0-9-.]+$",Email)
	PasswordValidation, _ := regexp.MatchString("[a-zA-Z0-9]{10,}",password) 
	
	if usernamevalidation != true {
		return "userName is not valid"
	}

	if FirstNamevalidation != true {
		return "First Name is not valid"
	}

	if LastNamevalidation != true{
		return "Last Name is not valid"
	}

	if EmailValidation != true {
		return "Email is not valid"
	}

	if PasswordValidation != true{
		return "Password is not valid"
	}

	return "valid"; 
}


/* Validating the User Input when logging in
*/
func loginvalidation(username string, password string)bool{
	usernamevalidation, _ := regexp.MatchString("[a-zA-Z0-9]{3,20}",username)
	PasswordValidation, _ := regexp.MatchString("[a-zA-Z0-9]{10,}",password) 
	if(usernamevalidation != true){
		return false
	}
	if(PasswordValidation != true){
		return false
	}
	return true
}

/*
Signup function responsible for registering a user into the database
*/
func singup(username string, firstname string, lastname string, email string, password string)bool{
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
log the user in if he/she has an existing account 
*/
func login(username string , password string)bool{
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
Add the Article in the DataBase
*/
func write_article(title string,content string,date string,username string){
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
		fmt.Println(row)
	  }
}
/* Return the Article of a user*/
func getArticles(username string) []Article{
	postgresconnection := "user="+user+ " " + "password=" +Password + " " + "dbname="+dbname + " " + "sslmode=disable"
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
	  }	
	  sqlStatement := `select title, content, date from article where username=$1`
	  rows,err := db.Query(sqlStatement, username)
	 if err != nil{
		panic(err) 
		fmt.Println(rows)
	 }
	 var Articles []Article 
	 i:= 0
	 defer rows.Close()
	 for rows.Next(){
		 var title string
		 var content string
		 var date string
		data := rows.Scan(&title, &content, &date)
		if data != nil{

		}
		fmt.Println(content)

		i++
	 }
if(i > 0){
	return Articles
}
return nil
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

func main() {
	http.HandleFunc("/Home", HomeHandler)
	http.HandleFunc("/profile",ProfileHandler)
	http.HandleFunc("/Signup", SignupHandler)
	http.HandleFunc("/login",  LoginHandler )
	http.HandleFunc("/logout",  LogoutHandler )
	http.ListenAndServe(":8000",nil)
}
