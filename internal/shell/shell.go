package shell

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"systemgroup.net/bootcamp/go/v1/shell/internal/models"
	"systemgroup.net/bootcamp/go/v1/shell/internal/servise"
)

func New() *Shell {
	shell := &Shell{
		Handlers: make(map[string]func(*Shell, []string) (string, error)),
		status:   true,
	}

	shell.register("cd", cd)
	shell.register("echo", echo)
	shell.register("type", type_)
	shell.register("cat", cat)
	shell.register("pwd", pwd)
	shell.register("exit", exit)
	shell.register("adduser", addUser)
	shell.register("login", login)
	shell.register("logout", logout)
	shell.register("history", history)
	return shell
}

func (s *Shell) Start() {
	for s.status {
		s.prompt()
		command := s.read()
		tokens, red , _:= s.tokenizer(command)
		msg, err := s.excute(tokens)
		s.show(red, msg, err)
	}
}

func (s *Shell) read() (str string) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func (s *Shell) prompt() {
	if s.CurrentUser == nil {
		fmt.Print("$ ")
		return
	}

	fmt.Printf("%s :$", s.CurrentUser.Username)
}

func (s *Shell) tokenizer(in string) (tokens []string, redirect []string, err error) {

	in += " "
	var str string
	var tmpStr string
	seenVar := false
	inQuote := false
	inDQote := false
	for i := 0 ; i < len(in) ; i++ {
		if in[i] == '\'' {
			if inDQote {

				str += string(in[i])
			}else if inQuote {

				inQuote = false
			}else {

				inQuote = true
			}

		}else if in[i] == '"' {
			if inQuote {

				str += string(in[i])
			}else if inDQote{

				inDQote = false
			}else {

				inDQote = true
			}
		}else if in[i] == '\\' {
			if inQuote {

				str += string(in[i])
			}else if inDQote && i+1 < len(in) && isScalbe(string(in[i+1])){
				i++
				str += string(in[i])
			} else if inDQote {

				str += string(in[i])
			}
			// in normal mode we don't add \ to string like unix terminal
		}else if in[i] == '$' {
			if inQuote {

				str += string(in[i])
			}else {
				seenVar = true
				tmpStr = str
				str = ""
			}
		}else if in[i] == ' '{
			if inQuote {

				str += string(in[i])
			}else if inDQote{
				if seenVar{
					str = tmpStr + getEnvVar(str)
					seenVar = false
					tmpStr = ""
				}
				str += string(in[i])
			}else if len(str) > 0{
				if seenVar{
					str = tmpStr + getEnvVar(str)
					seenVar = false
					tmpStr = ""
				}
				tokens = append(tokens, str)
				str = ""
				tmpStr = ""
			}
		}else{
			str += string(in[i])
		}
	}

	if inDQote || inQuote || seenVar {
		return nil, nil, fmt.Errorf("wrong input")
	}

	var last string
	if len(tokens) >= 3 {
		last = tokens[len(tokens)-2] 
	}

	if last == "1>>" || last == "1>" || last == ">" || last == ">>" || last == "2>>" || last == "2>" {
		return tokens[:len(tokens)-2], tokens[len(tokens)-2 :], nil
	}
	return tokens, nil, nil
}

func (s *Shell) register(handleName string, handle func(*Shell, []string) (string, error)) {
	s.Handlers[handleName] = handle
}

func (s *Shell) excute(args []string) (msg string,err error) {
	if len(args) <= 0 {
		return 
	}
	s.historyLogger(args[0])
	cmd, ok := s.Handlers[args[0]]

	if ok {
		msg, err = cmd(s, args[1:])
	} else {
		msg, err = s.systemCommand(args)
	}
	
	if err != nil && len(err.Error()) >= 5 && err.Error()[:5] ==  "exec:"{
		
		err = errors.New(args[0]+": command not found")
	}
	return msg, err
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

func (s *Shell) historyLogger(command string) error {
	if s.CurrentUser != nil {
		return servise.AddCommandHistory(s.CurrentUser.Username, command)
	}

	for i := range s.History {
		if s.History[i].Command == command {
			s.History[i].Count++
			s.History[i].CreatedAt = time.Now()
			return nil
		}
	}

	s.History = append(s.History, models.CommandHistory{
		Command:   command,
		Count:     1,
		CreatedAt: time.Now(),
	})
	return nil
}

func (s *Shell) show(tokens []string, msg string, err error) {
	write_msg, write_err := false, false
	if tokens != nil ||len(tokens) == 2 {
		if tokens[0] == ">" || tokens[0] == "1>" {
			servise.WriteToFile(tokens[1], msg, servise.Overwrites)
			write_msg = true
		}else if tokens[0] == ">>" || tokens[0] == "1>>"{
			servise.WriteToFile(tokens[1], msg, servise.Append)
			write_msg = true
		}else if tokens[0] == "2>"{
			servise.WriteToFile(tokens[1], err.Error(), servise.Overwrites)
			write_err = true
		}else if tokens[0] == "2>>"{
			write_err = true
			servise.WriteToFile(tokens[1], err.Error(), servise.Append)
		}
	}

	if err != nil && !write_err{
		fmt.Println(err.Error())
		return
	}
	if !write_msg && msg != ""{
		fmt.Println(msg)
	}
}


func isScalbe(input string)bool {
	char := []string{ "$", "`", "\"", "\\" }
	for _,val := range char {
		if val == input {
			return true
		}
	}
	return false
}

func getEnvVar(varName string) string {
    return os.Getenv(varName)
}
