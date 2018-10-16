package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main(){
	defer DB.Close()
	Routes()
}

func dbConn() *sql.DB {
	const (
		HOST     = "localhost"
		PORT     = 5432
		USER     = "postgres"
		PASSWORD = "1303"
		DBNAME   = "bluetalk"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

//func dbConn() *sql.DB {
//	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
//	if err != nil {
//		log.Fatalf("Error opening database: %q", err)
//	}
//	return db
//}

func tagFinder(body string) (tags []string) {
	split := strings.Split(body, "#")
	for i, sentence := range split {
		if !strings.HasPrefix(body, "#") {
			if i == 0 {
				continue
			}
		}
		words := strings.Split(sentence, " ")
		word := words[0]
		if !(strings.Contains(word, "(") && strings.Contains(word, ")")) {
			if strings.HasSuffix(word, ".") || strings.HasSuffix(word, "!") || strings.HasSuffix(word, ")") || strings.HasSuffix(word, "(") || strings.HasSuffix(word, ":") {
				word = word[0 : len(word)-1]
			}
		}
		tags = append(tags, word)
	}
	return tags
}
//
//func processTag(postId int, tags []string) {
//	for _, v := range tags {
//		v = strings.TrimSpace(v)
//		if len(v) == 0 {
//			continue
//		}
//		var tagId int
//		// we check to figure out if this tag_name already exists with
//		// getting its tag_id, if tag_id is not bigger than one, so
//		// there is no such tag_name
//		if tagId = alreadyTag(v); tagId < 1 {
//			tagId = insertTag(v)
//		}
//		_ = insertTagRel(tagId, postId)
//	}
//}

func getNow() string {
	splitTime := strings.Split(time.Now().String(), " ")
	now := splitTime[0]
	return now
}

func picSha(pic multipart.File) string {
	h := sha1.New()
	io.Copy(h, pic)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func processPostPic(file multipart.File) (picName string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	// create sha for file name
	picName = picSha(file) + ".jpg"
	path := filepath.Join(wd, "static", "pic", "stories", picName)
	nf, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer nf.Close()
	// copy
	file.Seek(0, 0)
	io.Copy(nf, file)
	return picName
}

func processProPic(file multipart.File, user User) (picName string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	// create sha for file name
	picName = picSha(file) + ".jpg"
	path := filepath.Join(wd, "static", "pic", "pros", picName)
	nf, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer nf.Close()
	// copy
	file.Seek(0, 0)
	io.Copy(nf, file)
	return picName
}

func detectFileType(file multipart.File) (valid bool) {
	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	var err error
	_, err = file.Read(buff)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filetype := http.DetectContentType(buff)

	switch filetype {
	case "image/jpeg", "image/jpg":
		return true

	case "image/gif":
		return true

	case "image/png":
		return true
	default:
		return false
	}
}

func getUniqueInt(intSlice []int) []int {
	keys := make(map[int]bool)
	var list []int
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
//
//func getUniqueTag(TagSlice []Tag) []Tag {
//	keys := make(map[Tag]bool)
//	var list []Tag
//	for _, entry := range TagSlice {
//		if _, value := keys[entry]; !value {
//			keys[entry] = true
//			list = append(list, entry)
//		}
//	}
//	return list
//}
