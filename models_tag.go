package main

import (
	"fmt"
	_ "github.com/lib/pq"
)

type Tag struct {
	Id              int
	Name, CreatedOn string
}
type TagRelations struct {
	TagId, PostId int
}

func getPostsByTagId(TagId int) (posts []Post) {
	rows, err := DB.Query("select postId from tagRel where tagId = $1", TagId)
	if err != nil {
		fmt.Println(err)
	}

	var ids []int

	for rows.Next() {
		var postId int
		rows.Scan(&postId)
		ids = append(ids, postId)
	}

	ids = getUniqueInt(ids)
	posts = getPostsByIds(ids)

	return posts
}

func getTagById(TagId int) (tag Tag) {
	rows, err := DB.Query("select * from tags where id = $1", TagId)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		rows.Scan(&tag.Id, &tag.Name, &tag.CreatedOn)
	}

	return tag
}

func getPostTags(postId int) (ItsTags []Tag) {
	rows, err := DB.Query("SELECT tagId FROM tagRel where postId = $1", postId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var tagId int
		err = rows.Scan(&tagId)
		if err != nil {
			fmt.Println(err)
		}

		// after we get a tagId, simultaneously we get the associated tag.
		rows, err := DB.Query("SELECT * FROM tags where id = $1", tagId)
		if err != nil {
			fmt.Println(err)
		}
		var tag Tag
		for rows.Next() {
			err = rows.Scan(&tag.Id, &tag.Name, &tag.CreatedOn)
			if err != nil {
				fmt.Println(err)
			}
			// and append it to the main variable
			ItsTags = append(ItsTags, tag)
		}
	}

	return ItsTags
}

func alreadyTag(tagName string) (ItsId int) {
	rows, err := DB.Query("select id from tags where name = $1", tagName)
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

func insertTag(tagName string) (ItsId int) {
	err := DB.QueryRow(
		"insert into tags(name,createdOn) VALUES($1,$2) returning id;",
		tagName, getNow()).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}

	return ItsId
}

func insertTagRel(tagId, postId int) (ItsId int) {
	err := DB.QueryRow(
		"insert into tagRel(tagId, postId) values($1, $2) returning tagId;",
		tagId, postId).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}

	return ItsId
}

func insertPicRel(postId int, pics []string) (ItsId int) {
	for _, pic := range pics{
		err := DB.QueryRow(
			"insert into picRel(pic, postId) values($1, $2) returning postId;",
			pic, postId).Scan(&ItsId)
		if err != nil {
			fmt.Println(err)
		}
		if ItsId < 1{
			fmt.Println("inseritng one of pic rels failed")
		}
	}


	return ItsId
}
