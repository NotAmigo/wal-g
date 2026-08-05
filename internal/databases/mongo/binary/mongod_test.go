package binary

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestMongodClientOptionsDoNotLimitBackupCursorLifetime(t *testing.T) {
	timeout := 10 * time.Minute
	clientOptions := mongodClientOptions("backup", "mongodb://localhost:27017", timeout)

	assert.Equal(t, timeout, *clientOptions.ServerSelectionTimeout)
	assert.Equal(t, timeout, *clientOptions.ConnectTimeout)
	assert.Nil(t, clientOptions.Timeout)
	assert.Equal(t, "backup", *clientOptions.AppName)
	assert.True(t, *clientOptions.Direct)
	assert.False(t, *clientOptions.RetryReads)
}

func TestBackupCursorErrorIsRetried(t *testing.T) {
	assert.True(t, backupCursorErrorIsRetried(errors.New("(Location50915) backup cursor is already open")))
	assert.True(t, backupCursorErrorIsRetried(errors.New("(BackupCursorOpenConflictWithCheckpoint) checkpoint in progress")))
	assert.False(t, backupCursorErrorIsRetried(errors.New("connection refused")))
}

func TestMakeBsonRsMembers(t *testing.T) {
	assert.Equal(t, bson.A{}, makeBsonRsMembers(RsConfig{}))
	assert.Equal(t, bson.A{bson.M{"_id": 0, "host": "localhost:1234"}}, makeBsonRsMembers(RsConfig{
		RsMembers:   []string{"localhost:1234"},
		RsMemberIDs: []int{0},
	}))
	assert.Equal(t,
		bson.A{
			bson.M{"_id": 0, "host": "localhost:1234"},
			bson.M{"_id": 1, "host": "localhost:5678"},
			bson.M{"_id": 2, "host": "remotehost:9876"},
		},
		makeBsonRsMembers(RsConfig{
			RsName:      "",
			RsMembers:   []string{"localhost:1234", "localhost:5678", "remotehost:9876"},
			RsMemberIDs: []int{0, 1, 2},
		}))
	assert.Equal(t,
		bson.A{
			bson.M{"_id": 4, "host": "localhost:1234"},
			bson.M{"_id": 5, "host": "localhost:5678"},
			bson.M{"_id": 0, "host": "remotehost:9876"},
		},
		makeBsonRsMembers(RsConfig{
			RsName:      "",
			RsMembers:   []string{"localhost:1234", "localhost:5678", "remotehost:9876"},
			RsMemberIDs: []int{4, 5, 0},
		}))
}
