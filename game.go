package goadventure

import "strings"

type Game struct {
	gameStateRepo GameStateRepo
	openingScene  *Scene
}

type Scene struct {
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

func CreateGame() *Game {
	roomOne := &Scene{"Welcome to room one. You can go north.", nil}
	roomTwo := &Scene{"You're in room two. You can go south or west", nil}
	roomThree := &Scene{"You're in room three. You can go east", nil}

	roomOne.LinkSceneViaCommand(roomTwo, Command{"go", "north"})

	roomTwo.LinkSceneViaCommand(roomOne, Command{"go", "south"})
	roomTwo.LinkSceneViaCommand(roomThree, Command{"go", "west"})

	roomThree.LinkSceneViaCommand(roomTwo, Command{"go", "east"})

	emptySceneMap := map[uint64]*Scene{}
	gameStateRepo := &InMemoryGameStateRepo{emptySceneMap}

	return &Game{
		gameStateRepo,
		roomOne,
	}
}

func (game *Game) Play(twitterUserId uint64, rawCommand string) string {
	var (
		currentScene *Scene
		nextScene    *Scene
		responseText string
	)

	currentScene = game.gameStateRepo.GetCurrentSceneForUser(twitterUserId)
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
		game.gameStateRepo.SetCurrentSceneForUser(twitterUserId, nextScene)
		responseText = nextScene.Description
	} else {
		responseText = "Sorry Dave, I can't let you do that"
	}

	return responseText
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

// Temporary storage for dev
type InMemoryGameStateRepo struct {
	scenes map[uint64]*Scene
}

func (repo *InMemoryGameStateRepo) GetCurrentSceneForUser(twitterUserId uint64) *Scene {
	return repo.scenes[twitterUserId]
}

func (repo *InMemoryGameStateRepo) SetCurrentSceneForUser(twitterUserId uint64, scene *Scene) {
	repo.scenes[twitterUserId] = scene
}
