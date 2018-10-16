package main

import (
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	Id                                                  int
	IsActive, Private											 bool
	Name, FirstName, LastName, Email, Password, Quote, LastLog, CreatedOn, DeactivatedOn string
	Pics []string
}

func searchUsers(word string) (users []User) {
	word = "%" + word + "%"

	rows, err := DB.Query("SELECT * FROM users WHERE name LIKE $1", word)
	if err != nil {
		fmt.Println(err)

	}
	var user User
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
		users = append(users, user)
	}
	return users
}

func getUserByName(Name string) (user User) {

	rows, err := DB.Query("select * from users where name = $1", Name)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
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
			rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
		}
		users = append(users, user)
	}
	return users
}

func insertUser(user User) (ItsId int) {
	err := DB.QueryRow(
		"INSERT INTO users(name,firstName,lastName,email,password,createdOn) VALUES($1,$2,$3,$4,$5,$6) returning id;",
		user.Name, user.FirstName, user.LastLog, user.Email, user.Password, getNow()).Scan(&ItsId)
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
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
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

func deletePostPic(storyId int) int64{
	stmt, err := DB.Prepare("delete from postPicRel where postId= $1")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec(storyId)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	return affect
}

func deleteUserPic(UserId int) int64 {
	stmt, err := DB.Prepare("delete from userPicRel where userId= $1")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec(UserId)
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
//
//func updateUser(upUser User) int64 {
//	// we should have previous information of the user before updating
//	// to be able to find its dependent tables and update them with new one,
//	// e.g. if we use upUser.Name to find its coms, it'll return nothing.
//	// because there is no comment record with updated name
//	notUpUser := getUserById(upUser.Id)
//	stmt, err := DB.Prepare("update users set name= $1, firstName= $2, lastName= $3 email= $4, quote= $5, private= $6 where id= $7")
//	if err != nil {
//		fmt.Println(err)
//	}
//	res, err := stmt.Exec(upUser.Name, upUser.FirstName, upUser.LastName, upUser.Email, upUser.Quote, upUser.Private, upUser.Id)
//	if err != nil {
//		fmt.Println(err)
//	}
//	UserAffect, err := res.RowsAffected()
//	if err != nil {
//		fmt.Println(err)
//	}
//	posts := getUserPosts(notUpUser.Id)
//	for i, v := range posts {
//		stmt, err = DB.Prepare("update posts set by= $1 where id= $2")
//		if err != nil {
//			fmt.Println(err)
//		}
//		res, err = stmt.Exec(upUser.Name, v.Id)
//		if err != nil {
//			fmt.Println(err)
//		}
//		PostAffect, err := res.RowsAffected()
//		if err != nil {
//			fmt.Println(err)
//		}
//		if PostAffect < 1 {
//			fmt.Println("PostAffect less than one at ", i)
//		}
//	}
//
//	coms := getUsersComs(notUpUser.Name)
//	for i, notUpUser := range coms {
//		stmt, err = DB.Prepare("update coms set by= $1 where id= $2")
//		if err != nil {
//			fmt.Println(err)
//		}
//		res, err = stmt.Exec(upUser.Name, notUpUser.Id)
//		if err != nil {
//			fmt.Println(err)
//		}
//		PostAffect, err := res.RowsAffected()
//		if err != nil {
//			fmt.Println(err)
//		}
//		if PostAffect < 1 {
//			fmt.Println("ComAffect less than one at ", i)
//		}
//	}
//	return UserAffect
//}

func getUserById(UserId int) (user User) {
	rows, err := DB.Query(
		"SELECT * FROM users where id = $1", UserId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Quote, &user.IsActive, &user.Private, &user.LastLog, &user.CreatedOn, &user.DeactivatedOn)
	}
	return user
}