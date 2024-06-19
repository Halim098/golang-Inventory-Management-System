package Model

import (
	"errors"
	"ims/Database"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Username   string `json:"username" binding:"required" gorm:"unique"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email" binding:"required" gorm:"unique"`
	Role       string `json:"role" binding:"required"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	Activities []Activity `gorm:"foreignKey:UserID"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) BeforeSave(*gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Save() (*User, error) {
	err := u.BeforeSave(Database.Database)
	if err != nil {
		return &User{}, err
	}

	result, _ := FindUserByUsername(u.Username)
	if result.Username == u.Username {
		return &User{}, errors.New("username already exists")
	}

	result, _ = FindUserByEmail(u.Email)
	if result.Email == u.Email {
		return &User{}, errors.New("email already exists")
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	err = Database.Database.Exec("INSERT INTO users (username, password, email, role, created_at, Updated_at) VALUES (?, ?, ?, ?, ?, ?)", u.Username, u.Password, u.Email, u.Role, u.CreatedAt, u.UpdatedAt).Error
	if err != nil {
		return &User{}, errors.New("failed to create user")
	}
	return u, nil
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func FindUserByEmail(email string) (User, error) {
	var user User
	err := Database.Database.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user)
	if err.Error != nil {
		return User{}, err.Error
	}

	if err.RowsAffected == 0 {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := Database.Database.Raw("SELECT * FROM users WHERE username = ?", username).Scan(&user)
	if err.Error != nil {
		return User{}, err.Error
	}

	if err.RowsAffected == 0 {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func GetUserByID(id uint) (User, error) {
	var user User
	err := Database.Database.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&user)
	if err.Error != nil {
		return User{}, err.Error
	}

	if err.RowsAffected == 0 {
		return User{}, errors.New("user not found")
	}

	return user, nil
}
