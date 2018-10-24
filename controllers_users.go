package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func UserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	// the user that is logged in
	var cUser User
	// the user that its details is passed through url
	var psUser User
	parameter := ps.ByName("id")
	if userId, err := strconv.Atoi(parameter); err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}

	if userId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if userId > 0 {
		cUser = getUserById(userId)
	}
	gotFlash := make(map[string]interface{})

	gotFlash["Username"] = cUser.Name
	gotFlash["Notis"] = getUserNotifications(cUser.Id)
	gotFlash["User"] = psUser
	gotFlash["Title"] = psUser.Name
	posts := getUserPosts(psUser.Name)
	for i := range posts {
		// we are limiting the length of its text and tags
		if len(posts[i].Text) > 100 {
			posts[i].Text = posts[i].Text[:100] + " ..."
		}
		if len(posts[i].Tags) > 3 {
			posts[i].Tags = posts[i].Tags[:3]
		}
	}
	gotFlash["Posts"] = posts
	gotFlash["CUserFollowState"], gotFlash["PSUserFollowState"] = followState(cUser.Id, psUser.Id)
	if err := tpl.ExecuteTemplate(w, "user.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func UserProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	// the user that is logged in
	var cUser User
	// the user that its details is passed through url
	var psUser User
	parameter := ps.ByName("id")
	if userId, err := strconv.Atoi(parameter); err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}
	if userId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if userId > 0 {
		cUser = getUserById(userId)
	}else{
		w.Write([]byte("Please First Register"))
		return
	}

	if cUser.Name == psUser.Name{
		w.Write([]byte("Are You Serious?!"))
		return
	}
	submit := r.FormValue("submit")
	if submit == "follow"{

		// first we send a notification for the guy who has been followed.
		if affect := insertNotification(psUser.Id, cUser.Name,0,"follow", r.FormValue("time")); affect < 1{
			fmt.Println("inserting follow rel didn't happened")
		}
		// then insert a rel for this request, if the followed guy's profile is private
		// the 'accepted' row of followRel table should be false, until the followed guy
		// changes it, but if the followed guy's profile is public, we fill the accepted
		// row of followRel with true keyword.
		if affect := insertFollowRel(cUser.Id, psUser.Id, !psUser.Private); affect < 1{
			fmt.Println("inserting follow rel didn't happened")
		}
		w.Write([]byte("done"))
		return
	}

	if submit == "unfollow"{
		// first we send a notification for the guy who has been followed.
		if affect := insertNotification(psUser.Id, cUser.Name, 0,"unfollow", r.FormValue("time")); affect < 1{
			fmt.Println("inserting follow rel didn't happened")
		}
		if affect := deleteFollowRel(cUser.Id, psUser.Id); affect < 1{
			fmt.Println("deleting follow rel didn't happened")
		}
		w.Write([]byte("done"))
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	// the user that is logged in
	var cUser User
	// the user that its details is passed through url
	var psUser User
	parameter := ps.ByName("id")
	if userId, err := strconv.Atoi(parameter); err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}

	if userId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if userId > 0 {
		cUser = getUserById(userId)
	} else if userId < 1 || cUser.Id != psUser.Id {
		http.Redirect(w, r, "/", 303)
		return
	}

	gotFlash := make(map[string]interface{})

	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
		fmt.Println(err)
	}

	if cong, err := session.PopString(w, "Cong"); err != nil {
		fmt.Println(err)
	} else {
		gotFlash["Cong"] = cong
	}

	if cong, err := session.PopString(w, "Err"); err != nil {
		fmt.Println(err)
	} else {
		gotFlash["Err"] = cong
	}

	posts := getUserPosts(cUser.Name)
	for i := range posts {
		// we are limiting the length of its text and tags
		if len(posts[i].Text) > 100 {
			posts[i].Text = posts[i].Text[:100] + " ..."
		}
		if len(posts[i].Tags) > 3 {
			posts[i].Tags = posts[i].Tags[:3]
		}
	}
	gotFlash["Posts"] = posts
	gotFlash["Notis"] = getUserNotifications(cUser.Id)
	gotFlash["User"], gotFlash["Title"], gotFlash["Username"] = cUser, cUser.Name, cUser.Name
	if err := tpl.ExecuteTemplate(w, "profile.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func ProfileProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	var psUser User
	parameter := ps.ByName("id")
	if UserId, err := strconv.Atoi(parameter); err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(UserId)
	}
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	} else if UserId < 1 || cUser.Id != psUser.Id {
		http.Redirect(w, r, "/", 303)
		return
	}
	submit := r.FormValue("submit")
	username := r.FormValue("username")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	bio := r.FormValue("bio")
	pri := r.FormValue("pri")

	var private bool
	if pri == "yes" {
		private = true
	}
	fmt.Println("submit: ", submit)
	if submit == "Deactivate" {
		session.Destroy(w)
		// first we remove the profile pic

		if affect := deactiveUserById(cUser.Id); affect < 1 {
			fmt.Println("Your Account Can't Be Deactivated.")
		}
		// if successful
		if err := session.PutString(w, "Cong", "Your Account Has Been Deactivated."); err != nil {
			fmt.Println(err)
		}
		w.Write([]byte("done"))
		return
	}
	if submit == "DelPP" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		UserPic := filepath.Join(wd, "static", "pic", "pros", cUser.Pic)
		if err := os.Remove(UserPic); err != nil {
			fmt.Println(err)
		}
		if affect := deleteUserPic(cUser.Id); affect < 1 {
			w.Write([]byte("Your Picture Can't Be Deleted."))
			return
		}
		// if successful
		w.Write([]byte("done"))
		return
	}
	if submit == "UpPass" {
		cPass := r.FormValue("cPass")
		newPass := r.FormValue("newPass")
		confirmPass := r.FormValue("confirmPass")
		hashed := cUser.Password
		// check to see if the password is correct
		if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(cPass)); err != nil {
			w.Write([]byte("Your Current Password Is Wrong."))
			return
		}

		if newPass != confirmPass {
			w.Write([]byte("Your New Password Doesn't Match With The Confirm Password."))
			return
		}
		if len(newPass) < 8 || len(newPass) > 50 {
			w.Write([]byte("The Length Of Your Password Should Be Between 8 And 50 Characters."))
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		if affect := updatePass(cUser.Id, string(hash)); affect < 1 {
			w.Write([]byte("Your New Password Didn't Change, Please Try Again."))
			return
		}

		w.Write([]byte("done"))
		return
	}
	if submit == "UpPP" {
		var maxFileSize int64 = 4 * 1000 * 1000 //limit upload file to 10m
		if r.ContentLength > maxFileSize {
			if err := session.PutString(w, "Err", "The Maximum Size Of Profile Picture Is 4M."); err != nil {
				fmt.Println(err)
			}
			http.Redirect(w, r, "/profile/"+cUser.Name, 302)
			return
		}

		file, FileHeader, err := r.FormFile("pic")
		if err != nil {
			fmt.Println(err)
		}

		// we examine the existence of the pic in two ways, first we see if
		// err is not equal to missingFile error, then get the length of the
		// filename, if both are ok
		if err != http.ErrMissingFile && len(FileHeader.Filename) != 0 {
			// we validate the file, if it's not a picture we return an err
			if valid := detectFileType(file); !valid {
				if err := session.PutString(w, "Err", "Please Upload An Image, Other Types Are Not Supported."); err != nil {
					fmt.Println(err)
				}
				http.Redirect(w, r, "/profile/"+cUser.Name, 302)
				return
			}
			// if validation is passed, we delete the previous pic
			wd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			}
			UserPic := filepath.Join(wd, "static", "pic", "pros", cUser.Pic)
			if cUser.Pic != ""{
				if err := os.Remove(UserPic); err != nil{
					fmt.Println(err)
				}
			}

			// now put the new one
			picName := processPic(file, cUser.Name, "pros")
			defer file.Close()

			// and then pass the picName to database
			if affect := updateProPic(picName, cUser.Id); affect < 1 {
				fmt.Println("updating profile pic in db didn't happened")
			}
		}
		if err := session.PutString(w, "Cong", "Your Photo Has Been Updated."); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/profile/"+cUser.Name, 302)
		return
	}

	upUser := User{Id: cUser.Id, Name: strings.TrimSpace(username), FirstName: strings.TrimSpace(firstName), LastName: strings.TrimSpace(lastName), Email: strings.TrimSpace(email), Phone: strings.TrimSpace(phone), Bio: strings.TrimSpace(bio), Private: private}
	//if validation of form failed
	if warnings := upUser.UpValidate(); len(warnings) > 0 {
		w.Write([]byte(warnings))
		return
	}

	// if validation was passed
	if affect := updateUser(upUser); affect < 1 {
		w.Write([]byte("Your Account Didn't get Updated."))
		return
	}
	if err := session.PutString(w, "Cong", "Your Profile has been updated"); err != nil {
		fmt.Println(err)
	}

	w.Write([]byte("done"))
}

func MakePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	var psUser User
	parameter := ps.ByName("id")
	userId, err := strconv.Atoi(parameter)
	if err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	} else if UserId < 1 || UserId != psUser.Id {
		http.Redirect(w, r, "/", 303)
		return
	}
	gotFlash := make(map[string]interface{})
	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
		fmt.Println(err)
	}
	gotFlash["Notis"] = getUserNotifications(cUser.Id)
	gotFlash["Title"], gotFlash["Username"] = "Make Post", cUser.Name

	if err := tpl.ExecuteTemplate(w, "makePost.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func MakePostProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	var psUser User
	parameter := ps.ByName("id")
	userId, err := strconv.Atoi(parameter)
	if err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	} else if UserId < 1 || UserId != psUser.Id {
		http.Redirect(w, r, "/", 303)
		return
	}
	sentFlash := make(map[string]interface{})
	text := r.FormValue("text")
	tags := tagFinder(text)
	var maxFileSize int64 = 15 * 1000 * 1000 //limit upload file to 15m
	if r.ContentLength > maxFileSize {
		sentFlash["PostErr"] = "The Maximum Size Of Post Picture Is 15M."
		sentFlash["Text"] = text
		sentFlash["ProField"] = "posts"
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/profile/"+cUser.Name+"/post", 303)
		return
	}

	var picNames []string
	for i := 1; i < 7; i++ {
		// cause we don't know if the file exist or no, ignore the error
		file, FileHeader, err := r.FormFile("pic-" + strconv.Itoa(i))
		// we examine the existence of the pic in two ways, first we see if
		// err is not equal to missingFile error, then get the length of the
		// filename, if both are ok
		if err != http.ErrMissingFile && len(FileHeader.Filename) != 0 {
			defer file.Close()
			if valid := detectFileType(file); !valid {
				sentFlash["PostErr"] = "Please Upload An Image, Other Types Are Not Supported."
				sentFlash["Text"] = text
				sentFlash["ProField"] = "posts"
				if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
					fmt.Println(err)
				}
				http.Redirect(w, r, "/profile/"+cUser.Name+"/post", 303)
				return
			}
			picNames = append(picNames, processPic(file, cUser.Name, "posts"))
		}
	}

	if len(picNames) == 0 {
		sentFlash["PostErr"] = "At least add one image to your post"
		sentFlash["Text"] = text
		sentFlash["ProField"] = "posts"
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/profile/"+cUser.Name+"/post", 303)
		return
	}

	newPost := Post{Text: text, By: cUser.Name}

	var postId int
	if postId = insertPost(newPost, picNames); postId < 1 {
		fmt.Println("insertPost failed.")
		return
	}

	processTag(postId, tags)
	http.Redirect(w, r, "/profile/"+cUser.Name+"/post", 302)
}

