package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gorilla/mux"
)

var db *gorm.DB
var err error


type Todo struct {
    gorm.Model
    Title  string
    Status bool
}


func initDB() {
    dsn := "root:kl18jda183079@tcp(127.0.0.1:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    db.AutoMigrate(&Todo{})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    var todos []Todo
    db.Find(&todos)
    tmpl, _ := template.ParseFiles("todo.html")
    tmpl.Execute(w, todos)
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
    title := r.FormValue("title")
    if title != "" {
        db.Create(&Todo{Title: title, Status: false})
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}


func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    var todo Todo
    db.First(&todo, id)
    todo.Status = !todo.Status
    db.Save(&todo)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}


func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    var todo Todo
    db.Delete(&todo, id)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
    initDB()

    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler).Methods("GET")
    r.HandleFunc("/todo", createTodoHandler).Methods("POST")
    r.HandleFunc("/todo/{id}/update", updateTodoHandler).Methods("POST")
    r.HandleFunc("/todo/{id}/delete", deleteTodoHandler).Methods("POST")

    fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
