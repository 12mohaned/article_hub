package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	_"github.com/lib/pq"
	"github.com/gorilla/sessions"
	"time"
	"regexp"
	 "github.com/gorilla/mux"
	)
/* Cookie Set-up and Information*/
var (
	key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)	
)

const ( 
	host ="localhost"
	port = 5432
	user ="postgres"
	Password = "tarekandamr12/"
	dbname = "article_hub";
)

type Article struct {
	Title string
	Content string 
	date string 
	user string
}

type Users struct {
	UserName string
	FirstName string
	LastName string
	Email string
	password string 
}

type Articles struct{
	Articles []Article
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
template,_ := template.ParseFiles("Home.html")
template.Execute(Response,nil)
}
/**
* ProfileHanlder
* * the Profile View where users can see their Articles
* Routing the User when editing the post
*/
func ProfileHandler(Response http.ResponseWriter, Request *http.Request){
	var template_name string
	template_name = "profile.html"
	// var template_name string
	session, _ := store.Get(Request, "cookie-name")
	auth, ok := session.Values["authenticated"].(bool)	
	if  !ok || !auth {
		http.Error(Response, "You Should Log in to Access this page", http.StatusForbidden)
		return
	}
	if Request.Method == "GET"{
	
	name := session.Values["username"]
	data := name.(string)
	articles := Articles{Articles : getArticles(data)}
	vars := mux.Vars(Request)
	title := vars["title"]
	flag := Checktitle(title,data)

	//check if the title of article is found in the database or not 
	if len(flag) != 0{
		template_name = "post.html"
		articles = Articles {Articles:flag}
		Request.ParseForm()
		Title := Request.FormValue("Title")
		Content := Request.FormValue("Content")
		if len(Title) > 0{
			updateTitle(title,Title)
		}

		if len(Content) > 0{
			updateContent(Content,Title)
		}
	}
	template,_ := template.ParseFiles(template_name)
	template.Execute(Response,articles)
}
}


func WriteArticleHandler(Response http.ResponseWriter,Request *http.Request){
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
	date := time.Now()
	DateFormat := date.Format("01-02-2006 15:04:05 Monday")
	name := session.Values["username"]
	data := name.(string)
	if flag{
		WriteArticle(Title,Content,DateFormat,data)
	}
	template,_ := template.ParseFiles("writeArticle.html")
	template.Execute(Response,nil)
}
/**
*  SignupHandler
* * Responsible for Rendering the Signup Form 
* * Redirect User to Home Page in case of valid input

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
	if verified{
		//hashedpassword := hash(user.password)
		flag = Signup(user.UserName, user.FirstName, user.LastName, user.Email, user.password)
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

/** 
*  LoginHandler
* * Responsible for logging users in and redirect 
* * them to Home page in case of correct Authentication 

*/
func LoginHandler(Response http.ResponseWriter,Request *http.Request){

	var template_name string
	var flag bool
	template_name = "login.html"
	name := Request.FormValue("username")
	password := Request.FormValue("Password")
		if loginvalidation(name,password) != false {
		flag := Login(name,password)

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

/** 
* ArticleValidation
* * Validating the Article Input when the user is trying to post an Article
*/
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

/** 
* Signupvalidation
* * Validating the User Input (username,password) when attempting to Create an Account
*/
func Signupvalidation(username string, FirstName string, LastName string, Email string, password string)bool {
	usernamevalidation, _ := regexp.MatchString("[a-zA-Z0-9]{3,20}",username)
	FirstNamevalidation, _ := regexp.MatchString("[a-zA-Z0-9]{3,20}",FirstName)
	LastNamevalidation, _ := regexp.MatchString("[a-zA-Z0-9]{3,20}",LastName)
	EmailValidation, _ := regexp.MatchString("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+[a-zA-Z0-9-.]+$",Email)
	PasswordValidation, _ := regexp.MatchString("[a-zA-Z0-9]{10,}",password) 
	
	if usernamevalidation != true {
		return false
	}

	if FirstNamevalidation != true {
		return false
	}

	if LastNamevalidation != true{
		return false
	}

	if EmailValidation != true {
		return false
	}

	if PasswordValidation != true{
		return false
	}

	return true; 
}

/** 
* loginvalidation
* * Validating the User Input (username,password) when attempting to log in
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

/** 
*  LogoutHandler 
* * Responsible for logging users out 
*/
func LogoutHandler(Response http.ResponseWriter,Request *http.Request){
	session, _ := store.Get(Request, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(Request, Response)
	template,_ := template.ParseFiles("Logout.html")
	template.Execute(Response,nil)
}


	/*
log the user in if he/she has an existing account 
*/
func Login(username string , password string)bool{
	postgresconnection := initConnection()
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
	  }	
	  sqlStatement := `Select username from Users where username= $1 and password = $2`
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
	postgresconnection := initConnection()
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
		}	
		  
	sqlStatement := `INSERT INTO Users (username, firstname, lastname, email,password)
				VALUES ($1, $2, $3, $4,$5)`
	row, err := db.Exec(sqlStatement, username,firstname,lastname,  email,password)	
	
	if err != nil{
		panic(err)
		fmt.Println(row)
	}
	flag = true
	return flag
		}

	func initConnection() string{
		postgresconnection := "user="+user+ " " + "password=" +Password + " " + "dbname="+dbname + " " + "sslmode=disable"
		return postgresconnection
		}

//Write Article
func WriteArticle(title string,content string,date string,username string){
	postgresconnection := initConnection()
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
		}	
		sqlStatement := `insert into Article (Title , Content, date,username)
			 values ($1,$2,$3,$4)`
		row,err := db.Query(sqlStatement, title,content,date,username)
		if err != nil {
			panic(err)
			fmt.Println(username)
			fmt.Println(row)
		}
		fmt.Println(row)
		}
