package web

import (
	"strconv"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/webook/internal/domain"
	"github.com/webook/internal/service"
)

const (
	emailRegexpPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordRegexpPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`
)

type UserHandler struct {
	emailRegexp    *regexp.Regexp
	passwordRegexp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		passwordRegexp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
		svc:            svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", u.Signup)
	server.POST("/users/signin", u.Login)
	server.GET("/users/profile/:id", u.Profile)
	server.POST("/users/edit", u.Edit)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	type SignUpRequest struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	isEmail, err := u.emailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !isEmail {
		ctx.JSON(400, gin.H{"error": "invalid email"})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(400, gin.H{"error": "password and confirm password do not match"})
		return
	}

	isPassword, err := u.passwordRegexp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !isPassword {
		ctx.JSON(400, gin.H{"error": "invalid password"})
		return
	}

	err = u.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:
		ctx.JSON(201, gin.H{})
	case service.ErrDuplicateEmail:
		ctx.JSON(400, gin.H{"error": "email already exists"})
	default:
		ctx.JSON(500, gin.H{"error": err.Error()})
	}

}

func (u *UserHandler) Login(ctx *gin.Context) {
	type SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req SignInRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := u.svc.Signin(ctx, req.Email, req.Password)

	switch err {
	case nil:
		session := sessions.Default(ctx)
		session.Set("userId", user.Id)
		session.Options(sessions.Options{
			MaxAge: 900,
		})

		if err := session.Save(); err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{
			"id":    user.Id,
			"email": user.Email,
		})
	case service.ErrInvalidUserOrPassword:
		ctx.JSON(400, gin.H{"error": "invalid user or password"})
	default:
		ctx.JSON(500, gin.H{"error": err.Error()})
	}

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid user ID"})
		return
	}
	user, err := u.svc.FindById(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"id":    user.Id,
		"email": user.Email,
	})
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditRequest struct {
		Id       int64  `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req EditRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := u.svc.FindById(ctx, req.Id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = u.svc.Update(ctx, domain.User{
		Id:       user.Id,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"id":    user.Id,
		"email": user.Email,
	})
}
