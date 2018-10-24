package main

import (
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	Id                                                  											 int
	IsOnline, Private											 									 bool
	Name, FirstName, LastName, Email, Phone, Password, Pic, Bio, LastLog, CreatedOn, DeactivatedOn string
}
type Noti struct {
	UserId, RelatedPostId 		int
	RelatedUsername, Condition, CreatedOn	string
}

func searchUsers(word string) (users []User) {
	word = "%" + word + "%"

	rows, err := DB.Query("SELECT * FROM users WHERE name LIKE $1", word)
	if err != nil {
		fmt.Println(err)

	}
	var user User
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Bio, &user.IsOnline, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
		users = append(users, user)
	}
	return users
}

func searchTags(word string) (posts []Post) {
	word = "%" + word + "%"

	rows, err := DB.Query("SELECT id FROM tags WHERE name LIKE $1", word)
	if err != nil {
		fmt.Println(err)

	}
	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}

	for _, id := range ids{
		rows, err = DB.Query("SELECT * FROM posts WHERE id = $1", id)
		if err != nil {
			fmt.Println(err)

		}
		var post Post
		for rows.Next() {
			rows.Scan(&post.Id, &post.Text, &post.By, &post.ViewCount, &post.Like, &post.CreatedOn, &post.UpdatedOn, &post.DeletedOn)
			pics := getPostPics(post.Id)
			post.Pics = pics
			posts = append(posts, post)
		}
	}

	return posts
}

func getUserByName(Name string) (user User) {
	rows, err := DB.Query("select * from users where name = $1", Name)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Bio, &user.IsOnline, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
	}
	return user
}

func getUsersByIds(ids []int) (users []User) {
	var user User
	for _, id := range ids {
		rows, err := DB.Query("select * from users where id = $1", id)
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Bio, &user.IsOnline, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
		}
		users = append(users, user)
	}
	return users
}

func insertUser(user User) (ItsId int) {
	err := DB.QueryRow(
		"INSERT INTO users(name,firstName,lastName,email,password,isOnline,createdOn) VALUES($1,$2,$3,$4,$5,true,$6) returning id;",
		user.Name, user.FirstName, user.LastName, user.Email, user.Password, getNow()).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}
	return ItsId
}

func insertFollowRel(followerId, followedId int, accepted bool) (ItsId int) {
	// if there is already a relation is our database with
	// these users, we don't duplicate it and just return 1
	err := DB.QueryRow(
		"INSERT INTO followRel VALUES($1,$2,$3) returning followerId;",
		followerId, followedId, accepted).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}
	return ItsId
}

func insertNotification(userId int, relatedUsername string, relatedPostId int, condition, time string) (ItsId int) {
	err := DB.QueryRow(
		"INSERT INTO notiRel VALUES($1,$2,$3,$4,$5) returning userId;",
		userId, relatedUsername, relatedPostId, condition, time).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}
	return ItsId
}

func getUserNotifications(userId int) (notis []Noti){
	rows, err := DB.Query("select * from notiRel where userId = $1 order by createdOn desc limit 5", userId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var noti Noti
		rows.Scan(&noti.UserId, &noti.RelatedUsername, &noti.RelatedPostId, &noti.Condition, &noti.CreatedOn)
		// we obtain the name of notification sender from the beginning of the text field

		notis = append(notis, noti)
	}
	return notis
}

func getAllUsers() (users []User) {
	rows, err := DB.Query("select * from users order by id desc")
	if err != nil {
		fmt.Println(err)
	}
	var user User
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Bio, &user.IsOnline, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
		users = append(users, user)
	}
	return users
}

func getUserFollowings(id int) (users []User){
	rows, err := DB.Query("select followedId from followRel where followerId = $1", id)
	if err != nil {
		fmt.Println(err)
	}
	var followed int
	for rows.Next() {
		rows.Scan(&followed)
		user := getUserById(followed)
		users = append(users, user)
	}
	return users
}

