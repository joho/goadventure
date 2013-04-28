package goadventure

import "strings"

type Game struct {
	storageEngine StorageEngine
	openingScene  *Scene
}

type Scene struct {
	id          int
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
	roomOne := &Scene{1, "Welcome to room one. You can go north.", nil}
	roomTwo := &Scene{2, "You're in room two. You can go south or west", nil}
	roomThree := &Scene{3, "You're in room three. You can go east", nil}

	roomOne.LinkSceneViaCommand(roomTwo, Command{"go", "north"})

	roomTwo.LinkSceneViaCommand(roomOne, Command{"go", "south"})
	roomTwo.LinkSceneViaCommand(roomThree, Command{"go", "west"})

	roomThree.LinkSceneViaCommand(roomTwo, Command{"go", "east"})

	return &Game{
		storageEngine,
		roomOne,
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
		// should usually be of format "@goadventure go north"
		words := strings.Fields(rawCommand)
		command := Command{words[1], words[2]}

		nextScene = currentScene.DoSomethingMagical(command)
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
	game.storageEngine.SetCurrentSceneIdForUser(twitterUserId, scene.id)
}

func (game *Game) GetCurrentSceneForUser(twitterUserId uint64) *Scene {
	sceneId := game.storageEngine.GetCurrentSceneIdForUser(twitterUserId)

	// I have a sneaking suspicion that this recursive tree walk
	// is not exactly *idiomatic* go
	visitedChoices := map[*Choice]bool{}
	var walker func(*Scene) *Scene
	walker = func(rootScene *Scene) (desiredScene *Scene) {
		if rootScene.id == sceneId {
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

func (scene *Scene) DoSomethingMagical(command Command) (nextScene *Scene) {
	for _, choice := range scene.choices {
		if choice.command == command {
			nextScene = choice.scene
			break
		}
	}
	return
}
