package server

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

type Users struct {
	ID       int64
	Username string
	Pass     string
}

type EventWatcher struct {
	db     *sql.DB
	mux    *http.ServeMux
	server *http.Server
}
type data struct {
	Title      string
	Coment     string
	Login      bool
	Register1  bool
	Register2  bool
	Conditions []*Users
}

type Point struct {
	X int
	Y int
}

var (
	//go:embed web/*.html
	content embed.FS
	tmpl    = template.Must(template.New("server").ParseFS(content, "web/*.html"))
)

func (ew *EventWatcher) HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleIndex")
	cs, err := ew.Conditions(r.Context(), 10)
	if err != nil {
		ew.error(w, err, http.StatusInternalServerError)
		return
	}

	data := &data{
		Title:      "Login",
		Login:      true,
		Register1:  true,
		Register2:  false,
		Conditions: cs,
	}

	if r.Method == "POST" {
		button := r.FormValue("button")
		fmt.Println("button:", button)

		switch button {
		case "Login":
			ew.HandleLogin(w, r, *data)

		case "Register":
			ew.HandleRegister(w, *data)

		case "Sign-up":
			ew.HandleSignUp(w, r, *data)
		default:

			return

		}

	} else {

		if err := tmpl.ExecuteTemplate(w, "login", data); err != nil {
			ew.error(w, err, http.StatusInternalServerError)
			return
		}

	}
}

// ログイン処理.
func (ew *EventWatcher) HandleLogin(w http.ResponseWriter, r *http.Request, data data) {
	id := r.FormValue("username")
	pass := r.FormValue("password")
	fmt.Println("name:", r.Form)

	for _, password := range data.Conditions {
		user := password.Username
		password := password.Pass
		if id == user && pass == password {
			fmt.Println("Login success")
			http.Redirect(w, r, "/game", http.StatusFound)
			return
		}
	}
	// Failed login
	fmt.Println("Login failed")
	data.Coment = "ユーザーネームまたはパスワードが間違っています"
	if err := tmpl.ExecuteTemplate(w, "login", data); err != nil {
		ew.error(w, err, http.StatusInternalServerError)
		return
	}
}

// 登録画面を生成.
func (ew *EventWatcher) HandleRegister(w http.ResponseWriter, data data) {
	fmt.Println("Register")
	data.Title = "Register"
	data.Login = false
	data.Register1 = false
	data.Register2 = true
	if err := tmpl.ExecuteTemplate(w, "login", data); err != nil {
		ew.error(w, err, http.StatusInternalServerError)
		return
	}
}

// 登録処理.
func (ew *EventWatcher) HandleSignUp(w http.ResponseWriter, r *http.Request, data data) {
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if user == "" {
		fmt.Println("Register failed")
		data.Title = "Register"
		data.Coment = "ユーザーネームが指定されていません"
		data.Login = false
		data.Register1 = false
		data.Register2 = true
		if err := tmpl.ExecuteTemplate(w, "login", data); err != nil {
			ew.error(w, err, http.StatusInternalServerError)
			return
		}
		return
	} else if pass == "" {
		fmt.Println("Register failed")
		data.Title = "Register"
		data.Coment = "パスワードが指定されていません"
		data.Login = false
		data.Register1 = false
		data.Register2 = true
		if err := tmpl.ExecuteTemplate(w, "login", data); err != nil {
			ew.error(w, err, http.StatusInternalServerError)
			return
		}
		return
	} else {
		id := r.FormValue("username")
		for _, password := range data.Conditions {
			user := password.Username
			if id == user {
				fmt.Println("Register failed")
				data.Title = "Register"
				data.Coment = "ユーザーネームが既に登録されています"
				data.Login = false
				data.Register1 = false
				data.Register2 = true
				if err := tmpl.ExecuteTemplate(w, "login", data); err != nil {
					ew.error(w, err, http.StatusInternalServerError)
					return
				}
				return

			}
		}

		c := &Users{
			Username: user,
			Pass:     pass,
		}
		if err := ew.AddCondition(r.Context(), c); err != nil {
			ew.error(w, err, http.StatusInternalServerError)
			return
		}
		// Successful login
		fmt.Println("Register success")
		http.Redirect(w, r, "/dashbord", http.StatusFound)
		return
	}
}

// ログイン後のゲーム画面生成.
func GameHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("server/web/index.html"))
	tmpl.Execute(w, nil)
}

// ランダムな座標を生成.
func PointHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	p := Point{
		X: rand.Intn(800),
		Y: rand.Intn(600),
	}
	fmt.Fprintf(w, `{"x": %d, "y": %d}`, p.X, p.Y)
}

// エラー処理.
func (ew *EventWatcher) error(w http.ResponseWriter, err error, code int) {
	log.Println("Error:", err)
	http.Error(w, http.StatusText(code), code)
}

// ハンドラの初期化.
func (ew *EventWatcher) InitHandlers() {
	ew.mux.HandleFunc("/", ew.HandleIndex)
	ew.mux.HandleFunc("/game", GameHandler)
	ew.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
}
