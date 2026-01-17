package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConfig *oauth2.Config
)

func homeHandler(c *gin.Context) {
	if c == nil {
		fmt.Printf("c is nil")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("c is nil")})
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"LoginURL": "/auth/login",
	})
}

func loginHandler(c *gin.Context) {
	url := oauthConfig.AuthCodeURL("random-state-string")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func redirectHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"redirected": "true"})
}

func main() {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	fmt.Println("clientId", clientId)
	fmt.Println("clientSecret", clientSecret)

	oauthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  os.Getenv("http://localhost:8080/auth/callback"),
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	gine_engine := gin.Default()
	gine_engine.LoadHTMLGlob("templates/*")
	gine_engine.GET("/", homeHandler)
	gine_engine.GET("/auth/login", loginHandler)
	gine_engine.GET("/auth/callback", redirectHandler)

	gine_engine.Run()
}
