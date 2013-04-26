package goadventure

// Temporary storage for dev
type InMemoryStorageEngine struct {
	scenes        map[uint64]*Scene
	tweetsHandled map[uint64]string
}

func CreateTweetLogger() TweetRepo {
	return &InMemoryStorageEngine{
		map[uint64]*Scene{},
		map[uint64]string{},
	}
}

func CreateGameStateRepo() GameStateRepo {
	return &InMemoryStorageEngine{
		map[uint64]*Scene{},
		map[uint64]string{},
	}
}

func (repo *InMemoryStorageEngine) GetCurrentSceneForUser(twitterUserId uint64) *Scene {
	return repo.scenes[twitterUserId]
}

func (repo *InMemoryStorageEngine) SetCurrentSceneForUser(twitterUserId uint64, scene *Scene) {
	repo.scenes[twitterUserId] = scene
}

func (repo *InMemoryStorageEngine) TweetAlreadyHandled(tweetId uint64) bool {
	_, present := repo.tweetsHandled[tweetId]
	return present
}

func (repo *InMemoryStorageEngine) StoreTweetHandled(tweetId uint64, tweetContents string) {
	repo.tweetsHandled[tweetId] = tweetContents
}
