package additions

import (
	"fmt"
	_ "github.com/lib/pq"
	"database/sql"
)

type Post struct {
	Id, ViewCount , Like                            int
	Text, By, CreatedOn, UpdatedOn, DeletedOn       string
	Pics []string
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
		posts = append(posts, post)
	}
	return posts
}

func getAllPosts(order string) (posts []Post) {
	var rows *sql.Rows
	var err error
	if order == "rate" {
		rows, err = DB.Query("select * from posts order by Like desc")
	}else{
		rows, err = DB.Query("select * from posts order by id desc")
	
	}
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

func getPostById(postId int) (post Post, ItsComments []Com, ItsTags []Tag, ItsPics []Pic) {
	rows, err := DB.Query("SELECT * FROM posts where id = $1", postId)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		rows.Scan(&post.Id, &post.Text, &post.By, &post., &post.ViewCount, &post.Like, &post.CreatedOn, &post.UpdatedOn, &post.DeletedOn)
	}
	
	ItsPics = getPostsComments(postId)
	ItsComments = getPostsPics(postId)
	ItsTags = getPostsTags(postId)
	return post, ItsComments, ItsTags, ItsPics
}

func insertPost(post Post) (ItsId int) {
	err := DB.QueryRow(
		"INSERT INTO posts(text,by,createdOn) VALUES($1,$2,$3) returning id;",
		post.Text, post.By, getNow()).Scan(&ItsId)
	if err != nil {
		fmt.Println(err)
	}

	return ItsId
}
