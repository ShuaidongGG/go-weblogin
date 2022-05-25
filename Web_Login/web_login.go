package main

import (
	"database/sql"
	"fmt"
	"net/http"
	_"mysql"
	"os/exec"
	"runtime"
	//"html/template"
	"text/template"
	//"github.com/jmoiron/sqlx"
)

var commands = map[string]string{
	"windows": "cmd /c start",
	"darwin": "open",
	"linux": "xdg-open",
}

func Open(url string) error{
	run,ok :=commands[runtime.GOOS]
	if !ok{
		return fmt.Errorf("No os like %s",runtime.GOOS)
	}
	cmd := exec.Command(run,url)
	return cmd.Start()
}
/*
func OpenW(url string) error{
	run,ok := commands["windows"]
	if !ok{
		return fmt.Errorf("No os like windows")
	}
	cmd := exec.Command(run,url)
	return cmd.Start()
}*/

var db *sql.DB

func initDB() (err error) {
	//dsn := "root:123456@tcp(127.0.0.1:3306)/db1"
	//dsn := "root:123456@tcp(127.0.0.1:3306)/db1"
	db, err = sql.Open("mysql", "root:123456@tcp(172.17.0.2:3306)/db1")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func check(account, password string) bool {
	sqlStr := `select * from login where account = ? and password =?`
	rows, err := db.Query(sqlStr, account, password)
	if err != nil {
		return false
	}
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}

func checkA(account string) bool {
	sqlStr := `select * from login where account = ?`
	rows, err := db.Query(sqlStr, account)
	if err != nil {
		return false
	}
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}

func addUser(account, password string) {
	sqlStr := `insert into login(account,password)values(?,?)`
	rows, err := db.Query(sqlStr, account, password)
	if err != nil {
		fmt.Println("sql is illegal")
	}
	defer rows.Close()
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("Login.html")
		t.Execute(w,nil)
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			return
		}
		account := r.Form.Get("account")
		password := r.Form.Get("password")
		flag := check(account, password)
		if account != "" {
		if flag == true {
		//fmt.Fprintln(w, "yes")
		//Suc(w,r)
		//t,_ :=template.ParseFiles("Suc.html")
		//_ = t.Execute(w,nil)
		//Open("http://localhost:9090/Suc")
		// w.Header().Set("Location", "/Suc") 
		// w.WriteHeader(301)
		t,_ := template.ParseFiles("Suc.html")
		t.Execute(w,nil)
		} else {
		//fmt.Fprintln(w, "no")
		//Failed(w,r)
		//Open("http://localhost:9090/Failed")
		// w.Header().Set("Location", "/Failed") 
		// w.WriteHeader(301)
		w.Write([]byte("<script>alert('账号或者密码不正确')</script>"))
		t, _ := template.ParseFiles("Login.html")
		t.Execute(w,nil)
	}
	}
	}
	
	
}
 
func Reg(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t,_ := template.ParseFiles("Reg.html")
		t.Execute(w,nil)
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			return
		}
		account := r.Form.Get("account")
		password := r.Form.Get("password")
		flag := checkA(account)
		if flag == false {
			addUser(account,password)
			w.Header().Set("location","/Login")
			w.WriteHeader(301)
		} else {
			w.Write([]byte("<script>alert('该用户名已存在')</script>"))
			t, _ := template.ParseFiles("Reg.html")
			t.Execute(w,nil)
		}

	}
}
/*
func Suc(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles("Suc.html")
	t.Execute(w,nil)
}

func Failed(w http.ResponseWriter, r *http.Request) {
	t,_ := template.ParseFiles("Failed.html")
	t.Execute(w,nil)
}*/

func main() {
	err := initDB()
	if err != nil {
		fmt.Println("err", err)
		return
	}
	http.HandleFunc("/Login", Login)
	http.HandleFunc("/Reg",Reg)
	//http.HandleFunc("/Suc",Suc)
	//http.HandleFunc("/Failed",Failed)
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("err", err)
		return
	}
}
