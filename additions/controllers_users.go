package additions

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func UserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else if UserId > 0 {
		cUser = getUserById(UserId)
	}
	parameter := ps.ByName("id")
	var user User
	if userId, err := strconv.Atoi(parameter); err != nil {
		user = getUserByName(parameter)
	} else {
		user = getUserById(userId)
	}
	gotFlash := make(map[string]interface{})

	gotFlash["UserName"] = cUser.Name
	gotFlash["User"] = user
	gotFlash["Title"] = user.Name
	gotFlash["Posts"] = getUserPosts(user.Id)
	if err := tpl.ExecuteTemplate(w, "user.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func Profile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	// the user that is logged in
	var cUser User
	// the user that its details is passed through url
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

	if cong, err := session.PopString(w, "Cong"); err != nil {
		fmt.Println(err)
	} else {
		gotFlash["Cong"] = cong
	}
	gotFlash["Posts"] = getUserPosts(cUser.Id)
	gotFlash["User"], gotFlash["Title"], gotFlash["UserName"] = cUser, cUser.Name, cUser.Name

	if err := tpl.ExecuteTemplate(w, "profile.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func EditProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	gotFlash["User"], gotFlash["Title"], gotFlash["UserName"] = cUser, cUser.Name, cUser.Name

	if err := tpl.ExecuteTemplate(w, "editProfile.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}

func EditProfileProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	} else if UserId < 1 || cUser != psUser {
		http.Redirect(w, r, "/", 303)
		return
	}
	sentFlash := make(map[string]interface{})
	submit := r.FormValue("submit")
	username := r.FormValue("username")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	quote := r.FormValue("quote")
	per := r.FormValue("private")

	var private bool
	if per == "true" {
		private = true
	}
	if submit == "Deactive" {
		session.Destroy(w)
		// first we remove the profile pic
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		if affect := deactiveUserById(cUser.Id); affect < 1 {
			fmt.Println("Your Account Can't Be Deleted.")
		}
		// if successful
		if err := session.PutString(w, "Cong", "Your Account Has Been Deactivated."); err != nil {
			fmt.Println(err)
		}
		w.Write([]byte("done"))
		return
	}
	fmt.Println("submit: ", submit)
	//if submit == "DeleteImg" {
	//	fmt.Println("got to delete img")
	//	wd, err := os.Getwd()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	UserPic := filepath.Join(wd, "static", "pic", "pros", strconv.Itoa(cUser.Id), cUser.Pic)
	//	if err := os.Remove(UserPic); err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println("userid: ", cUser.Id)
	//	if affect := deleteUsersPic(cUser.Id); affect < 1 {
	//		fmt.Println("Your Picture Can't Be Deleted.")
	//	}
	//	// if successful
	//	w.Write([]byte("done"))
	//	return
	//}

	if submit == "Update Password"{
		cPass := r.FormValue("cPass")
		newPass := r.FormValue("newPass")
		confirmPass := r.FormValue("confirmPass")
		hashed := cUser.Password
		// check to see if the password is correct
		if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(cPass)); err != nil {
			sentFlash["Err"] = "Your Current Password Is Wrong."
			if err := session.PutObject(w, "sentFlash",sentFlash); err != nil {
				fmt.Println(err)
			}
			http.Redirect(w, r, "/profile/"+cUser.Name+"/edit", 303)
			return
		}

		if newPass != confirmPass {
			sentFlash["Err"] = "Your New Password Doesn't Match With The Confirm Password."
			if err := session.PutObject(w, "sentFlash",sentFlash); err != nil {
				fmt.Println(err)
			}
			http.Redirect(w, r, "/profile/"+cUser.Name+"/edit", 303)
			return
		}
		if len(newPass) < 8 || len(newPass) > 50 {
			sentFlash["Err"] = "The Length Of Your Password Should Be Between 8 And 50 Characters."
			if err := session.PutObject(w, "sentFlash",sentFlash); err != nil {
				fmt.Println(err)
			}
			http.Redirect(w, r, "/profile/"+cUser.Name+"/edit", 303)
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		if affect := updatePass(cUser.Id, string(hash)); affect < 1 {
			sentFlash["Err"] = "Your New Password Didn't Change, Please Try Again."
			if err := session.PutObject(w, "sentFlash",sentFlash); err != nil {
				fmt.Println(err)
			}
			http.Redirect(w, r, "/profile/"+cUser.Name+"/edit", 303)
			return
		}

		if err := session.PutString(w, "Cong", "Your Password Has Been Updated."); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/profile/"+cUser.Name, 303)
		return
	}

	//var maxFileSize int64 = 4 * 1000 * 1000 //limit upload file to 10m
	//if r.ContentLength > maxFileSize {
	//	sentFlash["Err"] = "The Maximum Size Of Profile Picture Is 4M."
	//	if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
	//		fmt.Println(err)
	//	}
	//	http.Redirect(w, r, "/profile/"+strconv.Itoa(cUser.Id)+"/edit", 303)
	//	return
	//}
	//
	//file, FileHeader, err := r.FormFile("pic")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//var picName string
	//if err != http.ErrMissingFile && len(FileHeader.Filename) != 0{
	//	if valid := detectFileType(file); !valid {
	//		sentFlash["Err"] = "Please Upload An Image, Other Types Are Not Supported."
	//		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
	//			fmt.Println(err)
	//		}
	//		http.Redirect(w, r, "/profile/"+strconv.Itoa(cUser.Id)+"/edit", 303)
	//		return
	//	}
	//	picName = processProPic(file, cUser)
	//	defer file.Close()
	//}

	upUser := User{Id: cUser.Id, Name: strings.TrimSpace(username), FirstName: strings.TrimSpace(firstName), LastName: strings.TrimSpace(lastName), Email: strings.TrimSpace(email), Quote: strings.TrimSpace(quote), Private: private}
	//if picName != "" {
	//	upUser.Pic = picName
	//} else {
	//	upUser.Pic = cUser.Pic
	//}
	// if validation of form failed
	// we give it just a password, to avoid validation from returning err for password being empty
	if warnings := upUser.UpValidate(); len(warnings) > 0 {
		sentFlash["Private"] = private
		sentFlash["Quote"] = quote
		sentFlash["Errs"] = warnings
		sentFlash["UserName"] = username
		sentFlash["FirstName"] = firstName
		sentFlash["LastName"] = lastName
		sentFlash["Email"] = email
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/profile/"+cUser.Name+"/edit", 303)
		return
	}

	// if validation was passed
	if affect := updateUser(upUser); affect < 1 {
		sentFlash["Errs"] = "Your Account Didn't Updated."
		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/profile/"+cUser.Name+"/edit", 303)
		return
	}
	if err := session.PutString(w, "Cong", "Your Profile Has Been Updated."); err != nil {
		fmt.Println(err)
	}
	// in case user has changed username, we don't use the previous one.
	http.Redirect(w, r, "/profile/"+username, 302)
}