// /* Return the Article of a user*/
func getArticles(username string)[]Article{
	postgresconnection := initConnection()
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
	  }	
	  sqlStatement := `select title, content, date from Article where username=$1`
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

func Checktitle(title string,name string )[]Article{
	postgresconnection := initConnection()
	db,err := sql.Open("postgres",postgresconnection) 
	var Articles []Article
	if err != nil {
		panic(err)
	  }	
	  sqlStatement := `select username, title, content, date from Article where title=$1 and username = $2`
	  rows,err := db.Query(sqlStatement, title,name)
	 if err != nil{
		panic(err) 
	 }
	for rows.Next(){
	var Title string
	var content string
	var date string
	var username string 
	data := rows.Scan(&username,&Title, &content, &date)
	article := Article{Title:title, Content:content, date:date, user:username}
	if data != nil{
	}
	Articles = append(Articles,article)
	}
	return Articles
   }
// /* Update the Title of an Article*/
func updateTitle(oldtitle string, newtitle string)bool{
	postgresconnection := initConnection()
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
	  }	
	sqlStatement := `update article set title = $1 where title = $2`
	rows,err := db.Query(sqlStatement,newtitle,oldtitle)
	
	if err != nil{
		panic(err)
		fmt.Println(rows)
	}
	return true
}
// /* Update the Content of an Article*/
func updateContent(content string, title string)bool{
	postgresconnection := initConnection()
	db,err := sql.Open("postgres",postgresconnection) 
	if err != nil {
		panic(err)
	  }	
	sqlStatement := `update article set content = $1 where title = $2`
	rows,err := db.Query(sqlStatement,content,title)
	if err != nil{
		panic(err)
		fmt.Println(rows)
	}
	return true
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/Home",HomeHandler)
	mux.HandleFunc("/profile",ProfileHandler)
	mux.HandleFunc("/edit/{title}",ProfileHandler)
	mux.HandleFunc("/write",WriteArticleHandler)
	mux.HandleFunc("/Signup",SignupHandler)
	mux.HandleFunc("/login",LoginHandler )
	mux.HandleFunc("/logout",LogoutHandler )
	http.Handle("/",mux)
	http.ListenAndServe(":8000",nil)
}
