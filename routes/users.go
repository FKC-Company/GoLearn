package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nyamka11/backEnd/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func Register(ctx context.Context, c *gin.Context) {
	var b models.User
	if err := c.Bind(&b); err != nil {
		fmt.Errorf("%#v", err)
	}

	err := b.InsertG(ctx, boil.Infer())

	users, err := models.Users().AllG(ctx)
	if err != nil {
		fmt.Errorf("Get todo error: %v", err)
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

	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"users": users,
	})
}

func parseInt(arg string) int {
	i1, err := strconv.Atoi(arg)
	if err == nil {
		fmt.Println(i1)
	}

	return i1
}
