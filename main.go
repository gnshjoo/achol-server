package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const VerifyMessage = "verified"

func AuthHAndler(next HandlerFunc) HandlerFunc {
	ignore := []string{"/login", "public/index.html"}
	return func(c *Context) {
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s){
				next(c)
				return
			}
		}
		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			c.Redirect("/login")
			return
		} else if err != nil {
			c.RenderErr(http.StatusInternalServerError, err)
			return

		} else if Verify(VerifyMessage, v.Value) {
			next(c)
			return
		}
		c.Redirect("/login")
	}
}

func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}

func CheckLogin(username, password string) bool {
	// 로그인 처리
	const (
		USERNAME = "tester"
		PASSWORD = "12345"
	)
	return username == USERNAME && password == PASSWORD
}

// 인증토큰 생성
func Sign(message string) string {
	secretKey := []byte("goland-book-secret-key2")
	if len(secretKey) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}

func main() {

	type User struct {
		Id string
		Addressid string
	}

	s := NewServer()

	s.HandleFunc("GET", "/", func(c *Context){
		c.RenderTemplate("/public/index.html", map[string]interface{}{"time": time.Now()})
	})


	s.HandleFunc("GET","/user/:id", logHandler(recoverHandler(func(c *Context){
		if c.Params["id"] == "0" {
			panic("id is zero")
		}
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v\n", c.Params["id"])
	})) )

	s.HandleFunc("POST", "/users", logHandler(recoverHandler(parseFormHandler(parseJsonBodyHandler(func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, c.Params)
	})))))

	s.HandleFunc("GET", "/users/:id", func(c *Context){
		u := User{Id: c.Params["id"].(string)}
		c.RenderXml(u)
	})

	s.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *Context) {
		u:= User{c.Params["user_id"].(string), c.Params["address_id"].(string)}
		c.RenderJson(u)
	})

	s.HandleFunc("GET", "/login", func (c *Context) {
		c.RenderTemplate("/public/login.html", map[string]interface{}{"message":"로그인이 필요합니다."})
	})

	s.HandleFunc("POST", "/login", func(c *Context){
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name: "X_AUTH",
				Value: Sign(VerifyMessage),
				Path: "/",
			})
			c.Redirect("/")
		}
		c.RenderTemplate("/public/login.html", map[string]interface{}{"message": "id 또는 password가 일치하지 않습니다."})
	})

	s.Use(AuthHAndler)

	s.Run(":8080")

}
