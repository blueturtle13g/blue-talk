package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"os"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"path/filepath"
)

//func UserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	session := sessionManager.Load(r)
//	var cUser User
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	}
//	parameter := ps.ByName("id")
//	var user User
//	if userId, err := strconv.Atoi(parameter); err != nil {
//		user = getUserByName(parameter)
//	} else {
//		user = getUserById(userId)
//	}
//	gotFlash := make(map[string]interface{})
//
//	gotFlash["UserName"] = cUser.Name
//	gotFlash["User"] = user
//	gotFlash["Title"] = user.Name
//	gotFlash["Posts"] = getUserPosts(user.Id)
//	if err := tpl.ExecuteTemplate(w, "user.gohtml", gotFlash); err != nil {
//		fmt.Println(err)
//	}
//}

func ProfileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	// the user that is logged in
	var cUser User
	// the user that its details is passed through url
	var psUser User
	parameter := ps.ByName("id")
	if userId, err := strconv.Atoi(parameter);err != nil {
		psUser = getUserByName(parameter)
	} else {
		psUser = getUserById(userId)
	}

	if userId, err := session.GetInt("UserId"); err != nil {
		http.Redirect(w, r, "/", 303)
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
	gotFlash["Posts"] = getUserPosts(cUser.Name)
	gotFlash["User"], gotFlash["Title"], gotFlash["Username"] = cUser, cUser.Name, cUser.Name
	if err := tpl.ExecuteTemplate(w, "profile.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

//func EditProfileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
//	gotFlash["User"], gotFlash["Title"], gotFlash["Username"] = cUser, cUser.Name, cUser.Name
//
//	if err := tpl.ExecuteTemplate(w, "editProfile.gohtml", gotFlash); err != nil {
//		fmt.Println(err)
//	}
//}
//func MakePostProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
//	sentFlash := make(map[string]interface{})
//
//	r.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
//	fhs := r.MultipartForm.File["myfiles"]
//	for _, fh := range fhs {
//		f, err := fh.Open()
//		// f is one of the files
//	}
//}

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
	quote := r.FormValue("quote")
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
		fmt.Println("got to delete img")
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
	if submit == "UpPass"{
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
	if submit == "UpPP"{
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
		if err != http.ErrMissingFile && len(FileHeader.Filename) != 0{
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
			if err := os.Remove(UserPic); err != nil {
				fmt.Println(err)
			}
			if affect := deleteUserPic(cUser.Id); affect < 1 {
				fmt.Println("Your Picture Can't Be Deleted.")
			}

			// now put the new one
			picName := processProPic(file, cUser)
			defer file.Close()

			// and then pass the picName to database
			if affect := updateProPic(picName, cUser.Id); affect < 1{
				if err := session.PutString(w, "Err", "It didn't get completed."); err != nil {
					fmt.Println(err)
				}
				http.Redirect(w, r, "/profile/"+cUser.Name, 302)
				return
			}
		}
		if err := session.PutString(w, "Cong", "Your Photo Has Been Updated."); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/profile/"+cUser.Name, 302)
		return
	}
	if submit == "post"{
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
			http.Redirect(w, r, "/profile/"+cUser.Name, 303)
			return
		}

		r.ParseMultipartForm(15 << 20)

		fhs := r.MultipartForm.File["pics"]
		var picNames []string
		for _, fh := range fhs {
			file, err := fh.Open()

			if err != http.ErrMissingFile && len(fh.Filename) != 0{
				if valid := detectFileType(file); !valid {
					sentFlash["PostErr"] = "Please Upload An Image, Other Types Are Not Supported."
					sentFlash["Text"] = text
					sentFlash["ProField"] = "posts"
					if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
						fmt.Println(err)
					}
					http.Redirect(w, r, "/profile/"+cUser.Name, 303)
					return
				}
				defer file.Close()
				picNames = append(picNames, processPostPic(file))
			}else{
				sentFlash["PostErr"] = "At least add one image to your post"
				sentFlash["Text"] = text
				sentFlash["ProField"] = "posts"
				if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
					fmt.Println(err)
				}
				http.Redirect(w, r, "/profile/"+cUser.Name, 303)
				return
			}
		}


		newPost := Post{Text: text, By: cUser.Name}

		var postId int
		if postId = insertPost(newPost); postId < 1 {
			fmt.Println("insertPost failed.")
			return
		}

		processTag(postId, tags)
		http.Redirect(w, r, "/profile/"+cUser.Name, 302)
		return
	}

	upUser := User{Id: cUser.Id, Name: strings.TrimSpace(username), FirstName: strings.TrimSpace(firstName), LastName: strings.TrimSpace(lastName), Email: strings.TrimSpace(email) , Phone: strings.TrimSpace(phone), Quote: strings.TrimSpace(quote), Private: private}
	//if validation of form failed
	if warnings := upUser.UpValidate(); len(warnings) > 0 {
		fmt.Println("warnings: ",warnings)
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
