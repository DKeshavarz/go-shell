package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"systemgroup.net/bootcamp/go/v1/shell/internal/servise"
)

func New() *Shell {
	shell := &Shell{
		Handlers: make(map[string]func(*Shell, []string) (string, error)),
		status: true,
	}

	shell.register("cd", cd)
	shell.register("echo", echo)
	shell.register("type", type_)
	shell.register("cat", cat)
	shell.register("pwd", pwd)
	shell.register("exit", exit)
	shell.register("adduser",addUser)
	shell.register("login",login)
	shell.register("logout",logout)
	shell.register("history",history)
	return shell
}

func (s *Shell) Start() {
	for s.status {
		s.show()
		command := s.read()
		tokens, _ := s.tokenizer(command)
		s.excute(tokens)
	}
}

func (s *Shell) read() (str string) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func (s *Shell) show() {
	if s.CurrentUser == nil {
		fmt.Print("$ ")
		return
	}

	fmt.Printf("%s :$", s.CurrentUser.Username)
}

func (s *Shell) tokenizer(input string) (tokens []string, err error) {

	input += " "
	var str []rune

	for _, val := range input {

		if val == ' ' && len(str) > 0 {
			tokens = append(tokens, string(str))
			str = make([]rune, 0)
		} else if val != ' ' {
			str = append(str, val)
		}
	}

	return tokens, nil
}

func (s *Shell) register(handleName string, handle func(*Shell, []string) (string, error)) {
	s.Handlers[handleName] = handle
}

func (s *Shell) excute(args []string) {
	if len(args) <= 0 {
		return
	}
	fmt.Println("log : ", s.historyLogger(args[0]))

	cmd, ok := s.Handlers[args[0]]

	var msg string
	var err error
	if ok {
		msg, err = cmd(s, args[1:])
	} else {
		msg, err = s.systemCommand(args)
	}
	fmt.Println("msg:", msg)
	fmt.Println("err:", err)
}

func (s *Shell) systemCommand(args []string) (string, error) {

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

func (s *Shell) historyLogger(command string)(error){
	if s.CurrentUser != nil {
		return servise.AddCommandHistory(s.CurrentUser.Username, command)
	}

	s.History = append(s.History, command)
	return nil
}
