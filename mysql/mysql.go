package mysql

import (
	"../utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/core/errors"
	"log"
	"strings"
)

var Database *sql.DB

func Connect(url string) {
	var err error
	Database, err = sql.Open("mysql", url)
	utils.Check(err)
	fmt.Println("Connected to MySQL!")

	Execute("CREATE TABLE IF NOT EXISTS `answers` (`poll_id` text NOT NULL, `answer` text NOT NULL, `ip` text NOT NULL)")
	Execute("CREATE TABLE IF NOT EXISTS `polls` (`id` text NOT NULL, `question` text NOT NULL, `available_answers` text NOT NULL, `check_duplication` tinyint(1) NOT NULL)")
}

func Execute(sql string, params ... interface{}) int64 {
	if Database == nil {
		log.Fatal("You must connect to MySQL before use it!")
	}

	stmt, err := Database.Prepare(sql)
	utils.Check(err)

	res, err := stmt.Exec(params...)
	utils.Check(err)

	id, err := res.LastInsertId()
	utils.Check(err)

	return id
}

func Query(sql string, params ... interface{}) *sql.Rows {
	if Database == nil {
		log.Fatal("You must connect to MySQL before use it!")
	}

	stmt, err := Database.Prepare(sql)
	utils.Check(err)

	res, err := stmt.Query(params...)
	utils.Check(err)

	return res
}

func CreatePoll(poll Poll) {
	Execute("INSERT INTO `polls` (id, question, available_answers, check_duplication) VALUES (?, ?, ?, ?)", poll.Id, poll.Question, poll.ParseAvailableAnswers(), poll.CheckDuplicate)
}

func GetPoll(id string) (Poll, error){
	polls := Query("SELECT * FROM `polls` WHERE id = ?", id)

	var poll Poll
	for polls.Next() {
		var temp string
		err := polls.Scan(&poll.Id, &poll.Question, &temp, &poll.CheckDuplicate)
		utils.Check(err)

		poll.AvailableAnswers = strings.Split(temp, ";")
		break
	}

	if poll.Id == "" {
		return poll, errors.New("Poll with this ID does not exists!")
	}

	answers := Query("SELECT * FROM `answers` WHERE poll_id = ?", id)
	for answers.Next() {
		var answer Answer
		err := answers.Scan(&answer.PollId, &answer.Answer, &answer.Ip)
		utils.Check(err)
		poll.Answers = append(poll.Answers, answer)
	}

	return poll, nil
}

func CreateVote(answer Answer) {
	Execute("INSERT INTO `answers` (poll_id, answer, ip) VALUES (?, ?, ?)", answer.PollId, answer.Answer, answer.Ip)
}

func HadAlreadyVoted(id string, ip string) bool {
	res := Query("SELECT * FROM `answers` WHERE `poll_id` = ? AND `ip` = ?", id, ip)
	for res.Next() {
		return true
	}

	return false
}