package web

import (
	"../mysql"
	"../utils"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"strconv"
)

func Start() {
	app := iris.Default()

	tmpl := iris.Django("./website", ".html")
	tmpl.Reload(true)
	app.RegisterView(tmpl)
	app.StaticServe("./website", "/")

	app.Get("/", root)
	app.Get("/rest/{id:string}", rest)
	app.Get("/{id:string}", poll)
	app.Post("/create", create)
	app.Post("/vote", vote)

	fmt.Println("Everything seems to work!")
	err := app.Run(iris.Addr(":8080"))
	utils.Check(err)
}

func root(ctx context.Context) {
	ctx.ViewData("error", ctx.URLParam("error"))
	ctx.View("index.html")
}

func rest(ctx context.Context) {
	id := ctx.Params().Get("id")

	poll, err := mysql.GetPoll(id)
	if err != nil {
		_, _ = ctx.JSON(iris.Map{})
	}

	amount := make(map[string]int)

	for _, value := range poll.AvailableAnswers {
		amount[value] = 0
	}

	for _, value := range poll.Answers {
		amount[value.Answer] = amount[value.Answer] + 1
	}

	_, err = ctx.JSON(iris.Map{
		"id":                poll.Id,
		"question":          poll.Question,
		"available_answers": poll.AvailableAnswers,
		"answer_amount":     amount,
	})

	utils.Check(err)
}

func poll(ctx context.Context) {
	poll, err := mysql.GetPoll(ctx.Params().Get("id"))


	if err != nil {
		ctx.View("index.html")
		return
	}

	ctx.ViewData("error", ctx.URLParam("error"))
	ctx.ViewData("AvailableAnswers", poll.AvailableAnswers)
	ctx.ViewData("id", poll.Id)
	ctx.View("poll.html")
}

func vote(ctx context.Context) {
	poll, err := mysql.GetPoll(ctx.FormValue("id"))

	if err != nil {
		ctx.Redirect("/")
		return
	}

	ip := ctx.RemoteAddr()
	a := ctx.FormValue("answer")

	if mysql.HadAlreadyVoted(poll.Id, ip) && poll.CheckDuplicate {
		ctx.Redirect("/" + poll.Id + "?error=1")
		return
	}

	answer := mysql.Answer{
		PollId: poll.Id,
		Answer: a,
		Ip: ip,
	}

	mysql.CreateVote(answer)
	ctx.Redirect("/" + poll.Id)
}

func create(ctx context.Context) {
	amount, _ := strconv.Atoi(ctx.FormValue("answers-amount"))

	var answers []string
	notEmpty := readAnswers(ctx, amount, &answers)

	if notEmpty < 2 {
		ctx.Redirect("/?error=1")
		return
	}

	if notEmpty > 10 {
		ctx.Redirect("/?error=4")
		return
	}

	var poll mysql.Poll
	poll.Id = utils.GetRandId(10)
	poll.Question = ctx.FormValue("question")
	poll.AvailableAnswers = answers
	poll.CheckDuplicate = ctx.FormValue("check_duplication") == "on"

	if len(poll.Question) < 3 {
		ctx.Redirect("/?error=2")
		return
	}

	for index1, value1 := range poll.AvailableAnswers {
		if len(value1) < 3 {
			ctx.Redirect("/?error=2")
			return
		}
		for index2, value2 := range poll.AvailableAnswers {
			if value1 == value2 && index1 != index2 {
				ctx.Redirect("/?error=3")
				return
			}
		}
	}


	mysql.CreatePoll(poll)
	ctx.Redirect("/" + poll.Id)
}

func readAnswers(ctx context.Context, amount int, answers *[]string) int{
	notEmpty := 0
	for i := 1; i <= amount; i++ {
		current := ctx.FormValue(fmt.Sprintf("answer%d", i))

		if current != "" {
			*answers = append(*answers, current)
			notEmpty++
		}
	}

	return notEmpty
}