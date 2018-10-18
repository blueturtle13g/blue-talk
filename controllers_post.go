package main
//
//import (
//	"net/http"
//	"strconv"
//
//	"fmt"
//	"github.com/julienschmidt/httprouter"
//	"os"
//	"path/filepath"
//)
//
//func PostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	session := sessionManager.Load(r)
//	var cUser User
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	}
//
//	parameter := ps.ByName("id")
//	postId, _ := strconv.Atoi(parameter)
//	post, ItsComments, ItsTags, ItsPics := getPostById(postId)
//	if post.By != cUser.Name{
//		if affect := incrementViewCount(postId); affect < 1 {
//			fmt.Println("view count didn't increment.")
//		}
//	}
//	gotFlash := make(map[string]interface{})
//	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
//		fmt.Println(err)
//	}
//	gotFlash["UserName"] = cUser.Name
//	gotFlash["PostBy"] = getUserByName(post.By)
//	gotFlash["Post"] = post
//	gotFlash["Coms"] = ItsComments
//	gotFlash["Tags"] = getUniqueTag(ItsTags)
//	if err := tpl.ExecuteTemplate(w, "post.gohtml", gotFlash); err != nil {
//		fmt.Println(err)
//	}
//}
//
//func PostProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	session := sessionManager.Load(r)
//	parameter := ps.ByName("id")
//	postId, _ := strconv.Atoi(parameter)
//	sentFlash := make(map[string]interface{})
//	// if user is not registered and is trying to comment,
//	// we ask them to first log in
//	var cUser User
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	} else {
//		w.Write([]byte("Please First Register"))
//		return
//	}
//	cPost, _, _,_ := getPostById(postId)
//	submit := r.FormValue("submit")
//	text := r.FormValue("text")
//
//	// if the request was commenting and not rating.
//	if submit == "com" {
//		if len(text) < 3 || len(text) > 500 {
//			w.Write([]byte("The length of your comment can't be less than 3 or more than 500 characters."))
//			return
//		}
//		newCom := Com{By: cUser.Name, PostId: postId, Text: text}
//		if affect := insertCom(newCom); affect < 0 {
//			fmt.Println("inserting comment failed.")
//		}
//		w.Write([]byte("done"))
//		return
//	}
//	if submit == "resCom" {
//		comId, err := strconv.Atoi(r.FormValue("comId"))
//		if err != nil {
//			fmt.Println(err)
//		}
//		if len(text) < 3 || len(text) > 500 {
//			w.Write([]byte("The length of your comment can't be less than 3 or more than 500 characters"))
//			return
//		}
//		newComReply := Com{By: cUser.Name, ComId: comId, PostId: cPost.Id, Text: text}
//		if resComId := insertCom(newComReply); resComId < 1 {
//			fmt.Println("inserting resCom failed.")
//		}
//		w.Write([]byte("done"))
//		return
//	}
//	if submit == "EditCom" {
//		comId, err := strconv.Atoi(r.FormValue("comId"))
//		if err != nil {
//			fmt.Println(err)
//		}
//		if len(text) < 3 || len(text) > 500 {
//			w.Write([]byte("The length of your comment can't be less than 3 or more than 500 characters"))
//			return
//		}
//		// if we encounter an internal problem
//		if affect := updateCom(comId, text); affect < 1 {
//			fmt.Println("Updating Failed, Please Try Again.")
//			return
//		}
//		w.Write([]byte("done"))
//		return
//	}
//	if submit == "DeleteCom" {
//		comId, err := strconv.Atoi(r.FormValue("comId"))
//		if err != nil {
//			fmt.Println(err)
//		}
//		if affect := deleteComment(comId); affect < 1 {
//			w.Write([]byte("The Deletion Of Your Comment Failed, Please Try Again"))
//			return
//		}
//		w.Write([]byte("done"))
//		return
//	}
//}
//
//func EditPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	session := sessionManager.Load(r)
//	parameter := ps.ByName("id")
//	postId, err := strconv.Atoi(parameter)
//	if err != nil {
//		fmt.Println("Please pass an id")
//	}
//	psPost, _, _, _ := getPostById(postId)
//
//	var cUser User
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	} else if UserId < 1 || cUser.Name != psPost.By {
//		http.Redirect(w, r, "/", 303)
//		return
//	}
//	gotFlash := make(map[string]interface{})
//	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
//		fmt.Println(err)
//	}
//
//	gotFlash["Post"], gotFlash["UserName"] = psPost, cUser.Name
//
//	if err := tpl.ExecuteTemplate(w, "editPost.gohtml", gotFlash); err != nil {
//		fmt.Println(err)
//	}
//}
//
//func EditPostProcess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	session := sessionManager.Load(r)
//	parameter := ps.ByName("id")
//	postId, err := strconv.Atoi(parameter)
//	if err != nil {
//		fmt.Println("Please pass an id")
//	}
//	psPost, _, _, _ := getPostById(postId)
//
//	var cUser User
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	} else if UserId < 1 || cUser.Name != psPost.By {
//		http.Redirect(w, r, "/", 303)
//		return
//	}
//	sentFlash := make(map[string]interface{})
//	//if r.FormValue("submit") == "Delete" {
//	//	// first thing we should take care of, is deleting photo
//	//	wd, err := os.Getwd()
//	//	if err != nil {
//	//		fmt.Println(err)
//	//	}
//	//	PostPic := filepath.Join(wd, "static", "pic", "stories", strconv.Itoa(psPost.Id))
//	//	if err := os.Remove(PostPic); err != nil {
//	//		fmt.Println(err)
//	//	}
//	//	if affect := deletePost(psPost.Id); affect < 0 {
//	//		fmt.Println("deletion of post didn't happen")
//	//	}
//	//	if err := session.PutString(w, "Cong", "Your Post Is Deleted."); err != nil {
//	//		fmt.Println(err)
//	//	}
//	//	w.Write([]byte("done"))
//	//	return
//	//}
//	fmt.Println("submit: ", r.FormValue("submit"))
//	//if r.FormValue("submit") == "DeleteImg" {
//	//	wd, err := os.Getwd()
//	//	if err != nil {
//	//		fmt.Println(err)
//	//	}
//	//	PostPic := filepath.Join(wd, "static", "pic", "stories", strconv.Itoa(psPost.Id), psPost.Pic)
//	//	if err := os.Remove(PostPic); err != nil {
//	//		fmt.Println(err)
//	//	}
//	//	if affect := deletePostPic(psPost.Id); affect < 1 {
//	//		fmt.Println("Your Image Can't Be Deleted.")
//	//	}
//	//	w.Write([]byte("done"))
//	//	return
//	//}
//	var maxFileSize int64 = 4 * 1000 * 1000 //limit upload file to 10m
//	if r.ContentLength > maxFileSize {
//		sentFlash["Err"] = "The Maximum Size Of Post Picture Is 4M."
//		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
//			fmt.Println(err)
//		}
//		http.Redirect(w, r, "/post/"+strconv.Itoa(psPost.Id)+"/edit", 303)
//		return
//	}
//
//	file, FileHeader, err := r.FormFile("pic")
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer file.Close()
//	var picName string
//	if err != http.ErrMissingFile && len(FileHeader.Filename) != 0{
//		if valid := detectFileType(file); !valid {
//			sentFlash["Err"] = "Please Upload An Image, Other Types Are Not Supported."
//			if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
//				fmt.Println(err)
//			}
//			http.Redirect(w, r, "/post/"+strconv.Itoa(psPost.Id)+"/edit", 303)
//			return
//		}
//		picName = processPostPic(file)
//		defer file.Close()
//	}
//	text := r.FormValue("text")
//	tags := tagFinder(text)
//	selectedCat := r.FormValue("selectCat")
//	madeCat := r.FormValue("cat")
//
//	var warnings []string
//
//	upPost := Post{Id: postId, Text: text, By: cUser.Name}
//	sentFlash["TitlePost"] = title
//	sentFlash["Text"] = text
//	sentFlash["SelectedCat"] = selectedCat
//	sentFlash["Cat"] = madeCat
//	sentFlash["Cats"] = getAllCats("", "")
//	if warnings, cat = upPost.Validate(selectedCat, madeCat, "up", tags); len(warnings) > 0 {
//		sentFlash["Errs"] = warnings
//		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
//			fmt.Println(err)
//		}
//		http.Redirect(w, r, "/post/"+strconv.Itoa(psPost.Id)+"/edit", 303)
//		return
//	}
//
//	if affect := updatePost(upPost); affect < 1 {
//		sentFlash["Err"] = "There is something wrong with your data, Please check it and try again."
//		if err := session.PutObject(w, "sentFlash", sentFlash); err != nil {
//			fmt.Println(err)
//		}
//		http.Redirect(w, r, "/post/"+strconv.Itoa(psPost.Id)+"/edit", 303)
//		return
//	}
//	// if the process finished successfully.
//	processTag(postId, tags)
//	http.Redirect(w, r, "/post/"+strconv.Itoa(psPost.Id), 302)
//}
//
//func MakePostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	session := sessionManager.Load(r)
//	var cUser User
//	if UserId, err := session.GetInt("UserId"); err != nil {
//		fmt.Println(err)
//	} else if UserId > 0 {
//		cUser = getUserById(UserId)
//	} else {
//		http.Redirect(w, r, "/", 303)
//		return
//	}
//	gotFlash := make(map[string]interface{})
//	if err := session.PopObject(w, "sentFlash", &gotFlash); err != nil {
//		fmt.Println(err)
//	}
//	gotFlash["Title"], gotFlash["UserName"] = cUser.Name, cUser.Name
//
//	if err := tpl.ExecuteTemplate(w, "makePost.gohtml", gotFlash); err != nil {
//		fmt.Println(err)
//	}
//}
//
