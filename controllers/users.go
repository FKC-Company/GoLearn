package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nyamka11/backEnd/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx context.Context, c *gin.Context) {
	var b models.User
	if err := c.Bind(&b); err != nil {
		fmt.Errorf("%#v", err)
	}

	hashidPassword := passwordHash([]byte(b.Password))
	b.Password = hashidPassword

	err := b.InsertG(ctx, boil.Infer())

	users, err := models.Users().AllG(ctx)
	if err != nil {
		fmt.Errorf("Get user error: %v", err)
	}
	c.HTML(http.StatusOK, "list.html", map[string]interface{}{
		"users": users,
	})
}

func Delete(ctx context.Context, c *gin.Context) {
	userId := c.Param("user_id")

	var u models.User
	if err := c.Bind(&u); err != nil {
		fmt.Errorf("%#v", err)
	}

	i1, err := strconv.Atoi(userId)
	if err == nil {
		fmt.Println(i1)
	}

	u.UserID = i1
	af, err := u.DeleteG(ctx)
	fmt.Println("Get user affected:", af)

	if err != nil {
		fmt.Errorf("Get user error: %v", err)
	}

	users, err := models.Users().AllG(ctx)
	if err != nil {
		fmt.Errorf("Get user error: %v", err)
	}

	c.HTML(http.StatusOK, "list.html", map[string]interface{}{
		"users": users,
	})
}

func Update(ctx context.Context, c *gin.Context) {

	var u models.User
	if err := c.Bind(&u); err != nil {
		fmt.Errorf("%#v", err)
	}

	userId := c.Param("user_id")
	u.UserID = parseInt(userId)

	users, err := models.Users(Where("user_id = ?", u.UserID)).AllG(ctx)

	if err != nil {
		fmt.Errorf("Get user error: %v", err)
	}

	c.HTML(http.StatusOK, "edit.html", map[string]interface{}{
		"users": users,
	})
}

func UpdateExc(ctx context.Context, c *gin.Context) {
	userId := parseInt(c.Param("user_id"))

	var u models.User
	if err := c.Bind(&u); err != nil {
		fmt.Errorf("%#v", err)
	}

	u.UserID = userId

	af, err := u.UpdateG(ctx, boil.Whitelist("username", "email", "password"))
	println("Get users error: %v", af)
	if err != nil {
		fmt.Errorf("Get users error: %v", err)
	}

	users, err := models.Users().AllG(ctx)
	if err != nil {
		fmt.Errorf("Get users error: %v", err)
	}
	c.HTML(http.StatusOK, "list.html", map[string]interface{}{
		"users": users,
	})

}

func Login(ctx context.Context, c *gin.Context) {
	// session := sessions.Default(c)

	var u models.User
	if err := c.Bind(&u); err != nil {
		fmt.Errorf("%#v", err)
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Println("------------------")
	fmt.Println(username)
	fmt.Println(password)
	fmt.Println("------------------")

	// userId := c.Param("user_id")
	// u.UserID = parseInt(userId)

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	users, err := models.Users(Where("username=?", username)).OneG(ctx)
	if err != nil {
		fmt.Errorf("Get user error: %v", err)
	}

	// // Check for username and password match, usually from a database
	// if username != users.Username || password != users.Password {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	// 	return
	// }

	err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	} else {
		session := sessions.Default(c)
		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(http.StatusOK, gin.H{"hello": session.Get("hello")})
		// c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
	}

	// user, user_cookie_err := c.Cookie("user")
	// if user_cookie_err != nil {
	// 	user = "Guest"
	// }

	// c.SetCookie("user", username, 3600, "/", "localhost", false, true)
	// c.JSON(200, gin.H{"message": "Hello " + user})
}

func parseInt(arg string) int {
	i1, err := strconv.Atoi(arg)
	if err == nil {
		fmt.Println(i1)
	}

	return i1
}

func passwordHash(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
