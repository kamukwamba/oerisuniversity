package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)

type MessageAdmin struct {
	StInfo   StudentInfo
	Messages []MessageOut
}

type MessageOut struct {
	UUID         string
	Student_UUID string
	Sender_Name  string
	Message      string
	Sender       bool
	Seen_Student bool
	Seen_Admin   bool
	Date         string
}

func DeleteMessage(message_uuid string) bool {

	dbdelete := dbcode.SqlRead().DB

	deleted := true

	delete, err := dbdelete.Prepare("DELETE FROM messages WHERE uuid = ?")
	if err != nil {
		errout := fmt.Sprintf("line 30 erro: %s", err)
		ErrorPrintOut("messages", "DeleteMessages", errout)
	}
	defer delete.Close()

	_, errd := delete.Exec(message_uuid)
	if errd != nil {
		errout := fmt.Sprintf("line 42 erro: %s", err)
		ErrorPrintOut("messages", "DeleteMessages", errout)
	}

	return deleted

}

func UpdateMessages(student_uuid string) bool {
	dbupdate := dbcode.SqlRead().DB

	updated := true

	update, err := dbupdate.Prepare("UPDATE messages SET seen_admin = ? WHERE sender_uuid = ? ")

	if err != nil {
		errout := fmt.Sprintf("line 30 erro: %s", err)
		ErrorPrintOut("messages", "UpdateMessages", errout)
	}

	defer update.Close()

	_, errupdate := update.Exec(true, student_uuid)

	if errupdate != nil {
		errout := fmt.Sprintf("line 48 error: %s", err)
		ErrorPrintOut("messaes", "UpdateMessages", errout)
	}

	return updated

}

func ReadMessage(student_uuid string) MessageAdmin {
	dbread := dbcode.SqlRead().DB
	var message_out MessageOut

	student_data := GetStudentAllDetails(student_uuid)

	var message_out_list []MessageOut

	statement, err := dbread.Query("select uuid,sender_uuid,sender_name,sender,message,seen_student,seen_admin,date from messages")

	if err != nil {
		errorout := fmt.Sprintf("line 22: %s", err)
		ErrorPrintOut("messasges.go", "ReadMessages", errorout)
	}

	defer statement.Close()

	for statement.Next() {
		err = statement.Scan(&message_out.UUID, &message_out.Student_UUID, &message_out.Sender_Name, &message_out.Sender, &message_out.Message, &message_out.Seen_Student, &message_out.Seen_Admin, &message_out.Date)

		if err != nil {
			errorout := fmt.Sprintf("line 48: %s", err)
			ErrorPrintOut("messages.go", "ReadMessages", errorout)
		}

		if message_out.Student_UUID == student_uuid {
			message_out_list = append(message_out_list, message_out)

		}

	}

	student_msg := MessageAdmin{
		StInfo:   student_data,
		Messages: message_out_list,
	}

	fmt.Println(student_msg)

	return student_msg

}

func LoadMessages() []MessageOut {
	var message_out_list []MessageOut
	var message_out MessageOut

	var message_seen []MessageOut
	var message_unseen []MessageOut

	dbread := dbcode.SqlRead().DB

	statement, err := dbread.Query("select uuid,sender_uuid,sender_name, sender, message,seen_admin,date from messages")

	if err != nil {
		log.Fatal(err)
	}

	defer statement.Close()

	for statement.Next() {
		err := statement.Scan(&message_out.UUID, &message_out.Student_UUID, &message_out.Sender_Name, &message_out.Sender, &message_out.Message, &message_out.Seen_Admin, &message_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("line 42: %s", err)
			ErrorPrintOut("messages", "LoadMessages", error_out)
		}

		if message_out.Sender {
			if !message_out.Seen_Admin {
				if len(message_unseen) > 0 {
					for _, message := range message_unseen {
						if message_out.Student_UUID != message.Student_UUID {
							message_unseen = append(message_unseen, message_out)
						}
					}
				} else {
					message_unseen = append(message_unseen, message_out)
				}

			} else {
				if len(message_seen) > 0 {
					for _, message := range message_seen {
						if message_out.Student_UUID != message.Student_UUID {
							message_seen = append(message_seen, message_out)
						}
					}
				} else {
					message_seen = append(message_seen, message_out)
				}
			}

		}

		if len(message_unseen) > 0 {
			for _, message := range message_unseen {
				if len(message_out_list) > 0 {
					for _, messageout := range message_out_list {
						if message.Student_UUID != messageout.Student_UUID {
							message_out_list = append(message_out_list, message)
						}
					}
				} else {
					message_out_list = append(message_out_list, message)
				}
			}
		} else {
			for _, message := range message_unseen {
				message_out_list = append(message_out_list, message)
			}
		}

		if len(message_seen) > 0 {
			for _, message := range message_seen {
				if len(message_out_list) > 0 {
					for _, messageout := range message_out_list {
						if messageout.Student_UUID != message.Student_UUID {
							message_out_list = append(message_out_list, message)
						}
					}
				} else {
					message_out_list = append(message_out_list, message)
				}
			}
		} else {
			for _, message := range message_seen {
				message_out_list = append(message_out_list, message)
			}
		}
	}

	fmt.Println(message_out_list)

	return message_out_list

}

func AdminMessagesPage(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	messages_out := LoadMessages()

	err := tpl.ExecuteTemplate(w, "messageAdmin", messages_out)

	if err != nil {
		log.Fatal(err)
	}
}
