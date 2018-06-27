package httphandler

import (
	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context) {

	s := `
<html>
    <body>
		<form class="ui form" method="POST" action="/api/login">
			<h4>Login</h4>
			<div>
				<label>Email</label>
				<input type="text" placeholder="email" name="email">
			</div>
			<div>
				<label>Password</label>
				<input type="password" placeholder="password" name="password">
			</div>
			<input type="submit" value="Submmit">
		</form>
    </body>
</html>
	`

	c.Data(200, "text/html", []byte(s))
}
