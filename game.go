package goadventure

import "strings"

type Game struct {
	StateRepo
	openingScene *Scene
}

type StateRepo struct {
	scenes map[uint64]*Scene
}

type Scene struct {
	Description string
	choices     []*Choice
}

type Choice struct {
	command Command
	scene   *Scene
}

type Command struct {
	Verb    string
	Subject string
}

func CreateGame() *Game {
	roomOne := &Scene{"Welcome to room one. You can go north.", nil}
	roomTwo := &Scene{"You're in room two. You can go south", nil}

	choiceRoomOneToTwo := &Choice{
		Command{"go", "north"},
		roomTwo,
	}
	roomOne.choices = append(roomOne.choices, choiceRoomOneToTwo)

	choiceRoomTwoToOne := &Choice{
		Command{"go", "south"},
		roomOne,
	}
	roomTwo.choices = append(roomTwo.choices, choiceRoomTwoToOne)

	emptySceneMap := map[uint64]*Scene{}
	return &Game{
		StateRepo{emptySceneMap},
		roomTwo,
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
	return game.openingScene
}

func (repo *StateRepo) CurrentSceneForUser(twitterUserId uint64) *Scene {
	return repo.scenes[twitterUserId]
}

func (repo *StateRepo) SetCurrentSceneForUser(twitterUserId uint64, scene *Scene) {
	repo.scenes[twitterUserId] = scene
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

func parseCommand(rawCommand string) Command {
	// should usually be of format "@goadventure go north"
	words := strings.Fields(rawCommand)

	return Command{words[1], words[2]}
}
