package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aliforever/goproj/file"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

func GoPATH() string {
	return os.Getenv("GOPATH") + "/src/"
}

func TemplatePath() string {
	return GoPATH() + "github.com/aliforever/goproj/templates"
}

func ProjectPath(name string) string {
	return GoPATH() + name + "/"
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
	makeItem := flag.String("make", "make=make_item", "Enter Make Item")
	addItem := flag.String("add", "add=item", "Enter Add Item")
	golandWatchers := flag.String("golandWatchers", "1", "Should it enable go fmt and goimports ")
	flag.Parse()
	if *botName == "bot_username" {
		directoryName, err := CurrentDirectoryProjectName()
		if err != nil {
			fmt.Println(err)
			return
		}
		botName = directoryName
	}
	if *makeItem != "make=make_item" {
		if strings.Index(*makeItem, "menu") == 0 {
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
		} else if strings.Index(*makeItem, "model") == 0 {
			split := strings.Split(*makeItem, ":")
			err := CreateModelForBot(*botName, split[2], split[3], split[4])
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(fmt.Sprintf("New Model %s Added to %s Bot", split[1], *botName))
			}
			return
		} else if strings.Index(*makeItem, "inline_menu") == 0 {
			split := strings.Split(*makeItem, ":")
			line := 0
			menuType := "simple"
			if len(split) == 3 {
				lineInt, err := strconv.Atoi(split[2])
				if err == nil {
					line = lineInt
				}
			} else if len(split) == 4 {
				lineInt, err := strconv.Atoi(split[2])
				if err == nil {
					line = lineInt
				}
				menuType = split[3]
			}

			err := CreateInlineMenuForBot(*botName, split[1], line, menuType)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(fmt.Sprintf("New Menu %s Added to %s Bot", split[1], *botName))
			}
			return
		}
	}
	if *addItem != "add=item" {
		split := strings.Split(*addItem, ":")
		if len(split) != 2 {
			fmt.Println("wrong_command")
			return
		}
		if split[0] == "text" {
			AddTextToLanguage(*botName, split[1])
			return
		}
	}
	if *projectType == "bot" {
		watchers := true
		if *golandWatchers != "1" {
			watchers = false
		}
		CreateBotProject(*botName, *botToken, watchers)
	} else {
		fmt.Println("Project is not supported yet!")
	}
}

