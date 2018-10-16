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
	case len(user.Name) < 3:
		warnings = append(warnings, "Username cannot be less than 3.")
	case len(user.Name) > 20:
		warnings = append(warnings, "Username cannot be more than 20 characters.")
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
func (user User) UpValidate() (warnings []string) {
	if affect := alreadyName(user.Name, user.Id); affect > 0 {
		warnings = append(warnings, "This username already exists.")
	}
	if affect := alreadyEmail(user.Email, user.Id); affect > 0 {
		warnings = append(warnings, "This email already exists.")
	}
	if er := mailValidation(user.Email); er != "" {
		warnings = append(warnings, er)
	}
	if er := nameValidation(user.Name); er != "" {
		warnings = append(warnings, er)
	}
	switch {
	case len(user.Name) < 3:
		warnings = append(warnings, "Username Cannot Be Less Than 3 Characters.")
	case len(user.Quote) > 300:
		warnings = append(warnings, "The Length Of Your Quote Is Too Long(300).")
	case len(user.Name) > 20:
		warnings = append(warnings, "Username Cannot Be More Than 20 Characters.")
	case strings.Contains(user.Name, " "):
		warnings = append(warnings, "Username Cannot Have Space Between.")
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
		return "Please Insert A Valid Username Consist of English Letters."
	} else if err != nil {
		fmt.Println(err)
	}
	return ""
}
