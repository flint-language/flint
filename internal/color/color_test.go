package color

import "testing"

func TestColorDisabled(t *testing.T) {
	Enabled = false
	defer func() { Enabled = detectTerminal() }()

	out := Color("hello", Red, Bold)

	if out != "hello" {
		t.Fatalf("expected unchanged text, got %q", out)
	}
}

func TestColorEnabled(t *testing.T) {
	Enabled = true
	defer func() { Enabled = detectTerminal() }()

	out := Color("hi", Red)

	expected := string(Red) + "hi" + string(Reset)

	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestMultipleStyles(t *testing.T) {
	Enabled = true
	defer func() { Enabled = detectTerminal() }()

	out := Color("x", Bold, Green)

	expected := string(Bold) + string(Green) + "x" + string(Reset)

	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestColorNoStyles(t *testing.T) {
	Enabled = true
	defer func() { Enabled = detectTerminal() }()

	out := Color("text")

	expected := "text" + string(Reset)

	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestResetNotDuplicated(t *testing.T) {
	Enabled = true
	defer func() { Enabled = detectTerminal() }()

	out := Color("abc", Reset)

	expected := string(Reset) + "abc" + string(Reset)

	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestRedTextWrapper(t *testing.T) {
	Enabled = true
	defer func() { Enabled = detectTerminal() }()

	out := RedText("error")

	expected := string(Red) + "error" + string(Reset)

	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func TestBoldGreenWrapper(t *testing.T) {
	Enabled = true
	defer func() { Enabled = detectTerminal() }()

	out := BoldGreen("ok")

	expected := string(Bold) + string(Green) + "ok" + string(Reset)

	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}
