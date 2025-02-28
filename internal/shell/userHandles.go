package shell

import (
	"systemgroup.net/bootcamp/go/v1/shell/internal/models"
	"systemgroup.net/bootcamp/go/v1/shell/internal/servise"
)

func login(s *Shell, args []string) (msg string, err error) {
	user, err := createUserFromArgs(args)
	if err != nil {
		return msg, err
	}
	
	user, err = servise.GetUser(user)
	if err != nil {
		return msg, err
	}
	
	s.CurrentUser = user
	return
}

func addUser(s *Shell, args []string) (msg string, err error) {
	user, err := createUserFromArgs(args)
	if err != nil {
		return msg, err
	}

	err = user.Validate()
	if err != nil {
		return msg, err
	}

	err = servise.CreateUser(user)
	if err != nil {
		return msg, err
	}

	s.CurrentUser = user
	return
}

func logout(s *Shell, args []string) (msg string, err error) {
	s.CurrentUser = nil
	return
}

// ---------------------- helper -----------------------
func createUserFromArgs(args []string) (*models.User, error) {
	if len(args) > 2 {
		return nil, tooManyArgumentERR
	}
	if len(args) <= 0 {
		return nil, tooFewArgumentERR
	}

	var userName, password string

	userName = args[0]
	if len(args) == 2 {
		password = args[1]
	}

	user := &models.User{
		Username: userName,
		Password: password,
	}

	return user, nil
}
