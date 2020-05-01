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

/* Cookie Set-up and Information*/
var (
	key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)	
)

const ( 
	host ="localhost"
	port = 5432
	user ="postgres"
	Password = ""
	dbname = "article_hub";
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
date := time.Now()
DateFormat := date.Format("01-02-2006 15:04:05 Monday")
name := session.Values["username"]
data := name.(string)
if flag{
	Write_article(Title,Content,DateFormat,data)
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
	data := name.(string)
	getArticles(data)
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


func LogoutHandler(Response http.ResponseWriter,Request *http.Request){

	session, _ := store.Get(Request, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(Request, Response)
	template,_ := template.ParseFiles("Logout.html")
	template.Execute(Response,nil)
}

/*
Validate the Article which is going to Posted by the User
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

func main() {
	http.HandleFunc("/Home", HomeHandler)
	http.HandleFunc("/profile",ProfileHandler)
	http.HandleFunc("/Signup", SignupHandler)
	http.HandleFunc("/login",  LoginHandler )
	http.HandleFunc("/logout",  LogoutHandler )
	http.ListenAndServe(":8000",nil)
}
