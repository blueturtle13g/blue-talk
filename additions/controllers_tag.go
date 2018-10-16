package additions

import (
	"net/http"
	"strconv"

	"fmt"
	"github.com/julienschmidt/httprouter"
)

func TagHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := sessionManager.Load(r)
	var cUser User
	if UserId, err := session.GetInt("UserId"); err != nil {
		fmt.Println(err)
	} else {
		cUser = getUserById(UserId)
	}
	parameter := ps.ByName("id")
	TagId, _ := strconv.Atoi(parameter)
	gotFlash := make(map[string]interface{})
	tag := getTagById(TagId)
	gotFlash["Title"] = "#" + tag.Name
	gotFlash["UserName"] = cUser.Name
	gotFlash["Posts"] = getPostsByTagId(TagId)
	gotFlash["Tag"] = tag

	if err := tpl.ExecuteTemplate(w, "tag.gohtml", gotFlash); err != nil {
		fmt.Println(err)
	}
}
