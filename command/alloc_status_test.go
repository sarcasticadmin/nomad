package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAllocStatusCommand_Implements(t *testing.T) {
	var _ cli.Command = &AllocStatusCommand{}
}

func TestAllocStatusCommand_Fails(t *testing.T) {
	srv, _, url := testServer(t, nil)
	defer srv.Stop()

	ui := new(cli.MockUi)
	cmd := &AllocStatusCommand{Meta: Meta{Ui: ui}}

	// Fails on misuse
	if code := cmd.Run([]string{"some", "bad", "args"}); code != 1 {
		t.Fatalf("expected exit code 1, got: %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, cmd.Help()) {
		t.Fatalf("expected help output, got: %s", out)
	}
	ui.ErrorWriter.Reset()

	// Fails on connection failure
	if code := cmd.Run([]string{"-address=nope", "foobar"}); code != 1 {
		t.Fatalf("expected exit code 1, got: %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, "Error querying allocation") {
		t.Fatalf("expected failed query error, got: %s", out)
	}
	ui.ErrorWriter.Reset()

	// Fails on missing alloc
	if code := cmd.Run([]string{"-address=" + url, "26470238-5CF2-438F-8772-DC67CFB0705C"}); code != 1 {
		t.Fatalf("expected exit 1, got: %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, "No allocation(s) with prefix or id") {
		t.Fatalf("expected not found error, got: %s", out)
	}

	// Fail on uneven identifier length
	if code := cmd.Run([]string{"-address=" + url, "26470238-5CF2-438F-8772-DC67CFB0705"}); code != 1 {
		t.Fatalf("expected exit 1, got: %d", code)
	}
	if out := ui.ErrorWriter.String(); !strings.Contains(out, "must be of even length.") {
		t.Fatalf("expected even length error, got: %s", out)
	}
}
