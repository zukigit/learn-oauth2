package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	goGithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	oauthConfig      *oauth2.Config
	oauthStateString string = "random-state-string"
)

func homeHandler(c *gin.Context) {
	if c == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("c is nil")})
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"LoginURL": "/auth/login",
	})
}

func loginHandler(c *gin.Context) {
	url := oauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func redirectHandler(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid state: %s", state)})
		return
	}

	code := c.Query("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("oauthConfig.Exchange() failed, err: %s", err.Error())})
		return
	}

	c.SetCookie("oauth_token", token.AccessToken, 3600, "/", "rocky10", false, true)
	c.Redirect(http.StatusFound, "/profile")
}

func authMiddleware(c *gin.Context) {
	token, err := c.Cookie("oauth_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if token == "" {
		c.Redirect(http.StatusFound, "/")
		c.Abort()

		return
	}

	c.Set("access_token", token)
	c.Next()
}

func profileHandler(c *gin.Context) {
	token := c.GetString("access_token")

	// Create GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := goGithub.NewClient(tc)

	// Get authenticated user
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.GetID(),
		"login":      user.GetLogin(),
		"name":       user.GetName(),
		"email":      user.GetEmail(),
		"avatar_url": user.GetAvatarURL(),
		"company":    user.GetCompany(),
		"location":   user.GetLocation(),
		"bio":        user.GetBio(),
	})
}

func main() {
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("http://rocky10/:8080/auth/callback"),
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}

	gine_engine := gin.Default()
	gine_engine.LoadHTMLGlob("templates/*")
	gine_engine.GET("/", homeHandler)
	gine_engine.GET("/auth/login", loginHandler)
	gine_engine.GET("/auth/callback", redirectHandler)
	gine_engine.GET("/profile", authMiddleware, profileHandler)

	gine_engine.Run()
}
