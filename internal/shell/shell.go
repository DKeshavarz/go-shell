package shell

import (
	"bufio"
	"fmt"
	"os"

	"systemgroup.net/bootcamp/go/v1/shell/internal/models"
)

type Shell struct {
	History []string
	CurrentUser *models.User
}

func New()*Shell{
	return &Shell{}
}

func (s *Shell)Start(){
	contine := true
	for contine {
		s.Show()
		command := s.Read()
		fmt.Print(command)
	}
}

func (s *Shell)Read()(str string){
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func  (s *Shell)Show(){
	if s.CurrentUser == nil {
		fmt.Print("$ ")
		return
	}

	fmt.Printf("%s :$",s.CurrentUser.UserName)
}