package main

// CGO_ENABLED=0 go build main.go
import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
	"gopkg.in/gin-gonic/gin.v1"
)

//Data asdsadas
type Data struct {
	Five   int64
	Hourly int64
	Daily  int64
	First  int64
	Next   time.Time
	Local  time.Time
}

//Datos a exportar
var Datos Data

func updateF() {
	currentTime := time.Now().Local()
	fmt.Println("Its one more hour, The Current time is ", currentTime.Format("02-01-2006"))
	Datos.Five++
}
func updateH() {
	Datos.Hourly++
	fmt.Println("Its been 8 hrsalready")

}
func updateD() {
	fmt.Println("Another day under the sun")
	currentTime := time.Now().Local()
	fmt.Println("Current date is: ", currentTime.Format("02-01-2006"))
	Datos.Daily++
	//This is updated daily so we check if it is our desired date every day
	day := currentTime.Day()
	if day == 17 || day == 18 {
		Datos.First++
		fmt.Println("This is the :", day, ": Day of the Moth!")
	}
}

func ServeHTTP(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	//var req *http.Request = c.Req

	//Parsing HTML
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err)
	}

	Datos.Local = time.Now().Local()
	items := struct {
		Five   string
		Hourly string
		Daily  string
		First  string
		Next   time.Time
		Local  time.Time
	}{
		Five:   strconv.FormatInt(Datos.Five, 10),
		Hourly: strconv.FormatInt(Datos.Hourly, 10),
		Daily:  strconv.FormatInt(Datos.Daily, 10),
		First:  strconv.FormatInt(Datos.First, 10),
		Next:   Datos.Next,
		Local:  Datos.Local,
	}

	err = t.ExecuteTemplate(w, "index.html", items)
	if err != nil {
		fmt.Println(err)
	}

}
func task() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

// Listen (just commented to avoid annoying vscode underlines)
func Listen() {
	//fs := http.FileServer(http.Dir("public"))
	//http.Handle("/static/", fs)
	//http.HandleFunc("/", ServeHTTP)
	//port := os.Getenv("PORT")
	//fmt.Println(port)
	//http.ListenAndServe(port, nil)
}

func main() {
	//go Listen()
	r := gin.Default()
	//r.Static("/css", "../templates/css")
	r.GET("/", func(c *gin.Context) {
		ServeHTTP(c)
	})
	//go r.Run()
	//r.Run(":" + os.Getenv("PORT"))
	go r.Run(":8000")
	go updateD()
	//gocron.Every(10).Seconds().Do(updateD)
	// Do jobs without params
	gocron.Every(2).Minutes().Do(updateD)
	gocron.Every(5).Minutes().Do(updateF)
	gocron.Every(1).Hour().Do(updateH)
	gocron.Every(1).Day().Do(updateD)
	// remove, clear and next_run
	_, time := gocron.NextRun()
	Datos.Next = time
	fmt.Println(time)

	// function Start start all the pending jobs
	<-gocron.Start()
}
