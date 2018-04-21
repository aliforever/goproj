package main

import (
	"flag"
	"fmt"
	"goproj/file"
	"os"
	"strings"
	"os/exec"
)

func main() {
	projectType := flag.String("type", "new_project", "Enter Project Name'")
	botName := flag.String("username", "bot_username", "Enter Bot Username")
	botToken := flag.String("token", "bot_token", "Enter Bot Token")
	flag.Parse()
	if *projectType == "bot" {
		CreateBotProject(*botName, *botToken)
	} else {
		fmt.Println("Project is not supported yet!")
	}
}

func CreateBotProject(username, token string) {
	goPath := os.Getenv("GOPATH") + "/src/"
	currentPath := goPath + "goproj/"
	goSourcePath := goPath + username
	apiBytes, _ := file.FileGetContents(currentPath + "templates/bot/api.temp")
	databaseBytes, _ := file.FileGetContents(currentPath + "templates/bot/database.temp")
	keyboardsBytes, _ := file.FileGetContents(currentPath + "templates/bot/keyboards.temp")
	languageBytes, _ := file.FileGetContents(currentPath + "templates/bot/language.temp")
	engineBytes, _ := file.FileGetContents(currentPath + "templates/bot/engine.temp") // %BOTUSERNAME% %BOTUSERNAME_CAPS%
	engineBytes = []byte(strings.Replace(string(engineBytes), "%BOTUSERNAME%", username, -1))
	engineBytes = []byte(strings.Replace(string(engineBytes), "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1))
	configBytes, _ := file.FileGetContents(currentPath + "templates/bot/config.temp") // %BOTUSERNAME% %BOTTOKEN%
	configBytes = []byte(strings.Replace(string(configBytes), "%BOTUSERNAME%", username, -1))
	configBytes = []byte(strings.Replace(string(configBytes), "%BOTTOKEN%", token, -1))
	mainBytes, _ := file.FileGetContents(currentPath + "templates/bot/main.temp") // %BOTUSERNAME%
	mainBytes = []byte(strings.Replace(string(mainBytes), "%BOTUSERNAME%", username, -1))
	mainBytes = []byte(strings.Replace(string(mainBytes), "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1))
	userBytes, _ := file.FileGetContents(currentPath + "templates/bot/user.temp")
	persianBytes, _ := file.FileGetContents(currentPath + "templates/bot/persian.temp")

	configPath := goSourcePath + "/configs/"
	funcsPath := goSourcePath + "/funcs/"
	langPath := goSourcePath + "/lang/"
	modelPath := goSourcePath + "/models/"
	os.Mkdir(goSourcePath, os.ModePerm)
	os.Mkdir(goSourcePath + "/configs/", os.ModePerm)
	os.Mkdir(goSourcePath + "/funcs/", os.ModePerm)
	os.Mkdir(goSourcePath + "/lang/", os.ModePerm)
	os.Mkdir(goSourcePath + "/models/", os.ModePerm)
	if err := file.FilePutContents(configPath + "main.go", configBytes); err != nil {
		fmt.Println(err)
		return
	}
	file.FilePutContents(funcsPath + "engine.go", engineBytes)
	file.FilePutContents(funcsPath + "keyboards.go", keyboardsBytes)
	file.FilePutContents(modelPath + "api.go", apiBytes)
	file.FilePutContents(modelPath + "database.go", databaseBytes)
	file.FilePutContents(modelPath + "user.go", userBytes)
	file.FilePutContents(langPath + "language.go", languageBytes)
	file.FilePutContents(langPath + "persian.go", persianBytes)
	file.FilePutContents(goSourcePath + "/main.go", mainBytes)
	exec.Command("go fmt ", "-w", goSourcePath)
	exec.Command("goimports", "-w", goSourcePath)
	fmt.Println(fmt.Sprintf("Bot %s Created at %s", username, goSourcePath))
}