package shell

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		input string
		out   []string
		out2  []string
	}{
		{"exit 0    5  ", []string{"exit", "0", "5"}, nil},
		{"cd ..", []string{"cd", ".."}, nil},
		{"echo '$path   p ' ls     l s", []string{"echo", "$path   p ", "ls", "l", "s"}, nil},
		{"no    yes    no", []string{"no", "yes", "no"}, nil},
		{"'no$PATH'", []string{"no$PATH"}, nil},
		{"no$unsetvar", []string{"no"}, nil},
		{"\"no $unsetvar TestTokenizertest\"", []string{"no  TestTokenizertest"}, nil},
		{"echo test >> t.txt", []string{"echo", "test"}, []string{">>", "t.txt"}},
	}

	shell := New()
	for _, tc := range testCases {
		input := tc.input
		out := tc.out
		out2 := tc.out2

		realOut, realOut2, _ := shell.tokenizer(input)
		if !sliceAreEqul(realOut, out) {
			t.Fatal("expected : ", out, "\nbut got : ", realOut)
		}
		if !sliceAreEqul(realOut2, out2) {
			t.Fatal("expected : ", out2, "\nbut got : ", realOut2)
		}

	}
}

func TestIsScalbe(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{"\\", true},
		{"p", false},
		{"`", true},
		{"\"", true},
		{"a", false},
	}

	for _, tc := range testCases {
		in := tc.in
		out := tc.out

		realOut := isScalbe(in)
		if out != realOut {
			t.Fatal("expected : ", out, "\nbut got : ", realOut)
		}

	}
}

func TestNew(t *testing.T) {
	shell := New()

	if shell == nil {
		t.Fatal("Expected a new Shell instance, got nil")
	}

	if len(shell.Handlers) == 0 {
		t.Error("Expected handlers to be registered, got none")
	}

	if shell.status != true {
		t.Error("Expected shell status to be true, got false")
	}
}

func TestRead(t *testing.T) {
	shell := New()

	input := "test input\n"
	r, w, _ := os.Pipe()
	w.WriteString(input)
	w.Close()
	os.Stdin = r

	result := shell.read()
	if result != "test input" {
		t.Errorf("Expected 'test input', got '%s'", result)
	}
}

func TestTokenizer2(t *testing.T) {
	shell := New()

	testCases := []struct {
		input    string
		expected []string
		redirect []string
		err      bool
	}{
		{"echo hello", []string{"echo", "hello"}, nil, false},
		{"echo 'hello world'", []string{"echo", "hello world"}, nil, false},
		{"echo \"hello world\"", []string{"echo", "hello world"}, nil, false},
		{"echo hello > output.txt", []string{"echo", "hello"}, []string{">", "output.txt"}, false},
		{"echo hello >> output.txt", []string{"echo", "hello"}, []string{">>", "output.txt"}, false},
		{"echo hello 2> error.txt", []string{"echo", "hello"}, []string{"2>", "error.txt"}, false},
		{"echo hello 2>> error.txt", []string{"echo", "hello"}, []string{"2>>", "error.txt"}, false},
		{"echo 'unclosed quote", nil, nil, true},
		{"echo \"unclosed quote", nil, nil, true},
	}

	for _, tc := range testCases {
		tokens, redirect, err := shell.tokenizer(tc.input)

		if tc.err {
			if err == nil {
				t.Errorf("Expected error for input '%s', got none", tc.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
			continue
		}

		if !sliceAreEqul(tokens, tc.expected) {
			t.Errorf("For input '%s', expected tokens %v, got %v", tc.input, tc.expected, tokens)
		}

		if !sliceAreEqul(redirect, tc.redirect) {
			t.Errorf("For input '%s', expected redirect %v, got %v", tc.input, tc.redirect, redirect)
		}
	}
}

func TestSystemCommand(t *testing.T) {
	shell := New()

	output, err := shell.systemCommand([]string{"echo", "hello"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if strings.TrimSpace(output) != "hello" {
		t.Errorf("Expected 'hello', got '%s'", output)
	}

	_, err = shell.systemCommand([]string{"invalid_command"})
	if err == nil {
		t.Error("Expected error for invalid command, got none")
	}
}


func TestGetEnvVar(t *testing.T) {
	
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnvVar("TEST_VAR")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	result = getEnvVar("NON_EXISTENT_VAR")
	if result != "" {
		t.Errorf("Expected empty string for non-existent variable, got '%s'", result)
	}
}

func TestShow(t *testing.T) {
	shell := New()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	shell.show(nil, "test message", nil)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := strings.TrimSpace(buf.String())

	if output != "test message" {
		t.Errorf("Expected 'test message', got '%s'", output)
	}
}

//-------------------helper-------------------------------

func sliceAreEqul(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
