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
		s.show()
		command := s.read()
		tokens, _ := s.tokenizer(command)
		
		for _, val := range tokens {
			fmt.Println(val)
		}
	}
}

func (s *Shell)read()(str string){
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func (s *Shell)show(){
	if s.CurrentUser == nil {
		fmt.Print("$ ")
		return
	}

	fmt.Printf("%s :$",s.CurrentUser.UserName)
}

func (s *Shell)tokenizer(input string)(tokens []string,err error){

	input += " "
	var str []rune

	for _, val := range input {

		if val == ' ' && len(str) > 0{
			tokens = append(tokens, string(str))
			str = make([]rune, 0)
		}else {
			str = append(str, val)
		}
	}

	return tokens, nil
}