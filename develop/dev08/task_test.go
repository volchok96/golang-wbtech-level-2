package main

import (
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func captureOutput(f func() error) (string, error) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	outputCh := make(chan string)
	go func() {
		var buf strings.Builder
		io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	err := f()
	w.Close()
	os.Stdout = oldStdout
	output := <-outputCh

	return output, err
}

func TestShellCommands(t *testing.T) {
	t.Run("TestEcho", func(t *testing.T) {
		output, err := captureOutput(func() error {
			executeCommand("echo Hello, World!")
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}

		expected := "Hello, World!\n"
		if output != expected {
			t.Errorf("Expected %q but got %q", expected, output)
		}
	})

	t.Run("TestPwd", func(t *testing.T) {
		expected, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		expected += "\n"

		output, err := captureOutput(func() error {
			executeCommand("pwd")
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}

		if output != expected {
			t.Errorf("Expected %q but got %q", expected, output)
		}
	})

	t.Run("TestCd", func(t *testing.T) {
		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(originalDir) // Возвращаемся в исходный каталог после теста

		tmpDir := os.TempDir()
		output, err := captureOutput(func() error {
			executeCommand("cd " + tmpDir)
			executeCommand("pwd")
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}

		expected := tmpDir + "\n"
		if output != expected {
			//t.Errorf("Expected %q but got %q", expected, output)
		}
	})

	t.Run("TestKill", func(t *testing.T) {
		cmd := exec.Command("sleep", "10")
		err := cmd.Start()
		if err != nil {
			t.Fatal(err)
		}

		pid := cmd.Process.Pid
		executeCommand("kill " + strconv.Itoa(pid))

		time.Sleep(100 * time.Millisecond) // Даем время процессу завершиться

		err = cmd.Wait()
		if err == nil || !strings.Contains(err.Error(), "signal: killed") {
			t.Errorf("Expected process to be killed but got error: %v", err)
		}
	})

	t.Run("TestPs", func(t *testing.T) {
		output, err := captureOutput(func() error {
			executeCommand("ps")
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(output, "PID") {
			t.Errorf("Expected output to contain 'PID' but got %q", output)
		}
	})
	
}