package goadventure

import "strings"

type Game struct {
	StateRepo
}

type StateRepo struct {
	scenes map[uint64]*Scene
}

type Scene struct {
	Description string
}

type Command struct {
	Verb    string
	Subject string
}

func CreateGame() *Game {
	return &Game{
		StateRepo{
			map[uint64]*Scene{},
		},
	}
}

func (game *Game) Play(twitterUserId uint64, rawCommand string) string {
	var (
		currentScene *Scene
		nextScene    *Scene
		responseText string
	)

	currentScene = game.CurrentSceneForUser(twitterUserId)
	if currentScene == nil {
		// kick off the adventure
		nextScene = game.OpeningScene()
	} else {
		command := parseCommand(rawCommand)
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

func (game *Game) OpeningScene() *Scene {
	return &Scene{"Welcome. You are in an empty room. There is a door to the north."}
}

func (repo *StateRepo) CurrentSceneForUser(twitterUserId uint64) *Scene {
	return repo.scenes[twitterUserId]
}

func (repo *StateRepo) SetCurrentSceneForUser(twitterUserId uint64, scene *Scene) {
	repo.scenes[twitterUserId] = scene
}

func (scene *Scene) DoSomethingMagical(command Command) (nextScene *Scene) {
	if command.Verb == "go" {
		if command.Subject == "north" && strings.Contains(scene.Description, "north") {
			nextScene = &Scene{"You're in an empty room. There is a door to the south."}
		} else if command.Subject == "south" && strings.Contains(scene.Description, "south") {
			nextScene = &Scene{"You're in an empty room. There is a door to the north."}
		}
	}
	return
}

func parseCommand(rawCommand string) Command {
	// should usually be of format "@goadventure go north"
	words := strings.Fields(rawCommand)

	return Command{words[1], words[2]}
}
