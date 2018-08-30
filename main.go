package main

import (
	"flag"
	"fmt"
	"goproj/file"
	"os"
	"strings"
	"os/exec"
	"strconv"
	"errors"
)

func GoPATH() string {
	return  os.Getenv("GOPATH") + "/src/"
}

func TemplatePath() string {
	return GoPATH() + "goproj/templates"
}

func ProjectPath(name string)  string {
	return  GoPATH() + name + "/"
}

func CurrentDirectoryProjectName() (directory *string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	// for windows directory
	split := strings.Split(dir, `\src\`)
	if len(split) < 2 {
		// for linux directory
		split = strings.Split(dir, `/src/`)
		if len(split) < 2 {
			err = errors.New("wrong_directory")
			return
		}
	}
	directory = &split[1]
	return
}

func main() {
	projectType := flag.String("type", "new_project", "Enter Project Name'")
	botName := flag.String("username", "bot_username", "Enter Bot Username")
	botToken := flag.String("token", "bot_token", "Enter Bot Token")
	makeItem := flag.String("make", "menu=make_item", "Enter Make Item")
	flag.Parse()
	if *makeItem != "menu=make_item" {
		if *botName == "bot_username" {
			directoryName, err := CurrentDirectoryProjectName()
			if err != nil {
				fmt.Println(err)
				return
			}
			botName = directoryName
		}

		if strings.Contains(*makeItem, "menu") {
			split := strings.Split(*makeItem, ":")
			line := 0
			if len(split) == 3 {
				lineInt, err := strconv.Atoi(split[2])
				if err == nil {
					line = lineInt
				}
			}
			err := CreateMenuForBot(*botName, split[1], line)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(fmt.Sprintf("New Menu %s Added to %s Bot", split[1], *botName))
			}
			return
		}
	}
	if *projectType == "bot" {
		CreateBotProject(*botName, *botToken)
	} else {
		fmt.Println("Project is not supported yet!")
	}
}

func CreateMenuForBot(username, menuName string, line int) error {
	enginePath := ProjectPath(username) + "funcs/engine.go"
	interfacePath := ProjectPath(username) + "lang/language.go"
	persianPath := ProjectPath(username) + "lang/persian.go"
	englishPath := ProjectPath(username) + "lang/english.go"
	menuByte, err := file.FileGetContents(TemplatePath() + "/bot/menu.temp")
	if err != nil {
		return err
	}
	menu := strings.Replace(string(menuByte), "%MENU%", menuName, -1)
	menu = strings.Replace(menu, "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1)
	split := strings.Split(menu, "@language")
	menu = split[0]
	langParts := split[1]
	split = strings.Split(langParts, "--")
	langInterface := split[0]
	persian := split[1]
	english := split[2]
	currentInterfaceBytes, err := file.FileGetContents(interfacePath)
	currentPersianBytes, err := file.FileGetContents(persianPath)
	currentEnglishBytes, err := file.FileGetContents(englishPath)
	currentEngineBytes, err := file.FileGetContents(enginePath)
	if err != nil {
		return err
	}
	newInterfaceBytes := []byte((string(currentInterfaceBytes) + "\n" + langInterface))
	newPersianBytes := []byte((string(currentPersianBytes) + "\n" + persian))
	newEnglishBytes := []byte((string(currentEnglishBytes) + "\n" + english))
	file.FilePutContents(interfacePath, newInterfaceBytes)
	file.FilePutContents(persianPath, newPersianBytes)
	file.FilePutContents(englishPath, newEnglishBytes)
	var newEngineBytes []byte
	if line == 0 {
		newEngineBytes = []byte((string(currentEngineBytes) + "\n" + menu))
	} else {
		split := strings.Split(string(currentEngineBytes), "\n")

		if len(split) < line {
			newEngineBytes = []byte((string(currentEngineBytes) + "\n" + menu))
		} else {
			for i := range split {
				if i+1 == line {
					engine := strings.Join(split[:i+1], "\n") + "\n" + menu + "\n" + strings.Join(split[i+1:], "\n")
					newEngineBytes = []byte(engine)
					break
				}
			}
		}
	}

	file.FilePutContents(enginePath, newEngineBytes)
	return nil
}

func CreateBotProject(username, token string) {
	goPath := os.Getenv("GOPATH") + "/src/"
	currentPath := goPath + "goproj/"
	goSourcePath := goPath + username
	apiBytes, _ := file.FileGetContents(currentPath + "templates/bot/api.temp")
	databaseBytes, _ := file.FileGetContents(currentPath + "templates/bot/database.temp")
	keyboardsBytes, _ := file.FileGetContents(currentPath + "templates/bot/keyboards.temp")
	keyboardsBytes = []byte(strings.Replace(string(keyboardsBytes), "%BOTUSERNAME%", username, -1))
	languageBytes, _ := file.FileGetContents(currentPath + "templates/bot/language.temp")
	engineBytes, _ := file.FileGetContents(currentPath + "templates/bot/engine.temp") // %BOTUSERNAME% %BOTUSERNAME_CAPS%
	engineBytes = []byte(strings.Replace(string(engineBytes), "%BOTUSERNAME%", username, -1))
	engineBytes = []byte(strings.Replace(string(engineBytes), "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1))
	methodsBytes, _ := file.FileGetContents(currentPath + "templates/bot/methods.temp") // %BOTUSERNAME% %BOTUSERNAME_CAPS%
	methodsBytes = []byte(strings.Replace(string(methodsBytes), "%BOTUSERNAME%", username, -1))
	methodsBytes = []byte(strings.Replace(string(methodsBytes), "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1))
	configBytes, _ := file.FileGetContents(currentPath + "templates/bot/config.temp") // %BOTUSERNAME% %BOTTOKEN%
	configBytes = []byte(strings.Replace(string(configBytes), "%BOTUSERNAME%", username, -1))
	configBytes = []byte(strings.Replace(string(configBytes), "%BOTTOKEN%", token, -1))
	mainBytes, _ := file.FileGetContents(currentPath + "templates/bot/main.temp") // %BOTUSERNAME%
	mainBytes = []byte(strings.Replace(string(mainBytes), "%BOTUSERNAME%", username, -1))
	mainBytes = []byte(strings.Replace(string(mainBytes), "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1))
	userBytes, _ := file.FileGetContents(currentPath + "templates/bot/user.temp")
	persianBytes, _ := file.FileGetContents(currentPath + "templates/bot/persian.temp")
	englishBytes, _ := file.FileGetContents(currentPath + "templates/bot/english.temp")

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
	file.FilePutContents(funcsPath + "methods.go", methodsBytes)
	file.FilePutContents(funcsPath + "keyboards.go", keyboardsBytes)
	file.FilePutContents(modelPath + "api.go", apiBytes)
	file.FilePutContents(modelPath + "database.go", databaseBytes)
	file.FilePutContents(modelPath + "user.go", userBytes)
	file.FilePutContents(langPath + "language.go", languageBytes)
	file.FilePutContents(langPath + "persian.go", persianBytes)
	file.FilePutContents(langPath + "english.go", englishBytes)
	file.FilePutContents(goSourcePath + "/main.go", mainBytes)
	exec.Command("go fmt ", "-w", goSourcePath).Output()
	exec.Command("goimports", "-w", goSourcePath).Output()
	fmt.Println(fmt.Sprintf("Bot %s Created at %s", username, goSourcePath))
}
