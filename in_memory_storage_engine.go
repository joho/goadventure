package goadventure

// Temporary storage for dev
type InMemoryStorageEngine struct {
	scenes        map[uint64]string
	tweetsHandled map[uint64]string
}

func NewInMemoryStorageEngine() StorageEngine {
	return &InMemoryStorageEngine{
		map[uint64]string{},
		map[uint64]string{},
	}
}

func (repo *InMemoryStorageEngine) GetCurrentSceneKeyForUser(twitterUserId uint64) (string, error) {
	return repo.scenes[twitterUserId], nil
}

func (repo *InMemoryStorageEngine) SetCurrentSceneKeyForUser(twitterUserId uint64, sceneKey string) {
	repo.scenes[twitterUserId] = sceneKey
}

func (repo *InMemoryStorageEngine) TweetAlreadyHandled(tweetId uint64) bool {
	_, present := repo.tweetsHandled[tweetId]
	return present
}

func (repo *InMemoryStorageEngine) StoreTweetHandled(tweetId uint64, tweetContents string) {
	repo.tweetsHandled[tweetId] = tweetContents
}
