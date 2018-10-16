package main

import (
	"fmt"
	_ "github.com/lib/pq"
)

func alreadyName(username string, userId int) int {
	rows, err := DB.Query("select id from users where name = $1", username)
	if err != nil {
		fmt.Println(err)
	}

	var id int
	for rows.Next() {
		rows.Scan(&id)
	}
	if id == 0 || userId == id {
		return 0
	}
	return 1
}

func alreadyEmail(email string, userId int) int64 {
	rows, err := DB.Query("select id from users where email = $1", email)
	if err != nil {
		fmt.Println(err)
	}

	var id int
	for rows.Next() {
		rows.Scan(&id)
	}
	if id == 0 || userId == id {
		return 0
	}
	return 1
}

func checkUsername(uORe string) (ItsId int) {
	rows, err := DB.Query("select id from users where name = $1 or email = $1", uORe)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&ItsId)
		if err != nil {
			fmt.Println(err)
		}
	}
	return ItsId
}
