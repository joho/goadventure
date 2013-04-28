package goadventure

import (
	"fmt"
	"github.com/peterbourgon/diskv"
	"strconv"
)

type PersistentStorageEngine struct {
	tweetStore     *diskv.Diskv
	gameStateStore *diskv.Diskv
}

func NewPersistentStorageEngine() StorageEngine {
	storageEngine := new(PersistentStorageEngine)

	flatTransform := func(s string) []string { return []string{} }
	storageEngine.tweetStore = diskv.New(diskv.Options{
		BasePath:     "storage/tweets",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	storageEngine.gameStateStore = diskv.New(diskv.Options{
		BasePath:     "storage/game-states",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	return storageEngine
}

func (repo *PersistentStorageEngine) GetCurrentSceneIdForUser(twitterUserId uint64) (int, error) {
	rawValue, err := repo.gameStateStore.Read(strconv.FormatUint(twitterUserId, 10))
	if err != nil {
		fmt.Printf("Read \"%v\" as rawValue with \"%v\" as err\n", rawValue, err)
		return -1, err
	}
	value, err := strconv.Atoi(string(rawValue))
	return value, err
}

func (repo *PersistentStorageEngine) SetCurrentSceneIdForUser(twitterUserId uint64, sceneId int) {
	key := strconv.FormatUint(twitterUserId, 10)
	value := strconv.Itoa(sceneId)
	repo.gameStateStore.Write(key, []byte(value))
}

func (repo *PersistentStorageEngine) TweetAlreadyHandled(tweetId uint64) bool {
	_, err := repo.tweetStore.Read(strconv.FormatUint(tweetId, 10))
	return err == nil
}

func (repo *PersistentStorageEngine) StoreTweetHandled(tweetId uint64, tweetContents string) {
	key := strconv.FormatUint(tweetId, 10)
	repo.tweetStore.Write(key, []byte(tweetContents))
}
