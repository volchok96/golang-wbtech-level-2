package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// Смена директории (cd)
func changeDirectory(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("cd: missing argument")
	}
	return os.Chdir(args[1])
}

// Отображение текущей директории (pwd)
func printWorkingDirectory() (string, error) {
	return os.Getwd()
}

// Вывод аргументов в STDOUT (echo)
func echoOutput(args []string) string {
	return strings.Join(args[1:], " ")
}

// Завершение процесса по PID (kill)
func terminateProcess(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("kill: missing argument")
	}
	pid, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("kill: invalid PID")
	}
	return syscall.Kill(pid, syscall.SIGKILL)
}

// Отображение информации о процессах (ps)
func listProcesses() error {
	cmd := exec.Command("ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Graceful shutdown
func handleInterrupt() {
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)
	go func() {
		<-interruptSignal
		fmt.Println("\nReceived interrupt signal. Exiting...")
		os.Exit(0)
	}()
}

// Выполнение команды с поддержкой пайпов
func executeCommand(input string) {
	// Разделяем команды по пайпам
	commands := strings.Split(input, "|")
	totalCommands := len(commands)

	var previousCommand *exec.Cmd
	var previousOutput io.ReadCloser

	for i, command := range commands {
		command = strings.TrimSpace(command)
		args := strings.Fields(command)

		if len(args) == 0 {
			continue
		}

		if args[0] == "exit" || args[0] == "quit" {
			os.Exit(0) // Правильная обработка команды exit
		}

		var cmd *exec.Cmd

		switch args[0] {
		case "cd":
			if err := changeDirectory(args); err != nil {
				fmt.Println(err)
			}
			return
		case "pwd":
			dir, err := printWorkingDirectory()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(dir)
			}
			return
		case "echo":
			fmt.Println(echoOutput(args))
			return
		case "kill":
			if err := terminateProcess(args); err != nil {
				fmt.Println(err)
			}
			return
		case "ps":
			if err := listProcesses(); err != nil {
				fmt.Println(err)
			}
			return
		default:
			cmd = exec.Command(args[0], args[1:]...)
		}

		// Если это не первая команда в пайпе
		if previousCommand != nil {
			stdinPipe, err := cmd.StdinPipe() // канал (io.WriteCloser) подключен к stdin команды cmd
			if err != nil {
				fmt.Println("Error creating stdin pipe:", err)
				return
			}
			// Копируем данные из предыдущего выхода в новый stdin
			go func() {
				defer stdinPipe.Close()
				io.Copy(stdinPipe, previousOutput)
			}()
		}

		// Если это не последняя команда, то готовим пайп
		if i < totalCommands-1 {
			var err error
			previousOutput, err = cmd.StdoutPipe()
			if err != nil {
				fmt.Println("Error creating stdout pipe:", err)
				return
			}
			cmd.Stderr = os.Stderr
		} else {
			// Последняя команда выводит результат в стандартный поток
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		// Запускаем команду
		if err := cmd.Start(); err != nil { // Go вызывает системные вызовы fork() и exec()
			fmt.Println("Error starting command:", err)
			return
		}

		// Если предыдущая команда есть, ждем ее завершения
		if previousCommand != nil {
			if err := previousCommand.Wait(); err != nil {
				fmt.Println("Error waiting for previous command:", err)
				return
			}
		}

		// Устанавливаем текущую команду как предыдущую для следующей итерации
		previousCommand = cmd
	}

	// Ожидаем завершения последней команды
	if previousCommand != nil {
		if err := previousCommand.Wait(); err != nil {
			fmt.Println("Error waiting for command:", err)
		}
	}
}

// Главная функция, которая запускает шелл
func main() {
	handleInterrupt()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("shell> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		if strings.TrimSpace(input) == "exit" || strings.TrimSpace(input) == "quit" {
			break
		}

		executeCommand(input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
