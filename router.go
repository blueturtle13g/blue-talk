package main

import (
	"github.com/alexedwards/scs"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
	"fmt"
)

var (
	tpl            *template.Template
	sessionManager = scs.NewCookieManager("cPfu>HIUkVBA1M7W/gNo+ZEjtp0}Yz-~Gv4lmdXyTOQ{$9r3Rs2^#nwqC8i6JK5D")
	DB             = dbConn()
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	sessionManager.Lifetime(time.Hour * 6) // Set the maximum session lifetime to 1 hour.
	sessionManager.Persist(false)          // Persist the session after a user has closed their browser.
}


//func port() string {
//	port := os.Getenv("PORT")
//	port = ":" + port
//	return port
//}


func Routes() {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/out", LogOutHandler)
	router.GET("/about", AboutHandler)
	//router.GET("/tag:id", TagHandler)
	//router.GET("/user/:id", UserHandler)
	router.GET("/profile/:id", ProfileHandler)
	router.POST("/profile/:id", ProfileProcess)

	router.GET("/", IndexHandler)
	//router.POST("/", IndexProcess)

	router.GET("/contact", ContactHandler)
	router.POST("/contact", ContactProcess)

	router.GET("/search", SearchHandler)
	router.POST("/search", SearchProcess)
	//
	router.GET("/register", RegisterHandler)
	router.POST("/register", RegisterProcess)
	//
	router.GET("/in", LogInHandler)
	router.POST("/in", LogInProcess)
	//
	//router.GET("/profile/:id/edit", EditProfileHandler)
	//router.POST("/profile/:id/edit", EditProfileProcess)
	//
	//router.GET("/makePost", MakePostHandler)
	//router.POST("/makePost", MakePostProcess)
	//
	//router.GET("/post/:id", PostHandlerHandler)
	//router.POST("/post/:id", PostProcess)
	//
	//router.GET("/post/:id/edit", EditPostHandler)
	//router.POST("/post/:id/edit", EditPostProcess)
	if err := http.ListenAndServe(":8080", sessionManager.Use(router)); err != nil{
		fmt.Println(err)
	}
}
