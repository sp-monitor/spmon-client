package main

import (
	"fmt"
	"io"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"os/exec"
)

func uploadModule(w http.ResponseWriter, r *http.Request) {
	logger, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_DEBUG, "SPMON-CLIENT")
	if err != nil {
		log.Fatal(err)
	}

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Err(fmt.Sprintf("Error Retrieving File: %+v", err))
		return
	}

	defer file.Close()
	
	logger.Notice(fmt.Sprintf("Module Received: %+v", handler.Filename))

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	// moduleFile, err := os.CreateTemp("modules", handler.Filename)
	moduleFile, err := os.Create(fmt.Sprintf("modules/%s", handler.Filename))
	if err != nil {
		logger.Err(fmt.Sprintf("Error Creating Temp File: %+v", err))
	}
	defer moduleFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Err(fmt.Sprintf("Error Reading File: %+v", err))
	}
	moduleFile.Write(fileBytes)
	
	cmd := exec.Command("./nmp/nmp")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		logger.Err(fmt.Sprintf("Error Running NMP: %+v", err))
	}
}

func setupRoutes() {
	http.HandleFunc("/module/upload", uploadModule)
	http.ListenAndServe(":3999", nil)
}

func main() {
	fmt.Println("Server listening on port 3999")
	setupRoutes()
}