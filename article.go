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
	Contnet string 
	date 	string 
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
template,_ := template.ParseFiles("Home.html")
session, _ := store.Get(Request, "cookie-name")
auth, ok := session.Values["authenticated"].(bool)
if  !ok || !auth {
	http.Error(Response, "Forbidden", http.StatusForbidden)
	return
}
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
	template_name = "login.html"
	name := Request.FormValue("username")
	password := Request.FormValue("Password")
	flag := login(name,password)

	if flag{
		template_name = "Home.html"
		session, _ := store.Get(Request, "cookie-name")
		session.Values["authenticated"] = true
		session.Save(Request, Response)
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
log in the user if he/she has an existing account 
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
	http.HandleFunc("/Signup", SignupHandler)
	http.HandleFunc("/login",  LoginHandler )
	http.HandleFunc("/logout",  LogoutHandler )

	http.ListenAndServe(":8000",nil)

}
