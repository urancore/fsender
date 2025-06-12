package pkg

import (
	"fmt"
	"time"
)


// date | USER_IP | METHOD | message
func Log(userIP, method, message string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	datePart := fmt.Sprintf("\x1b[3;90m%s\x1b[0m", now)

	userIPPart := fmt.Sprintf("\x1b[37m%s\x1b[0m", userIP)

	var methodStyle string
	switch method {
	case "POST":
		methodStyle = "\x1b[46;37m"
	case "GET":
		methodStyle = "\x1b[42;37m"
	case "PUT":
		methodStyle = "\x1b[43;37m"
	case "PATCH":
		methodStyle = "\x1b[48;5;208;37m"
	default:
		methodStyle = "\x1b[37m"
	}

	methodPart := fmt.Sprintf("%s%s\x1b[0m", methodStyle, method)

	logMsg := fmt.Sprintf("%s | %s | %s | %s", datePart, userIPPart, methodPart, message)
	fmt.Println(logMsg)
}
