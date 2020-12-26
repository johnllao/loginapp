package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	AuthToken = "auth-token"
	SignedKey = "123"
)

var (
	indextempl *template.Template
	logintempl *template.Template
)

func main() {
	starthttp()
}

func starthttp() {
	var err error

	indextempl, err = template.New("index").Parse(indexhtml)
	if err != nil {
		log.Fatal(err)
	}

	logintempl, err = template.New("login").Parse(loginhtml)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("application started")
	var s = http.Server {
		Addr:    "localhost:8080",
		Handler: http.HandlerFunc(roothandler),
	}
	log.Fatal(s.ListenAndServe())
}

func roothandler(w http.ResponseWriter, r *http.Request) {

	var err error

	var authtok = r.Header.Get(AuthToken)
	if authtok == "" {
		var usr = r.FormValue("u")
		var pwd = r.FormValue("p")

		if usr == "" {
			w.Header().Set("Content-Type", "text/html")
			logintempl.Execute(w, nil)
			return
		}

		if usr != "admin" {
			w.Header().Set("Content-Type", "text/html")
			logintempl.Execute(w, "Invalid user name and/or password")
			return
		}

		if pwd != "admin" {
			w.Header().Set("Content-Type", "text/html")
			logintempl.Execute(w, "Invalid user name and/or password")
			return
		}

		// set claims
		var claims = make(jwt.MapClaims)
		claims["user"]   = usr
		claims["expiry"] = time.Now().Add(time.Hour * 1).Unix()

		// Create the token
		var tok = jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

		// Sign and get the complete encoded token as a string
		var signedtok string
		signedtok, err = tok.SignedString([]byte(SignedKey))
		if err != nil {
			w.Header().Set("Content-Type", "text/html")
			logintempl.Execute(w, err.Error())
			return
		}

		w.Header().Set(AuthToken, signedtok)
		w.Header().Set("Content-Type", "text/html")
		indextempl.Execute(w, nil)
		return
	} else {

		var tok *jwt.Token
		tok, err = jwt.Parse(authtok, func(token *jwt.Token) (interface{}, error) {
			return []byte(SignedKey), nil
		})
		if err != nil {
			w.Header().Set("Content-Type", "text/html")
			logintempl.Execute(w, err.Error())
			return
		}
		if !tok.Valid {
			w.Header().Set("Content-Type", "text/html")
			logintempl.Execute(w, "Invalid authentication token")
			return
		}

		w.Header().Set("Content-Type", "text/html")
		indextempl.Execute(w, nil)
	}

}


var indexhtml = `<!DOCTYPE html>
<html>
<head>
	<title>My Page</title>
</head>
<body>
	<div>
		Welcome!
	</div>
</body>
</html>
`

var loginhtml = `<!DOCTYPE html>
<html>
<head>
	<title>My Page</title>
	<style type="text/css">
	body, div {
		font-family: Arial;
		font-size: 12pt;
		margin: 0;
		padding: 0;
	}
	button, input[type=text], input[type=password] {
		padding: .5em;
		width: 100%;
	}
	.row {
		display: block;
		margin: .3em;
	}
	.cell {
		display: inline-block;
	}
	.button {
		width: 120px;
	}
	.label {
		width: 120px;
	}
	.textbox {
		width: 180px;
	}
	</style>
</head>
<body>
	<form method="POST" action="/">
		<div class="row">
			<div class="cell label"><label>User name</label></div>
			<div class="cell textbox"><input type="text" name="u" /></div>
		</div>
		<div class="row">
			<div class="cell label"><label>Password</label></div>
			<div class="cell textbox"><input type="password" name="p" /></div>
		</div>
		<div class="row">
			<div class="cell button">
				<button>Login</button>
			</div>
			<div class="cell">
				&nbsp;{{ . }}
			</div>
		</div>
	</form>
</body>
</html>
`
