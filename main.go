package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nyamka11/backEnd/controllers"
	"github.com/nyamka11/backEnd/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	ctx := context.Background()
	route := gin.Default()
	// route.GET("/", controllers.Users)

	store := cookie.NewStore([]byte("secret"))
	route.Use(sessions.Sessions("mysession", store))

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?parseTime=true")

	if err != nil {
		log.Fatalf("Cannot connect database: %v", err)
	}
	boil.SetDB(db)

	route.LoadHTMLGlob("view/*")
	usersGroup := route.Group("users")
	{
		usersGroup.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", map[string]interface{}{
				"title": "testTitle",
			})
		})
		usersGroup.GET("/list", func(ctx *gin.Context) {
			users, err := models.Users().AllG(ctx)
			if err != nil {
				fmt.Errorf("Get todo error: %v", err)
			}

			ctx.HTML(http.StatusOK, "list.html", map[string]interface{}{
				"users": users,
			})
		})
		usersGroup.POST("/list", func(c *gin.Context) {
			controllers.Register(ctx, c)
		})
		usersGroup.GET("/delete/:user_id", func(c *gin.Context) {
			controllers.Delete(ctx, c)
		})
		usersGroup.GET("/edit/:user_id", func(c *gin.Context) {
			controllers.Update(ctx, c)
		})
		usersGroup.POST("/edit/:user_id", func(c *gin.Context) {
			controllers.UpdateExc(ctx, c)
		})
	}

	authGroup := route.Group("auth")
	{
		authGroup.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", map[string]interface{}{})
		})
		authGroup.POST("/login", func(c *gin.Context) {
			controllers.Login(ctx, c)
		})
		authGroup.GET("/logout", func(c *gin.Context) {
			controllers.Logout(ctx, c)
		})
	}

	route.Run() // listen and serve on 0.0.0.0:8080
}