func AddTextToLanguage(username, textTitle string) error {
	interfacePath := ProjectPath(username) + "lang/language.go"
	persianPath := ProjectPath(username) + "lang/persian.go"
	englishPath := ProjectPath(username) + "lang/english.go"
	textBytes, _ := file.FileGetContents(TemplatePath() + "/bot/add/text.tmp")
	textContent := string(textBytes)
	textContent = strings.Replace(textContent, "%TITLE%", textTitle, -1)
	textContentInterface := textContent[strings.Index(textContent, "%INTERFACE%")+len("%INTERFACE%") : strings.Index(textContent, "%/INTERFACE%")]
	textContentInterface = strings.Replace(textContentInterface, "\n", "", -1)
	textContentEnglish := textContent[strings.Index(textContent, "%ENGLISH%")+len("%ENGLISH%") : strings.Index(textContent, "%/ENGLISH%")]
	textContentPersian := textContent[strings.Index(textContent, "%PERSIAN%")+len("%PERSIAN%") : strings.Index(textContent, "%/PERSIAN%")]
	interfaceBytes, err := file.FileGetContents(interfacePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	englishBytes, err := file.FileGetContents(englishPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	persianBytes, err := file.FileGetContents(persianPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	interfaceContent := string(interfaceBytes)
	persianContent := string(persianBytes)
	englishContent := string(englishBytes)
	interfaceEndBracketIndex := strings.LastIndex(interfaceContent, "}")
	interfaceContent = interfaceContent[:interfaceEndBracketIndex] + textContentInterface + "}"
	persianContent += "\n" + textContentPersian
	englishContent += "\n" + textContentEnglish
	file.FilePutContents(interfacePath, []byte(interfaceContent))
	file.FilePutContents(persianPath, []byte(persianContent))
	file.FilePutContents(englishPath, []byte(englishContent))
	exec.Command("go fmt ", "-w", ProjectPath(username)).Output()
	exec.Command("goimports", "-w", ProjectPath(username)).Output()
	return nil
}

func CreateModelForBot(username, modelFileName, modelStructName, modelTableName string) error {
	modelsPath := ProjectPath(username) + "models/"
	modelShortName := ""

	for _, char := range modelStructName {
		if !unicode.IsLower(char) {
			modelShortName += strings.ToLower(string(char))
		}
	}
	modelBytes, err := file.FileGetContents(TemplatePath() + "/bot/model.temp")
	if err != nil {
		return err
	}
	model := strings.Replace(string(modelBytes), "%MODEL_STRUCT_NAME%", modelStructName, -1)
	model = strings.Replace(model, "%MODEL_SHORT_NAME%", modelShortName, -1)
	model = strings.Replace(model, "%MODEL_TABLE_NAME%", modelTableName, -1)
	err = file.FilePutContents(modelsPath+modelFileName+".go", []byte(model))
	if err != nil {
		return err
	}
	exec.Command("go fmt ", "-w", ProjectPath(username)).Output()
	exec.Command("goimports", "-w", ProjectPath(username)).Output()
	return nil
}

func CreateInlineMenuForBot(username, menuName string, line int, menuType string) error {
	enginePath := ProjectPath(username) + "funcs/engine.go"
	interfacePath := ProjectPath(username) + "lang/language.go"
	persianPath := ProjectPath(username) + "lang/persian.go"
	englishPath := ProjectPath(username) + "lang/english.go"
	menuByte, err := file.FileGetContents(TemplatePath() + "/bot/inline_menu.temp")
	if err != nil {
		return err
	}
	fileContent := string(menuByte)
	if menuType == "simple" {
		idx := strings.Index(fileContent, "%SIMPLES_INLINE_MENU%") + len("%SIMPLES_INLINE_MENU%")
		endIdx := strings.Index(fileContent, "%/SIMPLES_INLINE_MENU%")
		content := fileContent[idx:endIdx]
		menu := strings.Replace(content, "%INLINE_MENU%", menuName, -1)
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
		lastBracketIndex := strings.LastIndex(string(currentInterfaceBytes), "}")
		newInterfaceBytes := []byte((string(currentInterfaceBytes)[:lastBracketIndex-1] + langInterface + string(currentInterfaceBytes[lastBracketIndex:])))
		newPersianBytes := []byte((string(currentPersianBytes) + persian))
		newEnglishBytes := []byte((string(currentEnglishBytes) + english))
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
		exec.Command("go fmt ", "-w", ProjectPath(username)).Output()
		exec.Command("goimports", "-w", ProjectPath(username)).Output()
		return nil
	}
	return errors.New("type_unknown")
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
	lastBracketIndex := strings.LastIndex(string(currentInterfaceBytes), "}")
	newInterfaceBytes := []byte((string(currentInterfaceBytes)[:lastBracketIndex-1] + langInterface + string(currentInterfaceBytes[lastBracketIndex:])))
	newPersianBytes := []byte((string(currentPersianBytes) + persian))
	newEnglishBytes := []byte((string(currentEnglishBytes) + english))
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
	exec.Command("go fmt ", "-w", ProjectPath(username)).Output()
	exec.Command("goimports", "-w", ProjectPath(username)).Output()
	return nil
}

func CreateBotProject(username, token string, addIdeaWatchers bool) {
	goPath := os.Getenv("GOPATH") + "/src/"
	goSourcePath := goPath + username
	if addIdeaWatchers {
		os.Mkdir(goSourcePath+"/.idea", os.ModePerm)
		content, _ := file.FileGetContents(TemplatePath() + "/bot/watcherTasks.xml")
		file.FilePutContents(goSourcePath+"/.idea/"+"watcherTasks.xml", content)
	}
	apiBytes, _ := file.FileGetContents(TemplatePath() + "/bot/api.temp")

	databaseBytes, _ := file.FileGetContents(TemplatePath() + "/bot/database.temp")
	keyboardsBytes, _ := file.FileGetContents(TemplatePath() + "/bot/keyboards.temp")
	keyboardsBytes = []byte(strings.Replace(string(keyboardsBytes), "%BOTUSERNAME%", username, -1))
	languageBytes, _ := file.FileGetContents(TemplatePath() + "/bot/language.temp")
	engineBytes, _ := file.FileGetContents(TemplatePath() + "/bot/engine.temp") // %BOTUSERNAME% %BOTUSERNAME_CAPS%
	engineBytes = []byte(strings.Replace(string(engineBytes), "%BOTUSERNAME%", username, -1))
	engineBytes = []byte(strings.Replace(string(engineBytes), "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1))
	methodsBytes, _ := file.FileGetContents(TemplatePath() + "/bot/methods.temp") // %BOTUSERNAME% %BOTUSERNAME_CAPS%
	methodsBytes = []byte(strings.Replace(string(methodsBytes), "%BOTUSERNAME%", username, -1))
	methodsBytes = []byte(strings.Replace(string(methodsBytes), "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1))
	configBytes, _ := file.FileGetContents(TemplatePath() + "/bot/config.temp") // %BOTUSERNAME% %BOTTOKEN%
	configBytes = []byte(strings.Replace(string(configBytes), "%BOTUSERNAME%", username, -1))
	configBytes = []byte(strings.Replace(string(configBytes), "%BOTTOKEN%", token, -1))
	mainBytes, _ := file.FileGetContents(TemplatePath() + "/bot/main.temp") // %BOTUSERNAME%
	mainBytes = []byte(strings.Replace(string(mainBytes), "%BOTUSERNAME%", username, -1))
	mainBytes = []byte(strings.Replace(string(mainBytes), "%BOTUSERNAME_CAPS%", strings.ToUpper(username), -1))
	userBytes, _ := file.FileGetContents(TemplatePath() + "/bot/user.temp")
	persianBytes, _ := file.FileGetContents(TemplatePath() + "/bot/persian.temp")
	englishBytes, _ := file.FileGetContents(TemplatePath() + "/bot/english.temp")

	configPath := goSourcePath + "/configs/"
	funcsPath := goSourcePath + "/funcs/"
	langPath := goSourcePath + "/lang/"
	modelPath := goSourcePath + "/models/"
	os.Mkdir(goSourcePath, os.ModePerm)
	os.Mkdir(goSourcePath+"/configs/", os.ModePerm)
	os.Mkdir(goSourcePath+"/funcs/", os.ModePerm)
	os.Mkdir(goSourcePath+"/lang/", os.ModePerm)
	os.Mkdir(goSourcePath+"/models/", os.ModePerm)
	if err := file.FilePutContents(configPath+"main.go", configBytes); err != nil {
		fmt.Println(err)
		return
	}
	file.FilePutContents(funcsPath+"engine.go", engineBytes)
	file.FilePutContents(funcsPath+"methods.go", methodsBytes)
	file.FilePutContents(funcsPath+"keyboards.go", keyboardsBytes)
	file.FilePutContents(modelPath+"api.go", apiBytes)
	file.FilePutContents(modelPath+"database.go", databaseBytes)
	file.FilePutContents(modelPath+"user.go", userBytes)
	file.FilePutContents(langPath+"language.go", languageBytes)
	file.FilePutContents(langPath+"persian.go", persianBytes)
	file.FilePutContents(langPath+"english.go", englishBytes)
	file.FilePutContents(goSourcePath+"/main.go", mainBytes)
	exec.Command("go fmt ", "-w", goSourcePath).Output()
	exec.Command("goimports", "-w", goSourcePath).Output()
	fmt.Println(fmt.Sprintf("Bot %s Created at %s", username, goSourcePath))
}
