package usecase

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/repository"
	"github.com/kizoukun/codingtest/web"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo *repository.UserRepository
}

func NewAuthUsecase() *AuthUsecase {
	return &AuthUsecase{
		userRepo: repository.NewUserRepository(),
	}
}

func (uc AuthUsecase) AuthLoginHandler(req web.LoginRequest, response *web.ResponseHttp) {
	// Handle user login

	user, err := uc.userRepo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Invalid email or password"
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Invalid email or password"
		return
	}

	// generate token
	token, err := generateToken(*user, 24*time.Hour)

	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to generate token"
		return
	}

	// no need to implement refresh token for this test
	response.StatusCode = http.StatusOK
	response.Data = web.LoginResponseData{
		Token: token,
	}
	response.Message = "Login successful"
	response.Success = true
}

func (uc AuthUsecase) AuthRegisterHandler(req web.RegisterRequest, response *web.ResponseHttp) {

	// get user by email
	existingUser, err := uc.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Email already in use"
		return
	}

	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to hash password"
		return
	}

	err = uc.userRepo.CreateUser(entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(bcryptedPassword),
		CreatedAt: time.Now(),
	})
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return
	}

	// future send an email verification

	response.StatusCode = http.StatusCreated
	response.Message = "User registered successfully"
	response.Success = true
}

func generateToken(user entity.User, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := web.JwtClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    "example-app",
			Subject:   user.Email,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_PRIVATE_KEY")))
}
