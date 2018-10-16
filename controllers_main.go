package main

import (
	"net/http"
	"strings"

	"fmt"
	"github.com/julienschmidt/httprouter"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	}
	gotFlash := make(map[string]interface{})
	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
		fmt.Println(err)
	}
	var err error
	if gotFlash["Cong"], err = session.PopString(w, "Cong"); err != nil {
		fmt.Println(err)
	}
	gotFlash["Title"] = "Blue Talk"
	gotFlash["Username"] = cUser.Name

	if err := tpl.ExecuteTemplate(w, "index.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

//func IndexProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	session := sessionManager.Load(r)
//	sentFlash := make(map[string]interface{})
//	// if the request was searching
//	if word := r.FormValue("search"); len(word) > 0 {
//		tags := searchTags(word)
//		stories := searchStories(word)
//		if len(tags) == 0 && len(stories) == 0 {
//			sentFlash["Err"] = "No Result Found For This Word."
//			sentFlash["Stories"] = getAllStories("", "")
//		} else {
//			// if the request was successful
//			sentFlash["Tags"] = tags
//			sentFlash["Stories"] = stories
//			sentFlash["LenS"] = len(stories)
//			sentFlash["LenT"] = len(tags)
//		}
//		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
//			fmt.Println(err)
//		}
//		http.Redirect(w, r, "/", 302)
//		return
//	}
//	// if the request was ordering
//	order := r.FormValue("order")
//	priority := r.FormValue("priority")
//	sentFlash["Stories"] = getAllStories(order, priority)
//	if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
//		fmt.Println(err)
//	}
//	http.Redirect(w, r, "/", 302)
//}

func AboutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	gotFlash := make(map[string]interface{})
	var cUser User
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	}
	gotFlash["Title"] = "About"
	gotFlash["Username"] = cUser.Name

	if err := tpl.ExecuteTemplate(w, "about.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func ContactHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	}
	gotFlash := make(map[string]interface{})
	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
		fmt.Println(err)
	}
	gotFlash["Title"] = "Contact"
	gotFlash["UserName"] = cUser.Name

	if err := tpl.ExecuteTemplate(w, "contact.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func ContactProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	sentFlash := make(map[string]interface{})
	var name, email string
	// to see if the user is already logged in, we don't ask for their information.
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser := getUserById(UserId)
		name = cUser.Name
		email = cUser.Email
	} else {
		name = r.FormValue("name")
		email = r.FormValue("email")
	}
	text := r.FormValue("text")
	newCnt := Cnt{Name: strings.TrimSpace(name), Email: strings.TrimSpace(email), Text: strings.TrimSpace(text)}
	if warnings := newCnt.Validate(); len(warnings) > 0 {
		sentFlash["Errs"] = warnings
		sentFlash["Name"] = name
		sentFlash["Email"] = email
		sentFlash["Text"] = text
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/contact", 302)
		return
	}
	if msgId := insertCnt(newCnt); msgId > 0 {
		sentFlash["Cong"] = "Congratulations dear " + name + ", your message has been sent successfully."
	} else {
		sentFlash["Err"] = "there is something wrong with your message, please check it and try again."
	}
	if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/contact", 302)
}

func SearchHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	}
	gotFlash := make(map[string]interface{})
	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
		fmt.Println(err)
	}
	gotFlash["Title"] = "Search"
	gotFlash["UserName"] = cUser.Name

	if err := tpl.ExecuteTemplate(w, "search.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func SearchProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessionManager.Load(r)
	sentFlash := make(map[string]interface{})
	var name, email string
	// to see if the user is already logged in, we don't ask for their information.
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser := getUserById(UserId)
		name = cUser.Name
		email = cUser.Email
	} else {
		name = r.FormValue("name")
		email = r.FormValue("email")
	}
	text := r.FormValue("text")
	newCnt := Cnt{Name: strings.TrimSpace(name), Email: strings.TrimSpace(email), Text: strings.TrimSpace(text)}
	if warnings := newCnt.Validate(); len(warnings) > 0 {
		sentFlash["Errs"] = warnings
		sentFlash["Name"] = name
		sentFlash["Email"] = email
		sentFlash["Text"] = text
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/contact", 302)
		return
	}
	if msgId := insertCnt(newCnt); msgId > 0 {
		sentFlash["Cong"] = "Congratulations dear " + name + ", your message has been sent successfully."
	} else {
		sentFlash["Err"] = "there is something wrong with your message, please check it and try again."
	}
	if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/search", 302)
}