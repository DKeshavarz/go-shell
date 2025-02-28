package servise

import (
    "fmt"
    "os"
)

type Mode int8
const (
	Append Mode = iota
	Overwrites
)
func WriteToFile(filename, message string, mode Mode) error {
    var file *os.File
    var err error

    if mode == Append {
        file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    } else {
        file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
    }

    if err != nil {
        return fmt.Errorf("could not open file: %w", err)
    }
    defer file.Close() 

    _, err = file.WriteString(message + "\n")
    if err != nil {
        return fmt.Errorf("could not write to file: %w", err)
    }

    return nil 
}
