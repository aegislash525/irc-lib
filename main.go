package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	// "time"
)

func main() {
	var (
		server  string
		port    string
		channel string
		nick    string
	)
	flag.StringVar(&server, "s", "irc.freenode.net", "server to connect")
	flag.StringVar(&port, "p", "6667", "port of the server")
	flag.StringVar(&nick, "n", "go-irc-usr", "your nickname")
	flag.Parse()

	conn, err := net.Dial("tcp", server+":"+port)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn) // here we recieve messages
	writer := bufio.NewWriter(conn) // here we insert our commands

	fmt.Print("\033[H\033[2J")

	fmt.Fprintf(conn, "NICK %s\r\n", nick)
	fmt.Fprintf(conn, "USER %s 0 * :%s\r\n", nick, nick)
	writer.Flush()

	// handle ping pong and print received messages
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}

			// copy last line PING :U`osMYiWW]
			// replace PING to PONG and send
			if strings.HasPrefix(line, "PING") {
				pong := strings.Replace(line, "PING", "PONG", 1)
				fmt.Fprintf(writer, "%s\r\n", pong)
				writer.Flush()
				continue
			}
			line = strings.Trim(line, "\r\n")
			// line = strings.Trim(line, ":")
			// line = strings.Split(line, ":")[1]
			// print messages received from server
			// fmt.Printf("[%s] %s\n",
			// 	time.Now().Format(time.TimeOnly), line)
			line = strings.Replace(line, ":*.freenode.net ", "", 1)
			fmt.Printf("%s\n", line)
		}
	}()

	writer.Flush()

	// user input
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		text, err := stdinReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			return
		}
		// trim user input
		text = strings.Trim(text, "\r\n")
		if text == "" {
			continue
		} else if text[0] == '/' {
			// COMMANDS
			cmd := strings.Split(text, "/")
			cmd = strings.Split(cmd[1], " ")
			switch cmd[0] {
			case "exit":
				return
			case "list":
				fmt.Fprintf(writer, "%s\r\n", "LIST")
				writer.Flush()
				continue
			case "help":
				fmt.Fprintf(writer, "%s\r\n", "HELP")
				writer.Flush()
				continue
			case "join":
				strs := strings.Split(text, " ")
				if len(strs) > 1 {
					fmt.Fprintf(writer, "JOIN %s\r\n", strs[1])
					channel = strs[1]
					writer.Flush()
				}
			default:
				continue
			}
		} else {
			fmt.Fprintf(writer, "PRIVMSG %s %s\r\n", channel, text)
			writer.Flush()
		}
	}
}
