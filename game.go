package goadventure

type Game struct {
}

type GameState struct {
	game       *Game
	screenName string
}

func (game Game) GetStateForUser(screenName string) (gameState *GameState) {
	gameState = &GameState{&game, screenName}
	return
}

func (gameState GameState) UpdateState(command string) (response string) {
	return "You are in an empty room. There are doors to the North and South of you."
}
