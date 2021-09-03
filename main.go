package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nyamka11/backEnd/models"
	"github.com/nyamka11/backEnd/routes"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	ctx := context.Background()
	route := gin.Default()
	// route.GET("/", controllers.Users)

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
			routes.Register(ctx, c)
		})

		usersGroup.GET("/delete/:user_id", func(c *gin.Context) {
			routes.Delete(ctx, c)
		})

		usersGroup.GET("/edit/:user_id", func(c *gin.Context) {
			routes.Update(ctx, c)
		})
	}

	route.Run() // listen and serve on 0.0.0.0:8080

}
