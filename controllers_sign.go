package main

import (
	"log"
	"net/http"
	"strings"

	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func LogInHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	gotFlash := make(map[string]interface{})
	if logged, err := session.Exists("UserId"); err != nil {
		fmt.Println(err)
	} else if logged {
		http.Redirect(w, r, "/", 302)
		return
	}
	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
		fmt.Println(err)
	}
	gotFlash["Title"] = "Log In"

	if err := tpl.ExecuteTemplate(w, "in.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func LogInProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	sentFlash := make(map[string]interface{})
	if logged, err := session.Exists("UserId"); err != nil {
		fmt.Println(err)
	} else if logged {
		http.Redirect(w, r, "/", 302)
		return
	}
	uORe := r.FormValue("uORe")
	sentFlash["UORE"] = uORe
	password := r.FormValue("password")
	var UserId int
	// check to see if the email or username is correct
	if UserId = checkUsername(uORe); UserId < 1 {
		sentFlash["Err"] = "Your User Name | Email Is Wrong."
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/in", 302)
		return
	}
	// username or email is correct
	cUser := getUserById(UserId)
	hashed := cUser.Password
	// check to see if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		sentFlash["Err"] = "your Password is wrong."
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/in", 303)
		return
	}
	// password is correct

	if affect := setUserOnline(UserId); affect < 1{
		fmt.Println("setting user to active mode failed.")
	}
	if err := session.PutInt(w, "UserId", cUser.Id); err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/profile/"+cUser.Name, 302)
}

func LogOutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	var time int
	var err error
	for k, v := range r.URL.Query() {
		if k == "time"{
			if time, err = strconv.Atoi(v[0]); err != nil{
				fmt.Println(err)
			}
		}
	}
	if userId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else {
		if affect := logOutProcess(userId, time); affect < 1 {
			fmt.Println("update last log didn't work")
		}
	}
	session.Destroy(w)
	w.Write([]byte("done"))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	gotFlash := make(map[string]interface{})
	if logged, err := session.Exists("UserId"); err != nil {
		fmt.Println(err)
	} else if logged {
		http.Redirect(w, r, "/", 302)
		return
	}
	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
		fmt.Println(err)
	}
	gotFlash["Title"] = "Register"

	if err := tpl.ExecuteTemplate(w, "register.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func RegisterProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	sentFlash := make(map[string]interface{})

	submit := r.FormValue("submit")
	username := r.FormValue("username")
	if submit == "CheckName" {
		// if submission is validating the username
		// first we validate it
		if er := nameValidation(username); er != "" {
			w.Write([]byte("Invalid"))
			return
		}
		// then we check to see if this name already exists in database

		// in the registration we don't have any userId data
		// so there is an err, but we ignore it, because when err
		// occurs the UserId gets 0 and we pass it to alreadyName func
		// and the value of 0 for userId means this user hasn't got any
		// previous username, so check all the userNames, but if the value
		// of userId is more than zero it means this user is changing its
		// username so in the alreadyName func we ignore him or herself name
		// in validation.
		UserId, _ := strconv.Atoi(r.FormValue("userId"))
		if affect := alreadyName(username, UserId); affect > 0 {
			w.Write([]byte("Taken"))
		} else {
			w.Write([]byte("Available"))
		}
		return
	}

	if logged, err := session.Exists("UserId"); err != nil {
		fmt.Println(err)
	} else if logged {
		http.Redirect(w, r, "/", 302)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	pwd := r.FormValue("password")
	confirm := r.FormValue("confirm")
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	password := string(hash)
	newUser := User{Name: strings.TrimSpace(username), FirstName: strings.TrimSpace(firstName), LastName: strings.TrimSpace(lastName) , Email: strings.TrimSpace(email), Password: strings.TrimSpace(password)}
	warnings := newUser.Validate(pwd, confirm)
	// if validation of form failed
	if len(warnings) > 0 {
		sentFlash["Errs"] = warnings
		sentFlash["Username"] = username
		sentFlash["Email"] = email
		sentFlash["FirstName"] = firstName
		sentFlash["LastName"] = lastName
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/register", 302)
		return
	}
	// if validation was passed
	UserId := insertUser(newUser)
	if err := session.PutInt(w, "UserId", UserId); err != nil {
		fmt.Println(err)
	}
	if err := session.PutString(w, "Cong", "Congratulations, You Successfully Registered In Blue Talk."); err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/profile/"+strconv.Itoa(UserId), 302)
}
