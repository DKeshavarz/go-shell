package shell

import (
	"os"
	"path/filepath"
	"testing"
)


func TestEcho(t *testing.T) {
    shell := New()

    msg, err := echo(shell, []string{"hello  ali"})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    if msg != "hello  ali" {
        t.Errorf("Expected 'hello ', got: '%s'", msg)
    }

    msg, err = echo(shell, []string{"hello", "ali","reza"})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    if msg != "hello ali reza" {
        t.Errorf("Expected 'hello world ', got: '%s'", msg)
    }
}

func TestCd(t *testing.T) {
    shell := New()

    tempDir := t.TempDir()
    msg, err := cd(shell, []string{tempDir})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
	if msg != "" {
        t.Errorf("Expected msg to be empty, got: %v", err)
    }
    currentDir, _ := os.Getwd()
    if currentDir != tempDir {
        t.Errorf("Expected current directory to be '%s', got: '%s'", tempDir, currentDir)
    }

    msg, err = cd(shell, []string{"/soblskjglkjlvj"})
    if err == nil {
        t.Error("Expected error, got nil")
    }

    msg, err = cd(shell, []string{})
    if err == nil {
        t.Error("Expected error, got nil")
    }

    msg, err = cd(shell, []string{"arg1", "arg2"})
    if err == nil {
        t.Error("Expected error, got nil")
    }
}

func TestPwd(t *testing.T) {
    shell := New()

    expectedDir, _ := os.Getwd()
    msg, err := pwd(shell, []string{})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    if msg != expectedDir {
        t.Errorf("Expected '%s', got: '%s'", expectedDir, msg)
    }

	msg, err = pwd(shell, []string{"one", "two"})
	if err == nil {
        t.Error("Expected have error, got: nil")
    }

}

func TestCat(t *testing.T) {
    shell := New()

    tempFile := filepath.Join(t.TempDir(), "test.txt")
    os.WriteFile(tempFile, []byte("let me innnnnnnnnnnn"), 0644)
    msg, err := cat(shell, []string{tempFile})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    if msg != "let me innnnnnnnnnnn" {
        t.Errorf("Expected 'hello world', got: '%s'", msg)
    }

    msg, err = cat(shell, []string{"/nonexistent"})
    if err == nil {
        t.Error("Expected error, got nil")
    }

    msg, err = cat(shell, []string{})
    if err == nil {
        t.Error("Expected error, got nil")
    }
}

func TestExit(t *testing.T) {
    shell := New()

    msg, err := exit(shell, []string{})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    if msg != "exit status 0" {
        t.Errorf("Expected 'exit status 0', got: '%s'", msg)
    }
    if shell.status != false {
        t.Error("Expected shell status to be false, got true")
    }

    shell.status = true
    msg, err = exit(shell, []string{"1"})
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    if msg != "exit status 1" {
        t.Errorf("Expected 'exit status 1', got: '%s'", msg)
    }
    if shell.status != false {
        t.Error("Expected shell status to be false, got true")
    }

    
    msg, err = exit(shell, []string{"random things"})
    if err == nil {
        t.Error("Expected error, got nil")
    }

    msg, err = exit(shell, []string{"1", "2"})
    if err == nil {
        t.Error("Expected error, got nil")
    }
}

