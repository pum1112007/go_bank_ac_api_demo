package userapi

import (
	"database/sql"
	"go_bank_ac_api_demo/user"
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
	AddBankAc(bkAc *user.BankAccount) error
	GetAllUserBkAc(id int) ([]user.BankAccount, error)
	// RemoveBkAc(u *user.User) error
	// Withdraw(u *user.User) error
	// Deposit(u *user.User) error
	// Transfers(u *user.User) error
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

//BankAccount API
func (h *Handler) addBankAc(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var bkAc user.BankAccount
	err = c.ShouldBindJSON(&bkAc)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	bkAc.UserID = id
	err = h.userService.AddBankAc(&bkAc)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, bkAc)
}

func (h *Handler) getAllUserBkAc(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	bkAc, err := h.userService.GetAllUserBkAc(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, bkAc)
}

// func (h *Handler) removeBkAc(c *gin.Context) {
// 	err := h.userService.RemoveBkAc()
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, users)
// }
// func (h *Handler) withdraw(c *gin.Context) {
// 	err := h.userService.Withdraw()
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, users)
// }
// func (h *Handler) deposit(c *gin.Context) {
// 	err := h.userService.Deposit()
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, users)
// }
// func (h *Handler) transfers(c *gin.Context) {
// 	users, err := h.userService.Transfers()
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, users)
// }
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
	r.DELETE("/users/:id", h.deleteUser)
	r.POST("/users/:id/bankAccounts", h.addBankAc)
	r.GET("/users/:id/bankAccounts", h.getAllUserBkAc)
	// r.DELETE("/bankAccount/:id", h.removeBkAc)
	// r.PUT("/bankAccount/:id/withdraw", h.withdraw)
	// r.PUT("/bankAccount/:id/deposit", h.deposit)
	// r.POST("/transfers", h.transfers)

	return r.Run(addr)
}
