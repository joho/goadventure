package goadventure

type PersistentStorageEngine struct {
	scenes        map[uint64]*Scene
	tweetsHandled map[uint64]string
}

func NewPersistentStorageEngine() StorageEngine {
	return &PersistentStorageEngine{
		map[uint64]*Scene{},
		map[uint64]string{},
	}
}

func (repo *PersistentStorageEngine) GetCurrentSceneForUser(twitterUserId uint64) *Scene {
	return repo.scenes[twitterUserId]
}

func (repo *PersistentStorageEngine) SetCurrentSceneForUser(twitterUserId uint64, scene *Scene) {
	repo.scenes[twitterUserId] = scene
}

func (repo *PersistentStorageEngine) TweetAlreadyHandled(tweetId uint64) bool {
	_, present := repo.tweetsHandled[tweetId]
	return present
}

func (repo *PersistentStorageEngine) StoreTweetHandled(tweetId uint64, tweetContents string) {
	repo.tweetsHandled[tweetId] = tweetContents
}
