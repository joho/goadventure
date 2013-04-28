package goadventure

// Temporary storage for dev
type InMemoryStorageEngine struct {
	scenes        map[uint64]int
	tweetsHandled map[uint64]string
}

func NewInMemoryStorageEngine() StorageEngine {
	return &InMemoryStorageEngine{
		map[uint64]int{},
		map[uint64]string{},
	}
}

func (repo *InMemoryStorageEngine) GetCurrentSceneIdForUser(twitterUserId uint64) int {
	return repo.scenes[twitterUserId]
}

func (repo *InMemoryStorageEngine) SetCurrentSceneIdForUser(twitterUserId uint64, sceneId int) {
	repo.scenes[twitterUserId] = sceneId
}

func (repo *InMemoryStorageEngine) TweetAlreadyHandled(tweetId uint64) bool {
	_, present := repo.tweetsHandled[tweetId]
	return present
}

func (repo *InMemoryStorageEngine) StoreTweetHandled(tweetId uint64, tweetContents string) {
	repo.tweetsHandled[tweetId] = tweetContents
}
