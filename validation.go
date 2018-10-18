package main

import (
	"fmt"
	"regexp"
	"strings"
)

func (user User) Validate(pwd, confirm string) (warnings []string) {
	if affect := alreadyName(user.Name, 0); affect != 0 {
		warnings = append(warnings, "This username already exists.")
	}
	if affect := alreadyEmail(user.Email, 0); affect != 0 {
		warnings = append(warnings, "This email already exists.")
	}
	if er := mailValidation(user.Email); er != "" {
		warnings = append(warnings, er)
	}
	switch {
	case len(user.Name) < 3 || len(user.FirstName) < 3 || len(user.LastName) < 3:
		warnings = append(warnings, "Names cannot be less than 3 Characters.")
	case len(user.Name) > 20 || len(user.FirstName) > 20 || len(user.LastName) > 20:
		warnings = append(warnings, "Names cannot be more than 20 characters.")
	case strings.Contains(user.Name, " "):
		warnings = append(warnings, "Username cannot have space between.")
	case len(pwd) < 8:
		warnings = append(warnings, "Password cannot be less than 8 characters.")
	case len(pwd) > 50:
		warnings = append(warnings, "Password cannot be more than 50 characters.")
	case pwd != confirm:
		warnings = append(warnings, "Password does not match with confirm-password.")
	}
	return warnings
}
func (user User) UpValidate() (warnings string) {
	if affect := alreadyName(user.Name, user.Id); affect > 0 {
		warnings += "This username already exists. \n"
	}
	if affect := alreadyEmail(user.Email, user.Id); affect > 0 {
		warnings += "This email already exists. \n"
	}
	if er := mailValidation(user.Email); er != "" {
		warnings +=  er+" \n"
	}
	if er := nameValidation(user.Name); er != "" {
		warnings +=  er+" \n"
	}
	switch {
	case len(user.Name) < 3 || len(user.FirstName) < 3 || len(user.LastName) < 3:
		warnings +="Names cannot be less than 3 Characters."
	case len(user.Phone) < 11 && len(user.Phone) != 0:
		warnings +="Phone cannot be less than 11 Characters."
	case len(user.Name) > 20 || len(user.FirstName) > 20 || len(user.LastName) > 20:
		warnings +="Names cannot be more than 20 characters."
	case len(user.Quote) > 300:
		warnings += "The Length Of Your Quote Is Too Long(300). \n"
	case strings.Contains(user.Name, " "):
		warnings += "Username Cannot Have Space Between. \n"
	}
	return warnings
}

//func (post Post) Validate(condition string, tags []string) (warnings []string) {
//	var postId int
//	if condition == "up" {
//		postId = post.Id
//	} else {
//		postId = 0
//	}
//	if affect := alreadyTitle(post.Title, postId); affect != 0 {
//		warnings = append(warnings, "This Title already exists, please choose another title for your post.")
//	}
//	if affect := alreadyBody(post.Title, postId); affect != 0 {
//		warnings = append(warnings, "This Title already exists, please choose another title for your post.")
//	}
//	switch {
//	case len(post.Text) < 100:
//		warnings = append(warnings, "Post can't be less than 100 characters.")
//	}
//}

func (msg Cnt) Validate() (warnings []string) {
	if len(msg.Name) < 3 || len(msg.Name) > 20 {
		warnings = append(warnings, "The length of your name can't be less than 3 or more than 20 characters.")
	}
	if len(msg.Text) == 0 {
		warnings = append(warnings, "You Are Sending Us An Empty Message, Please Check.")
	}
	if er := mailValidation(msg.Email); er != "" {
		warnings = append(warnings, er)
	}
	if er := nameValidation(msg.Name); er != "" {
		warnings = append(warnings, er)
	}
	return warnings
}

func mailValidation(email string) (er string) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(email) {
		return ""
	}
	return "you've entered an invalid email address, please check it."
}

func nameValidation(name string) (er string) {
	if valid, err := regexp.MatchString("^[a-zA-Z0]+([_ -]?[a-zA-Z0-9])*$", name); !valid {
		return "invalid"
	} else if err != nil {
		fmt.Println(err)
	}
	if strings.Contains(name, " ") {
		return "invalid"
	}
	return ""
}
