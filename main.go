package main

import (
	"fmt"
	"net/http"

	emailerifier "github.com/AfterShip/email-verifier"
	"github.com/gin-gonic/gin"
)

var verifier *emailerifier.Verifier

func main() {

	verifier := emailerifier.NewVerifier()
	//verifier = verifier.EnableSMTPCheck()

	verifier = verifier.EnableDomainSuggest()
	verifier = verifier.AddDisposableDomains([]string{"tractorjj.com"})

	router := gin.Default()

	router.LoadHTMLFiles("./templates/ver-email.html", "./templates/ver-result.html")

	router.GET("/verifyemail", EmailGetHandler)
	router.POST("/verifyemail", EmailPostHandler)

	router.Run("localhost:8080")
}

/* ----------------------------- Handlers ------------------------------------  */

func EmailGetHandler(c *gin.Context) {

	c.HTML(http.StatusOK, "ver-email.html", nil)
}

func EmailPostHandler(c *gin.Context) {

	email := c.PostForm("email")

	ret, err := verifier.Verify(email)
	if err != nil {

		fmt.Println("error in email syntax")
		c.HTML(http.StatusInternalServerError, "ver-email.html", gin.H{"Message": "unable to register email"})
		return
	}

	fmt.Println("email validation result ", ret)
	fmt.Println("your email : ", ret.Email, "\n Recheable :", ret.Reachable, "\n syntax :", ret.Syntax, "\nSMTP:")

	if !ret.Syntax.Valid {

		fmt.Println("syntax is invalid check again ")
		c.HTML(http.StatusBadRequest, "ver-email.html", gin.H{"Message": "your email address is  not valid syntax error"})
		return
	}

	if ret.Disposable {

		fmt.Println("we don't accept disposable email addresses")

		c.HTML(http.StatusBadRequest, "ver-email.html", gin.H{"Message": "we don't accept disposable email addresses"})
		return
	}

	c.HTML(http.StatusOK, "ver-result.html", gin.H{"email": email})

}
