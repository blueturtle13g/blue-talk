package main

import (
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	Id                                                  											 int
	IsActive, Private											 									 bool
	Name, FirstName, LastName, Email, Phone, Password, Pic, Quote, LastLog, CreatedOn, DeactivatedOn string
}

func searchUsers(word string) (users []User) {
	word = "%" + word + "%"

	rows, err := DB.Query("SELECT * FROM users WHERE name LIKE $1", word)
	if err != nil {
		fmt.Println(err)

	}
	var user User
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
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
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
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
			rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
		}
		users = append(users, user)
	}
	return users
}

func insertUser(user User) (ItsId int) {
	err := DB.QueryRow(
		"INSERT INTO users(name,firstName,lastName,email,password,createdOn) VALUES($1,$2,$3,$4,$5,$6) returning id;",
		user.Name, user.FirstName, user.LastName, user.Email, user.Password, getNow()).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}
	return ItsId
}

func getAllUsers(order string) (users []User) {
	rows, err := DB.Query("select * from users order by id desc")
	if err != nil {
		fmt.Println(err)
	}
	var user User
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
		users = append(users, user)
	}
	return users
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

func updateLastLog(userId int) int64 {
	stmt, err := DB.Prepare("update users set lastLog= $1 where id= $2")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec(getNow(), userId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
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
	stmt, err := DB.Prepare("update users set name= $1, firstName= $2, lastName= $3, email= $4, phone= $5, quote= $6, private= $7 where id= $8")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec(upUser.Name, upUser.FirstName, upUser.LastName, upUser.Email , upUser.Phone, upUser.Quote, upUser.Private, upUser.Id)
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
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password, &user.Pic, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
	}
	return user
}