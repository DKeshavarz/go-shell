package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"systemgroup.net/bootcamp/go/v1/shell/internal/commands"
)

func New()*Shell{
	shell := &Shell{
		Handlers: make(map[string]func([]string) (string, error)),
	}

	shell.register("cd",commands.Cd)
	shell.register("echo",commands.Echo)
	shell.register("type",commands.Type)
	shell.register("cat",commands.Cat)
	shell.register("pwd",commands.Pwd)
	shell.register("exit",commands.Exit)

	return shell
}

func (s *Shell)Start(){
	contine := true
	for contine {
		s.show()
		command := s.read()
		tokens, _ := s.tokenizer(command)
		s.excute(tokens)
	}
}

func (s *Shell)read()(str string){
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
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
		}else if val != ' '{
			str = append(str, val)
		}
	}

	return tokens, nil
}

func (s *Shell) register(handleName string, handle func([]string) (string, error)){
	s.Handlers[handleName] = handle
}

func (s *Shell) excute(args []string){
	if len(args) <= 0 {
		return
	}

	cmd, ok := s.Handlers[args[0]]

	var msg string
	var err error
	if ok {
		msg, err = cmd(args[1:])
	}else{
		msg, err = s.systemCommand(args)
	}
	fmt.Println("msg:", msg)
	fmt.Println("err:", err)
}

func (s *Shell) systemCommand(args []string)(string, error){

	command := exec.Command(args[0], args[1:]...)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return stderr.String(), err
	}

	return stdout.String(), nil
}