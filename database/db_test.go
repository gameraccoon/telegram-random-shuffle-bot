package database

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	testDbPath = "./testDb.db"
)

func dropDatabase(fileName string) {
	os.Remove(fileName)
}

func clearDb() {
	dropDatabase(testDbPath)
}

func connectDb(t *testing.T) *Database {
	assert := require.New(t)
	db := &Database{}

	err := db.Connect(testDbPath)
	if err != nil {
		assert.Fail("Problem with creation db connection:" + err.Error())
		return nil
	}
	return db
}

func createDbAndConnect(t *testing.T) *Database {
	clearDb()
	return connectDb(t)
}

func TestConnection(t *testing.T) {
	assert := require.New(t)
	dropDatabase(testDbPath)

	db := &Database{}

	assert.False(db.IsConnectionOpened())

	err := db.Connect(testDbPath)
	defer dropDatabase(testDbPath)
	if err != nil {
		assert.Fail("Problem with creation db connection:" + err.Error())
		return
	}

	assert.True(db.IsConnectionOpened())

	db.Disconnect()

	assert.False(db.IsConnectionOpened())
}

func TestSanitizeString(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	testText := "text'test''test\"test\\"

	db.SetDatabaseVersion(testText)
	assert.Equal(testText, db.GetDatabaseVersion())
}

func TestDatabaseVersion(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}

	{
		version := db.GetDatabaseVersion()
		assert.Equal(latestVersion, version)
	}

	db.SetDatabaseVersion("1.2")

	{
		version := db.GetDatabaseVersion()
		assert.Equal("1.2", version)
	}

	db.SetDatabaseVersion("1.4")
	db.Disconnect()

	{
		db = connectDb(t)
		version := db.GetDatabaseVersion()
		assert.Equal("1.4", version)
		db.Disconnect()
	}
}

func TestGetUserId(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	var chatId1 int64 = 321
	var chatId2 int64 = 123

	id1 := db.GetUserId(chatId1)
	id2 := db.GetUserId(chatId1)
	id3 := db.GetUserId(chatId2)

	assert.Equal(id1, id2)
	assert.NotEqual(id1, id3)

	assert.Equal(chatId1, db.GetUserChatId(id1))
	assert.Equal(chatId2, db.GetUserChatId(id3))
}

func TestCreateAndRemoveList(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	var chatId int64 = 123
	userId := db.GetUserId(chatId)

	{
		ids, texts := db.GetUserLists(userId)
		assert.Equal(0, len(ids))
		assert.Equal(0, len(texts))
	}
	
	assert.Equal(false, db.IsListExists(0))
	assert.Equal(false, db.IsListExists(1))

	listId := db.CreateList(userId, "testlist")
	assert.Equal(true, db.IsListExists(listId))
	{
		ids, texts := db.GetUserLists(userId)
		assert.Equal(1, len(ids))
		assert.Equal(1, len(texts))
		if len(ids) > 0 && len(texts) > 0 {
			assert.Equal(listId, ids[0])
			assert.Equal("testlist", texts[0])
			assert.Equal("testlist", db.GetListName(ids[0]))
		}
	}

	db.DeleteList(listId)
	assert.Equal(false, db.IsListExists(listId))
	{
		ids, texts := db.GetUserLists(userId)
		assert.Equal(0, len(ids))
		assert.Equal(0, len(texts))
	}
}

func TestAddingAndRemovingElementsToList(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	var chatId int64 = 123
	userId := db.GetUserId(chatId)
	db.CreateList(userId, "testlist")

	ids, _ := db.GetUserLists(userId)
	if len(ids) > 0 {
		listId := ids[0]
		{
			ids, texts := db.GetListItems(listId)
			assert.Equal(0, len(ids))
			assert.Equal(0, len(texts))
		}

		db.AddItemsToList(listId, []string{"one", "two"})
		{
			ids, texts := db.GetListItems(listId)
			assert.Equal(2, len(ids))
			assert.Equal(2, len(texts))
			if len(ids) > 1 && len(texts) > 1 {
				assert.Equal("one", texts[0])
				assert.Equal("two", texts[1])

				db.RemoveItem(ids[0])
			}
		}

		{
			ids, texts := db.GetListItems(listId)
			assert.Equal(1, len(ids))
			assert.Equal(1, len(texts))
			if len(ids) > 0 && len(texts) > 0 {
				assert.Equal("two", texts[0])
			}
		}
	}
}

func TestSetLastQuestion(t *testing.T) {
		assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	var chatId int64 = 123
	userId := db.GetUserId(chatId)

	{
		ids, texts := db.GetUserLists(userId)
		assert.Equal(0, len(ids))
		assert.Equal(0, len(texts))
	}

	listId := db.CreateList(userId, "testlist")
	db.AddItemsToList(listId, []string{"one", "two"})
	
	ids, _ := db.GetUserLists(userId)
	if len(ids) > 0 {
		listId := ids[0]
		{
			lastItemId, _ := db.GetLastItem(listId)
			assert.Equal(int64(-1), lastItemId)
		}

		itemId := int64(1)
		db.SetLastItem(listId, itemId)
		{
			lastItemId, lastItemText := db.GetLastItem(listId)
			assert.Equal(itemId, lastItemId)
			itemIds, itemTexts := db.GetListItems(listId)
			
			var foundItemIndex int = -1
			for i, id := range itemIds {
				if (id == itemId) {
					foundItemIndex = i
					break
				}
			}
			assert.NotEqual(-1, foundItemIndex)
			if foundItemIndex != -1 {
				assert.Equal(itemTexts[foundItemIndex], lastItemText)
			}
		}
	}
}