//func MakeMessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	session := sessionManager.Load(r)
//	var cUser User
//	var psUser User
//	parameter := ps.ByName("id")
//	userId, err := strconv.Atoi(parameter)
//	if err != nil {
//		psUser = getUserByName(parameter)
//	} else {
//		psUser = getUserById(userId)
//	}
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	} else if UserId < 1 || UserId != psUser.Id {
//		http.Redirect(w, r, "/", 303)
//		return
//	}
//	gotFlash := make(map[string]interface{})
//	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
//		fmt.Println(err)
//	}
//
//	gotFlash["Title"], gotFlash["Username"], gotFlash["Msgs"] = "Messages", cUser.Name, getUserMessages(cUser.Id)
//
//	if err := tpl.ExecuteTemplate(w, "makePost.gohtml", gotFlash); err != nil {
//		fmt.Println(err)
//	}
//}

func MessagesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	var psUser User
	parameter := ps.ByName("id")
	userId, err := strconv.Atoi(parameter)
	if err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	} else if UserId < 1 || UserId != psUser.Id {
		http.Redirect(w, r, "/", 303)
		return
	}
	gotFlash := make(map[string]interface{})
	gotFlash["Notis"] = getUserNotifications(cUser.Id)
	gotFlash["Title"], gotFlash["MsgHeaders"], gotFlash["Username"] = "Messages", getMsgHeaders(cUser.Name), cUser.Name
	if err := tpl.ExecuteTemplate(w, "messages.gohtml", gotFlash); err != nil{
		fmt.Println(err)
	}
}

