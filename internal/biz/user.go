package biz

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"tkratos-realworld/internal/conf"

	v1 "tkratos-realworld/api/realworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"tkratos-realworld/internal/pkg/middleware/auth"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type User struct {
	ID           uint
	Email        string
	Username     string
	Bio          string
	Image        string
	PasswordHash string
}

type Profile struct {
	ID        uint
	Username  string
	Bio       string
	Image     string
	Following bool
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
}

type ProfileRepo interface {
	GetProfile(ctx context.Context, username string) (*Profile, error)
	FollowUser(ctx context.Context, currentUserID uint, followingID uint) error
	UnfollowUser(ctx context.Context, currentUserID uint, followingID uint) error
	GetUserFollowingStatus(ctx context.Context, currentUserID uint, userIDs []uint) (following []bool, err error)
}

type UserUsecase struct {
	ur   UserRepo
	pr   ProfileRepo
	jwtc *conf.JWT

	log *log.Helper
}

func NewUserUsecase(
	ur UserRepo,
	pr ProfileRepo,
	logger log.Logger,
	jwtc *conf.JWT) *UserUsecase {
	return &UserUsecase{ur: ur, pr: pr, jwtc: jwtc, log: log.NewHelper(logger)}
}

type UserLogin struct {
	Email    string
	Username string
	Token    string
	Bio      string
	Image    string
}

func (uc *UserUsecase) Login(ctx context.Context, email, password string) (*UserLogin, error) {
	if len(email) == 0 {
		return nil, errors.New(422, "email", "cannot be empty")
	}
	u, err := uc.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !verifyPassword(u.PasswordHash, password) {
		return nil, errors.Unauthorized("user", "login failed")
	}

	return &UserLogin{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Token:    uc.generateToken(u.ID),
	}, nil
}

func (uc *UserUsecase) Register(ctx context.Context, username, email, password string) (*UserLogin, error) {
	u := &User{
		Email:        email,
		Username:     username,
		PasswordHash: hashPassword(password),
	}
	if err := uc.ur.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return &UserLogin{
		Email:    email,
		Username: username,
		Token:    uc.generateToken(u.ID),
	}, nil
}

func (uc *UserUsecase) GetCurrentUser(ctx context.Context) (*User, error) {
	cu := auth.FromContext(ctx)
	u, err := uc.ur.GetUserByID(ctx, cu.UserID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

type UserUpdate struct {
	Email    string
	Username string
	Password string
	Bio      string
	Image    string
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, uu *UserUpdate) (*UserLogin, error) {
	cu := auth.FromContext(ctx)
	u, err := uc.ur.GetUserByID(ctx, cu.UserID)
	if err != nil {
		return nil, err
	}
	u.Email = uu.Email
	u.Image = uu.Image
	u.PasswordHash = hashPassword(uu.Password)
	u.Bio = uu.Bio
	u, err = uc.ur.UpdateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return &UserLogin{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Token:    uc.generateToken(u.ID),
	}, nil
}

func hashPassword(pwd string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func verifyPassword(hashed, input string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input)); err != nil {
		return false
	}
	return true
}

func (uc *UserUsecase) generateToken(userID uint) string {
	return auth.GenerateToken(uc.jwtc.Secret, userID)
}
