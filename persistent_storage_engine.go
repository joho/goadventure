package goadventure

import (
	"github.com/peterbourgon/diskv"
)

type PersistentStorageEngine struct {
	diskvStore *diskv.Diskv
}

func NewPersistentStorageEngine() StorageEngine {
	storageEngine := new(PersistentStorageEngine)

	flatTransform := func(s string) []string { return []string{} }
	storageEngine.diskvStore = diskv.New(diskv.Options{
		BasePath:     "diskv-store",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	return storageEngine
}

func (repo *PersistentStorageEngine) GetCurrentSceneIdForUser(twitterUserId uint64) int {
	return 1
}

func (repo *PersistentStorageEngine) SetCurrentSceneIdForUser(twitterUserId uint64, sceneId int) {
}

func (repo *PersistentStorageEngine) TweetAlreadyHandled(tweetId uint64) bool {
	return false
}

func (repo *PersistentStorageEngine) StoreTweetHandled(tweetId uint64, tweetContents string) {
}
