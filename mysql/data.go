package mysql

import "strings"

type Poll struct {
	Id string
	Question string
	AvailableAnswers []string
	Answers []Answer
	CheckDuplicate bool
}

type Answer struct {
	PollId string
	Answer string
	Ip string
}

func (poll Poll) ParseAvailableAnswers() string {
	var s string

	for _, value := range poll.AvailableAnswers {
		s += value + ";"
	}

	return strings.TrimSuffix(s, ";")
}
