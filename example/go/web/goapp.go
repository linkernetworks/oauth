package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	clientId     string
	clientSecret string
)

func main() {
	flag.StringVar(&clientId, "i", "MB8CAQACAgw3AgMBAAECAgcZAgE7AgE1AgETAgExAgEx", "id of client")
	flag.StringVar(&clientSecret, "s", "MAkCAgw3AgMBAAE=", "secret id of client")
	flag.Parse()

	// Application home endpoint
	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		u := genDefaultValues()
		codeUrl := fmt.Sprintf("<a href=\"http://localhost:9096/services/oauth/authorize?%s\">Code</a><br/>", u.Encode())

		w.Write([]byte("<html><body>"))
		w.Write([]byte(codeUrl))

		w.Write([]byte("</body></html>"))
	})

	http.HandleFunc("/callback/code", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - CODE<br/>"))
		defer w.Write([]byte("</body></html>"))

		r.ParseForm()
		code := r.Form.Get("code")
		redirectUri := r.Form.Get("redirect_uri")
		if code == "" {
			w.Write([]byte("Nothing to do"))
			return
		}

		u, err := url.Parse("http://localhost:9096/services/oauth/token")
		if err != nil {
			w.Write([]byte("internal error: " + err.Error()))
			return
		}
		p := u.Query()
		p.Set("grant_type", "authorization_code")
		p.Set("code", code)
		p.Set("redirect_uri", "")
		p.Set("client_id", clientId)
		p.Set("client_secret", clientSecret)
		p.Set("redirect_uri", redirectUri)
		u.RawQuery = p.Encode()

		resp, err := http.Get(u.String())
		if err != nil {
			w.Write([]byte("get request to token endpoint error: " + err.Error()))
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.Write([]byte("parse response body error: " + err.Error()))
			return
		}

		w.Write(body)
	})

	http.ListenAndServe(":14000", nil)
}

func genDefaultValues() url.Values {
	u := url.Values{}
	u.Set("response_type", "code")
	u.Set("client_id", clientId)
	u.Set("state", "xyz")
	u.Set("user_email", "dev@linker.com")
	u.Set("redirect_uri", url.QueryEscape("http://localhost:14000/callback/code"))

	return u
}
