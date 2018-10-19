package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"database/sql"
)

type Post struct {
	Id, ViewCount , Like                            int
	Text, By, CreatedOn, UpdatedOn, DeletedOn       string
	Pics []string
	Tags []Tag
}


func updatePost(upPost Post) int64 {
	// first we delete previous relations to put new ones.
	stmt, err := DB.Prepare("delete from tagrel where postId= $1")
	if err != nil {
		fmt.Println(err)
	}

	res, err := stmt.Exec(upPost.Id)
	if err != nil {
		fmt.Println(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	if affect < 1 {
		fmt.Println("tagrelation didn't get cleared.")
	}
	
	stmt, err = DB.Prepare("update posts set text= $1, updatedOn= $2 where id= $3")
	if err != nil {
		fmt.Println(err)
	}
	res, err = stmt.Exec(upPost.Text, getNow(), upPost.Id)
	if err != nil {
		fmt.Println(err)
	}

	affect, err = res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	return affect
}

func getPostsByIds(ids []int) (posts []Post) {
	var post Post
	for _, id := range ids {
		rows, err := DB.Query("select * from posts where id = $1", id)
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			rows.Scan(&post.Id, &post.Text, &post.By, &post.ViewCount, &post.Like, &post.CreatedOn, &post.UpdatedOn, &post.DeletedOn)
		}
		post.Pics = getPostPics(post.Id)
		post.Tags = getPostTags(post.Id)
		posts = append(posts, post)
	}
	return posts
}

func getAllPosts(order string) (posts []Post) {
	var rows *sql.Rows
	var err error
	if order == "like" {
		rows, err = DB.Query("select * from posts order by Liked desc")
	}else{
		rows, err = DB.Query("select * from posts order by id desc")
	
	}
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


func getPostPics(postId int) (pics []string) {
	rows, err := DB.Query("select pic from picrel where postId = $1", postId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var pic string
		rows.Scan(&pic)
		pics = append(pics, pic)
	}
	return pics
}

//func getPostById(postId int) (post Post, ItsComments []Com, ItsTags []Tag, ItsPics []Pic) {
//	rows, err := DB.Query("SELECT * FROM posts where id = $1", postId)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	for rows.Next() {
//		rows.Scan(&post.Id, &post.Text, &post.By, &post., &post.ViewCount, &post.Like, &post.CreatedOn, &post.UpdatedOn, &post.DeletedOn)
//	}
//
//	ItsPics = getPostsComments(postId)
//	ItsComments = getPostsPics(postId)
//	ItsTags = getPostsTags(postId)
//	return post, ItsComments, ItsTags, ItsPics
//}

func insertPost(post Post, picNames []string) (ItsId int) {
	err := DB.QueryRow(
		"INSERT INTO posts(text,by,createdOn) VALUES($1,$2,$3) returning id;",
		post.Text, post.By, getNow()).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}

	if affect := insertPicRel(ItsId, picNames); affect == 0{
		fmt.Println("inserting picrel failed.")
	}
	return ItsId
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
