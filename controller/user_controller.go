package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/model"
	"github.com/wahyujatirestu/payshare/service"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) Register(ctx *gin.Context)  {
	var req struct {
		Name			string		`json:"name" binding:"required"`
		Username		string		`json:"username" binding:"required"`
		Email			string		`json:"email" binding:"required,email"`
		Phone			string		`json:"phone"`
		Password		string		`json:"password" binding:"required,min=8"`
		ConfirmPassword	string		`json:"confirmPassword" binding:"required"`
		Address			string		`json:"address"`
		Role			string		`json:"role" binding:"required,oneof=customer employee"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		Name:       req.Name,
		Username:   req.Username,
		Email:      req.Email,
		Phone:      req.Phone,
		Password: 	req.Password,
		Address:    req.Address,
		Role:       req.Role,
	}

	err := c.userService.Register(user, req.ConfirmPassword)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, gin.H{"message": "User registered successfully"})
}

func (c *UserController) GetAllUser(ctx *gin.Context)  {
	filters := make(map[string]interface{})

	if name := ctx.Query("name"); name != "" {
		filters["name"] = name
	}

	if email := ctx.Query("email"); email != "" {
		filters["email"] = email
	}

	if role := ctx.Query("role"); role != "" {
		filters["role"] = role
	}

	if username := ctx.Query("username"); username != "" {
		filters["username"] = username
	}

	user, err := c.userService.GetAllUser(filters)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200,gin.H{"data": user})
}


func (c *UserController) GetUserById(ctx *gin.Context)  {
	id := ctx.Param("id")
	user, err := c.userService.GetUserById(id)
	if err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"data": user})
}

func (c *UserController) UpdateUser(ctx *gin.Context)  {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	loggedInUser, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser := loggedInUser.(model.User)

	var req struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email" binding:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Role     string `json:"role" binding:"omitempty,oneof=customer employee"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if currentUser.ID.String() != id.String() && currentUser.Role != "employee" {
		ctx.JSON(403,gin.H{"error": "Forbidden"})
		return
	}

	if currentUser.Role != "employee" && req.Role != "" {
		ctx.JSON(403, gin.H{"error": "Forbidden: only employee can update role"})
		return
	}

	user := &model.User{
		ID:        id,
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		Address:   req.Address,
		Role: 	   req.Role,
	}

	err = c.userService.Update(user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User updated successfully"})
}

func (c *UserController) DeleteUser(ctx *gin.Context)  {
	id := ctx.Param("id")
	if err := c.userService.Delete(id); err != nil {
		ctx.JSON(404, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User deleted successfully"})
}
