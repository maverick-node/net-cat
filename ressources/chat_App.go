package netcat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const MaxUsers = 10

type User struct {
	Name string
	Conn net.Conn
}

var (
	chatLogo   string
	users      []User
	backUp     []string
	logBackUp  []string
	mu         sync.Mutex
)

// Checks if the name provided is a valid name
func isValidName(name string) bool {
	if len(strings.TrimSpace(name)) == 0 {
		return false
	}

	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}
	return true
}

func invalidName(conn net.Conn) {
	fmt.Fprintln(conn, "Invalid name. Name must:")
	fmt.Fprintln(conn, "- Not be empty")
	fmt.Fprintln(conn, "- Only contain letters, numbers, underscore (_), or hyphen (-)")
	fmt.Fprintln(conn, "Please try again.")
}

// Reads the name, checks the name if valid using isValidName function, checks the name if already exists
func readValidName(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	for {
		fmt.Fprint(conn, "[ENTER YOUR NAME]:")
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		name = strings.TrimSpace(name)
		if !isValidName(name) {
			invalidName(conn)
			continue
		}
		mu.Lock()
		nameExists := false
		for _, user := range users {
			if strings.EqualFold(user.Name, name) {
				nameExists = true
				break
			}
		}
		mu.Unlock()
		if nameExists {
			fmt.Fprintln(conn, "This name is already taken. Please choose another name.")
			log.Printf("Client fails to change their name %s: %s", conn.RemoteAddr().String(), "("+name+")")
			continue
		}
		return name, nil
	}
}

// Remove the user from the slice if logged out
func removeUser(name string) {
	mu.Lock()
	defer mu.Unlock()
	for i, user := range users {
		if user.Name == name {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
}

// Updates the name when the user change their name
func updateUserInList(oldName, newName string) {
	mu.Lock()
	defer mu.Unlock()
	for i := range users {
		if users[i].Name == oldName {
			users[i].Name = newName
			break
		}
	}
}

// broadcast sends notifications of various events (join, message, leave, change) to all users in the chat.
func broadcast(eventType string, content string, senderName string) {
	mu.Lock()
	defer mu.Unlock()
	time := time.Now()
	for _, user := range users {
		if eventType == "name" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n%s has joined our chat...\n", senderName)
		} else if eventType == "message" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:%s", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), senderName, content)
		} else if eventType == "leave" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n%s has left the chat\n", senderName)
			logBackUp = append(logBackUp, fmt.Sprintf("%s has left the chat\n", senderName))
		} else if eventType == "change" && user.Name != senderName {
			fmt.Fprintf(user.Conn, "\n%s has changed their name to %s\n", content, senderName)
			logBackUp = append(logBackUp, fmt.Sprintf("%s has changed their name to %s\n", content, senderName))
		}
		fmt.Fprintf(user.Conn, "[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), user.Name)
	}

	if eventType == "message" && content != "" {
		backUp = append(backUp, fmt.Sprintf("[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:%s", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), senderName, content))
	} else if (eventType == "change" || eventType == "leave") && content != "" {
		backUp = append(backUp, logBackUp...)
	}
}

// HandleClient manages the interaction with a connected client, including name validation and message handling.
func HandleClient(conn net.Conn) {
	mu.Lock()
	if len(users) >= MaxUsers {
		mu.Unlock()
		fmt.Fprint(conn, "Sorry, the chat room is full (maximum 10 users). Please try again later.\n")
		log.Printf("Server reaches the maximum users")
		conn.Close()
		return
	}
	mu.Unlock()

	defer conn.Close()

	var err error

	chatLogo, err = LoadChatLogo("./ressources/welcome.txt")
	if err != nil {
		log.Fatalf("Error loading chat logo: %v", err)
		log.Printf("Error loading chat logo: %v", err)
		return
	}
	fmt.Fprint(conn, string([]byte(chatLogo)))
	name, err := readValidName(conn)
	if err != nil {
		log.Printf("Error reading name: %v", err)
		return
	}

	mu.Lock()
	for _, v := range backUp {
		fmt.Fprint(conn, v)
	}
	mu.Unlock()

	mu.Lock()
	users = append(users, User{Name: name, Conn: conn})
	mu.Unlock()

	broadcast("name", "", name)
	log.Printf("New client connected: %s %s", conn.RemoteAddr(), "("+name+")")

	for {
		time := time.Now()
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			removeUser(name)
			broadcast("leave", "", name)
			log.Printf("Client disconnected: %s %s", conn.RemoteAddr(), "("+name+")")
			return
		}
		if message != "\n" {
			if message == "--name\n" {
				oldName := name
				for {
					fmt.Fprint(conn, "Please enter new username:")
					log.Printf("Client tries to change their name: %s %s", conn.RemoteAddr(), "("+name+")")
					newName, err := bufio.NewReader(conn).ReadString('\n')
					if err != nil {
						log.Printf("Error reading name: %v", err)
						return
					}
					newName = strings.TrimSpace(newName)

					if !isValidName(newName) {
						fmt.Fprintln(conn, "Invalid name. Please try again.")
						log.Printf("Client tries to enter an invalid name %s: %s", conn.RemoteAddr().String(), "("+name+")")
						continue
					}

					mu.Lock()
					nameExists := false
					for _, user := range users {
						if strings.EqualFold(user.Name, newName) {
							nameExists = true
							break
						}
					}
					mu.Unlock()

					if nameExists {
						fmt.Fprintln(conn, "This name is already taken. Please choose another name.")
						log.Printf("Client fails to change their name %s: %s", conn.RemoteAddr().String(), "("+name+")")
						continue
					}
					name = newName
					updateUserInList(oldName, newName)
					broadcast("change", oldName, newName)
					log.Printf("Client %s %s changed their name succesfully to %s", conn.RemoteAddr(), "("+oldName+")", "("+newName+")")
					break
				}
			} else {
				broadcast("message", message, name)
				log.Printf("Message received from %s %s: %s", conn.RemoteAddr().String(), "("+name+")", message)
			}
		} else if len(message) == 1 {
			fmt.Fprint(conn, "You cannot submit an empty message.")
			log.Printf("Client tries to send empty message %s: %s", conn.RemoteAddr().String(), "("+name+")")
			fmt.Fprintf(conn, "\n[%d-%.2d-%.2d %.2d:%.2d:%.2d][%s]:", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second(), name)
		}
	}
}
