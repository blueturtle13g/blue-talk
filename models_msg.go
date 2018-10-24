package main

import "fmt"

type Msg struct {
	Id                                                      int
	Text, FromName, ToName, ToPic, CreatedOn, UpdatedOn, DeletedOn string
	Seen                                                    bool
}

type MsgH struct {
	To, LastMsg, Pic string
}

func getMsgHeaders(username string) (uniqueMsgHeaders []MsgH) {
	// first we should find out, who has some messages with this guy
	rows, err := DB.Query("select toName from msgs where fromName = $1", username)
	if err != nil {
		fmt.Println(err)
	}
	var toNames []string
	for rows.Next() {
		var toName string
		rows.Scan(&toName)
		toNames = append(toNames, toName)
	}
	// now we delete repeated destinations(to) from our list
	// so we have all unique destinations(tos)
	uniqueToNames := getUniqueString(toNames)
	// now it's time to get the last message for each destinations(tos)
	for _, uniqueToName := range uniqueToNames{
		// here we put a limitation to get just the last one.
		var pic string
		rows, err := DB.Query("select pic from users where name = $1;", uniqueToName)
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			rows.Scan(&pic)
		}
		rows, err = DB.Query("select text from msgs where fromName = $1 and toName = $2 order by id desc limit 1;", username, uniqueToName)
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			var text string
			rows.Scan(&text)
			uniqueMsgHeaders = append(uniqueMsgHeaders, MsgH{uniqueToName, text, pic})
		}
	}
	return uniqueMsgHeaders
}

func getMsgs(username, toName string) (msgs []Msg) {
	// first we should find out, who has some messages with this guy
	rows, err := DB.Query("select * from msgs where fromName = $1 and toName = $2 order by id asc", username, toName)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var msg Msg
		rows.Scan(&msg.Id, &msg.Text, &msg.FromName, &msg.ToName, &msg.Seen, &msg.CreatedOn, &msg.UpdatedOn, &msg.DeletedOn )
		rows, err := DB.Query("select pic from users where name = $1;", msg.ToName)
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			rows.Scan(&msg.ToPic)
		}
		msgs = append(msgs, msg)
	}
	return msgs
}
