package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func Echo(arg []string)(msg string, err error){
	for _,val := range arg{

		msg += val + " "
	}
	return
}

func Type(args []string)(msg string, err error){
	if len(args) == 0 {
		return msg,tooFewArgumentERR
	}
	if len(args) >= 2{
		return msg,tooManyArgumentERR
	}

	cmd := args[0]
	builtins := map[string]bool{
		"exit": true, "echo": true, "type": true,
		"cd": true, "pwd": true, "cat": true,
		"adduser": true, "login": true, "logout": true,
		"history": true,
	}

	if builtins[cmd] {
		msg = fmt.Sprintf("%s is a shell builtin",cmd)
		return
	}

	path, err := exec.LookPath(cmd)
	if err != nil {
		msg = fmt.Sprintf("%s: command not found", cmd)
	} else {
		msg = fmt.Sprintf("%s is %s", cmd, path)
	}
	return
}

func Cd(args []string)(msg string, err error){
	if len(args) == 0 {
		return msg,tooFewArgumentERR
	}
	if len(args) >= 2{
		return msg,tooManyArgumentERR
	}

	dir := args[0]

	if err = os.Chdir(dir); err != nil {
		return
	}

	return
}

func Pwd(args []string)(msg string, err error){
	msg, err = os.Getwd()
	return
}

func Cat(args []string)(string, error){
	var msg string
	
	if len(args) == 0 {
		return msg,tooFewArgumentERR
	}
	
	for _, file := range args {
		data, err := os.ReadFile(file)
		if err != nil {
			return  msg, err
		}
		msg += string(data)
	}

	return msg, nil
}

func Exit(args []string) (msg string, err error) {
	code := 0
	if len(args) <= 0{
		msg = fmt.Sprintf("exit status %d",code)
		return
	} 
	
	if len(args) ==  1{
		if code, err = strconv.Atoi(args[0]); err != nil {
			msg = fmt.Sprintf("exit status %d",code)
			return
		}
	}
	
	return msg, tooManyArgumentERR
}


//---------------------- helper ---------------------------
func findExecutableInPath(cmd string) (string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", fmt.Errorf("PATH environment variable not set")
	}
	
	pathDirs := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range pathDirs {
		if dir == "" {
			continue
		}
		
		fullPath := filepath.Join(dir, cmd)

		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}

		if info.Mode().IsRegular() && (info.Mode().Perm()&0111 != 0) {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("command not found")
}