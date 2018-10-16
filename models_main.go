package main

import (
	"fmt"
	_ "github.com/lib/pq"
)

type Cnt struct {
	Id                           int
	Name, Email, Text, CreatedOn string
}

func insertCnt(cnt Cnt) (ItsId int) {
	err := DB.QueryRow(
		"INSERT INTO cnt(name,email,text,createdOn) VALUES($1,$2,$3,$4) returning id;",
		cnt.Name, cnt.Email, cnt.Text, getNow()).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}
	return ItsId
}