func deleteFollowRel(followerId, followedId int) int64 {
	stmt, err := DB.Prepare("delete from followRel where followerId = $1 and followedId = $2")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec(followerId, followedId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
}

func deactiveUserById(UserId int) int64 {
	stmt, err := DB.Prepare("update users set password= $1, deactivatedOn= $2 where id= $3")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec("", getNow(), UserId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
}

func logOutProcess(userId int, time int) int64 {
	stmt, err := DB.Prepare("update users set lastLog= $1, isOnline = false where id= $2")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec(time, userId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
}

func setUserOnline(userId int) int64 {
	stmt, err := DB.Prepare("update users set isOnline = true where id= $1")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec( userId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
}

func followState(cUserId, psUserId int) (cUserFollowState, psUserFollowState string) {
	rows, err := DB.Query("select followerId, accepted from followRel where followerId = $1 and followedId = $2", cUserId, psUserId)
	if err != nil {
		fmt.Println(err)
	}
	var id int
	var accepted bool
	for rows.Next() {
		err = rows.Scan(&id, &accepted)
		if err != nil {
			fmt.Println(err)
		}
	}
	if id > 0 && accepted == false{
		cUserFollowState = "sent"
	}else if id > 0 && accepted == true{
		cUserFollowState = "accepted"
	}else{
		cUserFollowState = "nil"
	}

	rows, err = DB.Query("select followerId, accepted from followRel where followerId = $1 and followedId = $2", psUserId, cUserId)
	if err != nil {
		fmt.Println(err)
	}
	var id2 int
	var accepted2 bool
	for rows.Next() {
		err = rows.Scan(&id2, &accepted2)
		if err != nil {
			fmt.Println(err)
		}
	}
	if id2 > 0 && accepted2 == false{
		psUserFollowState = "sent"
	}else if id2 > 0 && accepted2 == true{
		psUserFollowState = "accepted"
	}else{
		psUserFollowState = "nil"
	}
	return cUserFollowState, psUserFollowState
}

func deleteUserPic(userId int) int64 {
	stmt, err := DB.Prepare("update users set pic= $1 where id= $2")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec("", userId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
}

func updatePass(UserId int, password string) int64 {
	stmt, err := DB.Prepare("update users set password= $1 where id= $2")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec(password, UserId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
}

func updateUser(upUser User) int64 {
	// we should have previous information of the user before updating
	// to be able to find its dependent tables and update them with new one,
	// e.g. if we use upUser.Name to find its coms, it'll return nothing.
	// because there is no comment record with updated name
	notUpUser := getUserById(upUser.Id)
	stmt, err := DB.Prepare("update users set name= $1, firstName= $2, lastName= $3, email= $4, phone= $5, bio= $6, private= $7 where id= $8")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec(upUser.Name, upUser.FirstName, upUser.LastName, upUser.Email , upUser.Phone, upUser.Bio, upUser.Private, upUser.Id)
	if err != nil {
		fmt.Println(err)
	}
	UserAffect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	posts := getUserPosts(notUpUser.Name)
	for i, post := range posts {
		stmt, err = DB.Prepare("update posts set by= $1 where id= $2")
		if err != nil {
			fmt.Println(err)
		}
		res, err = stmt.Exec(upUser.Name, post.Id)
		if err != nil {
			fmt.Println(err)
		}
		PostAffect, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err)
		}
		if PostAffect < 1 {
			fmt.Println("PostAffect less than one at ", i)
		}
	}

	coms := getUserComs(notUpUser.Name)
	for i, notUpUser := range coms {
		stmt, err = DB.Prepare("update coms set by= $1 where id= $2")
		if err != nil {
			fmt.Println(err)
		}
		res, err = stmt.Exec(upUser.Name, notUpUser.Id)
		if err != nil {
			fmt.Println(err)
		}
		PostAffect, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err)
		}
		if PostAffect < 1 {
			fmt.Println("ComAffect less than one at ", i)
		}
	}
	return UserAffect
}

func updateProPic(pic string, userId int) int64{
	stmt, err := DB.Prepare("update users set pic= $1 where id= $2")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec(pic, userId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
}


func getUserPosts(username string) (posts []Post) {
	rows, err := DB.Query("select * from posts where by = $1", username)
	if err != nil {
		fmt.Println(err)
	}
	var post Post
	for rows.Next() {
		rows.Scan(&post.Id, &post.Text, &post.By, &post.ViewCount, &post.Like, &post.CreatedOn, &post.UpdatedOn, &post.DeletedOn)
		post.Pics = getPostPics(post.Id)
		post.Tags = getPostTags(post.Id)
		posts = append(posts, post)
	}
	return posts
}

func getUserById(UserId int) (user User) {
	rows, err := DB.Query(
		"SELECT * FROM users where id = $1", UserId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Bio, &user.IsOnline, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
	}
	return user
}