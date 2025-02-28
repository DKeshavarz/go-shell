package shell

import "systemgroup.net/bootcamp/go/v1/shell/internal/models"

type Shell struct {
	History     []string
	CurrentUser *models.User
	Handlers    map[string]func(*Shell, []string) (string, error)
	status      bool
}
