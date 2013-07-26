package goadventure

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

type Game struct {
	storageEngine StorageEngine
	openingScene  *Scene
}

type Scene struct {
	key         string
	Description string
	choices     []*Choice
}

type Command struct {
	Verb    string
	Subject string
}

type Choice struct {
	command Command
	scene   *Scene
}

func CreateGame(storageEngine StorageEngine) *Game {

	lines := make(chan string)
	scenes := make(map[string]*Scene)

	go parseLinesIntoSceneTree(scenes, lines)

	readLinesFromScript("script", lines)

	return &Game{
		storageEngine,
		scenes["start"],
	}
}

func readLinesFromScript(scriptPath string, lines chan string) {
	file, err := os.Open(scriptPath)
	if err != nil {
		log.Fatal("Error opening script file.")
	}
	lineReader := bufio.NewReaderSize(file, 20)
	for line, isPrefix, e := lineReader.ReadLine(); e == nil; line, isPrefix, e = lineReader.ReadLine() {
		fullLine := string(line)
		if isPrefix {
			for {
				line, isPrefix, _ = lineReader.ReadLine()
				fullLine += string(line)
				if !isPrefix {
					break
				}
			}
		}
		// add a line to the game/parse
		lines <- fullLine
	}
	close(lines)
}

func parseLinesIntoSceneTree(scenes map[string]*Scene, lines chan string) {
	/*
		Script data format is as below.
		First room must have key of "start"
		All rooms must be defined before they can be referenced in a choice

		===
		key
		description (less than 120 chars)
		- verb noun : key (of next room)

		example:

		===
		start
		Welcome to room one. You can go north.
		- go north : north_room
		===
		north_room
		Welcome to room two. You can go south or west.
		- go south : start
		- go west : west_room
		===
		west_room
		Welcome to room three. You can go east.
		- go east : north_room
	*/

	maxDescriptionLength := 120

	var scene *Scene
	for line := range lines {
		if line == "===" {
			if scene != nil {
				scenes[scene.key] = scene
			}
			scene = new(Scene)
		} else {
			if strings.HasPrefix("-", line) {
				// populate choices
				segments := strings.Split(line, ":")

				sceneKey := strings.ToLower(strings.Trim(segments[1], " "))

				rawCommandWords := strings.Split(strings.Trim(segments[0], " -"), " ")
				command := Command{rawCommandWords[0], rawCommandWords[1]}

				scene.LinkSceneViaCommand(scenes[sceneKey], command)
			} else {
				// key has not been set ("" is string zero value)
				if scene.key == "" {
					scene.key = strings.ToLower(line)
				} else {
					scene.Description = line
					length := len(line)
					if length > maxDescriptionLength {
						log.Fatalf("Scene \"%v\" description exceeds max length %v with %v chars", scene.key, maxDescriptionLength, length)
					}
				}
			}
		}
	}
}

func (game *Game) Play(twitterUserId uint64, rawCommand string) string {
	var (
		currentScene *Scene
		nextScene    *Scene
		responseText string
	)

	currentScene = game.GetCurrentSceneForUser(twitterUserId)
	if currentScene == nil {
		// kick off the adventure
		nextScene = game.openingScene
	} else {
		nextScene = currentScene.DoSomethingMagical(rawCommand)
	}

	if nextScene != nil {
		game.SetCurrentSceneForUser(twitterUserId, nextScene)
		responseText = nextScene.Description
	} else {
		responseText = "Sorry Dave, I can't let you do that"
	}

	return responseText
}

func (game *Game) SetCurrentSceneForUser(twitterUserId uint64, scene *Scene) {
	game.storageEngine.SetCurrentSceneKeyForUser(twitterUserId, scene.key)
}

func (game *Game) GetCurrentSceneForUser(twitterUserId uint64) *Scene {
	sceneKey, _ := game.storageEngine.GetCurrentSceneKeyForUser(twitterUserId)
	if sceneKey == "" {
		return nil
	}

	// I have a sneaking suspicion that this recursive tree walk
	// is not exactly *idiomatic* go
	visitedChoices := map[*Choice]bool{}
	var walker func(*Scene) *Scene
	walker = func(rootScene *Scene) (desiredScene *Scene) {
		if rootScene.key == sceneKey {
			desiredScene = rootScene
		} else {
			for _, choice := range rootScene.choices {
				if !visitedChoices[choice] {
					visitedChoices[choice] = true
					desiredScene = walker(choice.scene)
					if desiredScene != nil {
						break
					}
				}
			}
		}
		return
	}

	return walker(game.openingScene)
}

func (scene *Scene) LinkSceneViaCommand(nextScene *Scene, command Command) {
	choice := &Choice{command, nextScene}
	scene.choices = append(scene.choices, choice)
}

func (scene *Scene) DoSomethingMagical(rawCommand string) (nextScene *Scene) {
	command, err := scene.parseCommand(rawCommand)
	if err != nil {
		return
	}

	// remind the user where they are if they want to look around
	if command.Verb == "look" && command.Subject == "around" {
		nextScene = scene
		return
	}

	// check all choices from the current scene for a command match
	for _, choice := range scene.choices {
		if choice.command == command {
			nextScene = choice.scene
			break
		}
	}
	return
}

func (scene *Scene) parseCommand(rawCommand string) (command Command, err error) {
	// should usually be of format "@goadventuregame go north"
	words := strings.Fields(rawCommand)
	if len(words) < 3 {
		log.Printf("Bad command received: \"%v\" on scene \"%v\"\n", rawCommand, scene.Description)
		err = errors.New("Not enough words in command")
	} else {
		command = Command{words[1], words[2]}
	}
	return
}
