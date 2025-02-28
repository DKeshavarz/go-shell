package shell

import (
	"testing"

	"systemgroup.net/bootcamp/go/v1/shell/internal/models"
	"systemgroup.net/bootcamp/go/v1/shell/internal/servise"
)




func TestALL(t *testing.T) {
	s := New()
	var testUser = []models.User{
		{Username: "test_ali" ,Password: "123"},
		{Username: "test_ali2",Password: ""},
	}
	
	servise.CreateUser(&testUser[0])
	args := []string{"test_ali", "123"}

	msg, err := login(s, args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.CurrentUser.Username != testUser[0].Username {
		t.Fatalf("expected CurrentUser Username to be %s, got %s", testUser[0].Username, s.CurrentUser.Username)
	}
	if msg != "" {
		t.Fatalf("expected empty message, got %v", msg)
	}

	argsWrongPass := []string{"testuser", "wrongpass"}
	msg, err = login(s, argsWrongPass)
	if err == nil {
		t.Fatalf("expected error when logging in with wrong password, got none")
	}

	argsNonExisting := []string{"nonexistentuser"}
	msg, err = login(s, argsNonExisting)
	if err == nil {
		t.Fatalf("expected error when logging in with nonexistent user, got none")
	}

	argsValid := []string{testUser[1].Username, testUser[1].Password}
	msg, err = addUser(s, argsValid)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.CurrentUser.Username != "test_ali2" {
		t.Fatalf("expected CurrentUser Username to be 'test_ali2', got %s", s.CurrentUser.Username)
	}
	if msg != "" {
		t.Fatalf("expected empty message, got %v", msg)
	}

	msg, err = logout(s, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.CurrentUser != nil {
		t.Fatalf("expected CurrentUser to be nil after logout, got %v", s.CurrentUser)
	}
	if msg != "" {
		t.Fatalf("expected empty message, got %v", msg)
	}

}

