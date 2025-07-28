package applications

import (
	"errors"
	"mime/multipart"
	"net/http"
	"zavrsni/yo-yo-car/core/utils"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/repositories"
	"zavrsni/yo-yo-car/shared/storage"
)

const (
	defaultProfilePictureUrl = "https://storage.googleapis.com/yoyo-car-no2.firebasestorage.app/profile_pictures/default-profile-picture.png"
)

func NewUserApplication(
	userRepository repositories.UserRepository,
) *User {
	return &User{
		userRepository: userRepository,
	}
}

type User struct {
	Application
	userRepository repositories.UserRepository
}

func (a *User) GetUserById(userId string) (*models.User, Exception) {
	user, err := a.userRepository.GetById(userId)

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return user, nil
}

func (a *User) GetUserByEmail(email string) (*models.User, Exception) {
	user, err := a.userRepository.GetByEmail(email)

	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}

	return user, nil
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Provider  string `json:"provider"`
}

/*CreateUser is the User controller method that handles the POST request*/
func (a *User) CreateUser(request *CreateUserRequest) (*models.User, Exception) {
	user, err := a.userRepository.Persist(models.NewUser(request.FirstName, request.LastName, request.Email, request.Password, request.Provider, defaultProfilePictureUrl))
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
	}
	return user, nil
}

type UpdateUserRequest struct {
	UserID    string `json:"-"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Vehicle   string `json:"vehicle"`
}

func (a *User) UpdateUser(request *UpdateUserRequest) (*models.User, Exception) {
	user, err := a.userRepository.GetById(request.UserID)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}

	if request.LastName != "" {
		user.LastName = request.LastName
	}

	if request.Email != "" {
		user, _ = a.userRepository.GetByEmail(request.Email)
		if user != nil {
			return nil, NewApplicationException(http.StatusNotFound, errors.New("an account with this email already exists"))
		}
		user.Email = request.Email
	}

	if request.Vehicle != "" {
		user.Vehicle = request.Vehicle
	}

	user, err = a.userRepository.Update(user)
	if err != nil {
		return nil, NewApplicationException(http.StatusNotFound, err)
	}

	return user, nil
}

func (a *User) DeleteUser(userId string) Exception {
	err := a.userRepository.Delete(userId)

	if err != nil {
		return NewApplicationException(http.StatusInternalServerError, err)
	}

	return nil
}

func (u *User) UploadProfilePicture(userID string, file multipart.File, header *multipart.FileHeader) (string, Exception) {
	url, err := storage.UploadProfilePicture(file, header, userID)
	if err != nil {
		return "", NewApplicationException(http.StatusInternalServerError, err)
	}

	if err := u.userRepository.UpdateProfilePicture(userID, url); err != nil {
		return "", NewApplicationException(http.StatusInternalServerError, err)
	}

	return url, nil
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

func (a *User) ChangePassword(userID string, request *ChangePasswordRequest) Exception {
	user, err := a.userRepository.GetById(userID)
	if err != nil {
		return NewApplicationException(http.StatusNotFound, err)
	}

	if !utils.VerifyPassword(user.Password, request.CurrentPassword) {
		return NewApplicationException(http.StatusInternalServerError, errors.New("current password is incorrect"))
	}

	hashedNewPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return NewApplicationException(http.StatusInternalServerError, errors.New("could not hash new password"))
	}

	user.Password = hashedNewPassword

	if _, err = a.userRepository.Update(user); err != nil {
		return NewApplicationException(http.StatusInternalServerError, errors.New("failed to update password"))
	}

	return nil
}