//func MessagesProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	session := sessionManager.Load(r)
//	var cUser User
//	var psUser User
//	parameter := ps.ByName("id")
//	userId, err := strconv.Atoi(parameter)
//	if err != nil {
//		psUser = getUserByName(parameter)
//	} else {
//		psUser = getUserById(userId)
//	}
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	} else if UserId < 1 || UserId != psUser.Id {
//		http.Redirect(w, r, "/", 303)
//		return
//	}
//}

func MessageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	var psUser User
	parameter := ps.ByName("id")
	parameter2 := ps.ByName("to")
	userId, err := strconv.Atoi(parameter)
	if err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}
	userId2, err2 := strconv.Atoi(parameter2)
	var user2 User
	if err2 != nil {
		user2 = getUserByName(parameter2)
	} else {
		user2 = getUserById(userId2)
	}
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	} else if UserId < 1 || UserId != psUser.Id {
		http.Redirect(w, r, "/", 303)
		return
	}
	if user2.Id < 1{
		http.Redirect(w, r, "/profile/"+cUser.Name+"/messages", 303)
		return
	}
	gotFlash := make(map[string]interface{})
	gotFlash["Notis"] = getUserNotifications(cUser.Id)
	gotFlash["Msgs"], gotFlash["Username"], gotFlash["Title"] = getMsgs(cUser.Name, user2.Name), cUser.Name, "Direct"
	if err := tpl.ExecuteTemplate(w, "message.gohtml", gotFlash); err != nil{
		fmt.Println(err)
	}
}


//func MessageProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	session := sessionManager.Load(r)
//	var cUser User
//	var psUser User
//	parameter := ps.ByName("id")
//	userId, err := strconv.Atoi(parameter)
//	if err != nil {
//		psUser = getUserByName(parameter)
//	} else {
//		psUser = getUserById(userId)
//	}
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	} else if UserId < 1 || UserId != psUser.Id {
//		http.Redirect(w, r, "/", 303)
//		return
//	}
//}


func FollowingHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	var psUser User
	parameter := ps.ByName("id")
	userId, err := strconv.Atoi(parameter)
	if err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	} else if UserId < 1 || UserId != psUser.Id {
		http.Redirect(w, r, "/", 303)
		return
	}
	gotFlash := make(map[string]interface{})
	gotFlash["Notis"] = getUserNotifications(cUser.Id)
	gotFlash["Users"], gotFlash["Title"], gotFlash["Username"] = getUserFollowings(cUser.Id), "Following", cUser.Name
	if err := tpl.ExecuteTemplate(w, "following.gohtml", gotFlash); err != nil{
		fmt.Println(err)
	}
}
