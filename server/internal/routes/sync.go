package routes

import (
	"fmt"
	"net/http"
	"time"
)

func (c *Controller) uploadFile(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(10 << 20) // max of 10MB
	if err != nil {
		fmt.Printf("Failed to parse file: %v", err)
		return
	}

	file, header, err := request.FormFile("file")
	if err != nil {
		fmt.Printf("Could not get file: %v", err)
		return
	}
	defer file.Close()

	modifiedDate := request.FormValue("modified-date")
	if modifiedDate == "" {
		modifiedDate = time.Now().Format(time.RFC3339)
	}

	fmt.Printf("File: %s", header.Filename)
	fmt.Printf("Size: %d", header.Size)
	fmt.Printf("Header: %v", header.Header)
	fmt.Printf("modified: %s", modifiedDate)
}
