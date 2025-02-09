package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConfirmAction() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you want to perform this action on all instances? (yes/no): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	return strings.ToLower(response) == "yes"
}
