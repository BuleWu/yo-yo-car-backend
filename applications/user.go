package applications

import (
	"github.com/google/uuid"
	"net/http"
	"zavrsni/yo-yo-car/models"
	"zavrsni/yo-yo-car/repositories"
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

func (a *User) GetUser(userId string) (interface{}, Exception) {
	user, err := a.userRepository.Get(userId)

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
}

/*CreateUser is the User controller method that handles the POST request*/
func (a *User) CreateUser(request *CreateUserRequest) (interface{}, Exception) {
	user := models.NewUser(uuid.NewString())
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Email = request.Email
	user.Password = request.Password
	user, err := a.userRepository.Persist(user)
	if err != nil {
		return nil, NewApplicationException(http.StatusInternalServerError, err)
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
