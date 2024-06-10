// internal/api/handlers/auth_handler.go

package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"marketplace-api/internal/models"
	"marketplace-api/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService        *services.UserService
	distributorService *services.DistributorService
	storeService       *services.StoreService
	jwtSecret          string
	log                *logrus.Logger
}

type CustomClaims struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Name        string `json:"name"`
	CompanyName string `json:"company_name"`
	Details     string `json:"details"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	BIN         string `json:"bin"`
	jwt.StandardClaims
}

func NewAuthHandler(userService *services.UserService, distributorService *services.DistributorService, storeService *services.StoreService, jwtSecret string, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{userService: userService, distributorService: distributorService, storeService: storeService, jwtSecret: jwtSecret, log: log}
}

// Register godoc
// @Summary      Register user
// @Description  register user by role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body  models.RegisterInput  true "query params"
// @Success      200  {message}  User registered successfully
// @Failure      400  {error}  Bad request
// @Failure      404  {error}  Not found
// @Failure      500  {error}  Internal server error
// @Router       /auth/register [post]
func (ah *AuthHandler) Register(c *gin.Context) {
	var registerInput models.RegisterInput
	if err := c.BindJSON(&registerInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	err := ah.userService.RegisterUser(&registerInput)
	if err != nil {
		ah.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary      Login
// @Description  Logs in a user and returns a token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.LoginCredentials true "Login credentials"
// @Success      200  string  token
// @Failure      400  string  Bad request
// @Failure      404  string  Not found
// @Failure      500  string  Internal server error
// @Router       /auth/login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	var credentials models.LoginCredentials
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate user credentials
	user, err := ah.userService.ValidateCredentials(credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if credentials.Role != "admin" {
		status := ah.userService.GetStatusById(user.ID)
		if status == false {
			c.JSON(http.StatusForbidden, gin.H{"error": "account is not activated"})
			return
		}
	}

	// Generate JWT token
	var token string
	token, err = generateJWTToken(user, ah.jwtSecret)
	if err != nil {
		ah.log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
func (ah *AuthHandler) LogOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (ah *AuthHandler) GetStoreByID(c *gin.Context) {
	storeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || storeId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	store, err := ah.distributorService.GetStoreByID(storeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": store})
}

func (ah *AuthHandler) GetDistributorByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	distributor, err := ah.distributorService.GetDistributorByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": distributor})
}
func (ah *AuthHandler) ChangePassword(c *gin.Context) {
	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	userID := c.GetInt64("user_id")
	err := ah.userService.UpdatePassword(userID, input.CurrentPassword, input.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

// ChangeEmail godoc
// @Summary      Change users email
// @Description  Returns distributors profile
// @Tags         user
// @Security     BearerToken
// @Accept       json
// @Produce      json
// @Param        request body InputChangeEmail true "sdfg"
// @Success      200  string  email updated successfully
// @Failure      400  {string}  Bad request
// @Failure      404  {string}  Not found
// @Failure      500  {string}  Internal server error
// @Router       /distributor/profile [get]
func (ah *AuthHandler) ChangeEmail(c *gin.Context) {
	var input InputChangeEmail
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	userID := c.GetInt64("user_id")
	err := ah.userService.UpdateEmail(userID, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "email updated successfully"})
}
func (ah *AuthHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetInt64("user_id")
	err := ah.userService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func generateJWTToken(user *models.User, jwtSecret string) (string, error) {
	// Create a new JWT token
	customClaims := CustomClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	customClaims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Generate encoded token and return it as string
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// InputChangeEmail model
type InputChangeEmail struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
