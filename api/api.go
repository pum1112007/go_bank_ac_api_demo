package api

import (
	"GO_BANK_AC_API_DEMO/user"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	FindByID(id int) (*user.User, error)
	All() ([]user.User, error)
	Insert(u *user.User) error
	Update(u *user.User) error
	Delete(u *user.User) error
}

type Handler struct {
	userService UserService
}

func (h *Handler) getUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := h.userService.FindByID(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) allUser(c *gin.Context) {
	users, err := h.userService.All()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	u, err := h.userService.FindByID(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var update struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
	}

	err = c.ShouldBindJSON(&update)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if update.FirstName != nil {
		u.FirstName = *update.FirstName
	}
	if update.LastName != nil {
		u.LastName = *update.LastName
	}
	if update.Email != nil {
		u.Email = *update.Email
	}

	err = h.userService.Update(u)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = h.userService.Delete(&user.User{
		ID: id,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (h *Handler) createUser(c *gin.Context) {
	var u user.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = h.userService.Insert(&u)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, u)
}

func StartServer(addr string, db *sql.DB) error {
	r := gin.Default()
	h := &Handler{
		userService: &user.Service{
			DB: db,
		},
	}
	r.GET("/users", h.allUser)
	r.GET("/users/:id", h.getUser)
	r.POST("/users", h.createUser)
	r.PUT("/users/:id", h.updateUser)
	//r.POST("/users/:id/bankAccount", )
	r.DELETE("/users/:id", h.deleteUser)
	// r.GET("/users/:id/bankAccounts",h.)
	// r.DELETE("/bankAccount/:id")
	// r.PUT("/bankAccount/:id/withdraw")
	// r.PUT("/bankAccount/:id/deposit")
	// r.POST("/transfers")

	return r.Run(addr)
}
