package goadventure

import (
	"github.com/peterbourgon/diskv"
	"log"
	"strconv"
)

type PersistentStorageEngine struct {
	tweetStore     *diskv.Diskv
	gameStateStore *diskv.Diskv
}

func NewPersistentStorageEngine() StorageEngine {
	storageEngine := new(PersistentStorageEngine)

	// TODO set up a folder structure for < 1k entries per folder
	// actual TODO set up a proper datastore
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
	rawValue, err := repo.gameStateStore.Read(repo.translateKey(twitterUserId))
	if err != nil {
		log.Printf("Read \"%v\" as rawValue with \"%v\" as err\n", rawValue, err)
		return -1, err
	}
	value, err := strconv.Atoi(string(rawValue))
	return value, err
}

func (repo *PersistentStorageEngine) SetCurrentSceneIdForUser(twitterUserId uint64, sceneId int) {
	value := strconv.Itoa(sceneId)
	repo.gameStateStore.Write(repo.translateKey(twitterUserId), []byte(value))
}

func (repo *PersistentStorageEngine) TweetAlreadyHandled(tweetId uint64) bool {
	_, err := repo.tweetStore.Read(repo.translateKey(tweetId))
	return err == nil
}

func (repo *PersistentStorageEngine) StoreTweetHandled(tweetId uint64, tweetContents string) {
	repo.tweetStore.Write(repo.translateKey(tweetId), []byte(tweetContents))
}

func (_ *PersistentStorageEngine) translateKey(rawKey uint64) string {
	return strconv.FormatUint(rawKey, 10)
}
