package goadventure

type Game struct {
}

func (game Game) Play(twitterUserId uint64, command string) string {
	return "You are in an empty room. There are doors to the North and South of you."
}
