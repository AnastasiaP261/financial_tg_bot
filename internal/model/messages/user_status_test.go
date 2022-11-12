package messages

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// base64('{"status":"someStatus","command":"/command 123"}') = 'eyJzdGF0dXMiOiJzb21lU3RhdHVzIiwiY29tbWFuZCI6Ii9jb21tYW5kIDEyMyJ9'

func TestModel_getUserInfo(t *testing.T) {
	ctx := context.Background()

	_, _, statusStore := mocksUp(t)
	model := New(nil, nil, statusStore)

	statusStore.EXPECT().GetString(ctx, "123status").Return("eyJzdGF0dXMiOiJzb21lU3RhdHVzIiwiY29tbWFuZCI6Ii9jb21tYW5kIDEyMyJ9", nil)

	info, err := model.getUserInfo(ctx, 123)

	assert.NoError(t, err)
	assert.Equal(t, userInfo{
		Status:  "someStatus",
		Command: "/command 123",
	}, info)
}

func TestModel_setUserInfo(t *testing.T) {
	t.Run("установка статуса", func(t *testing.T) {
		ctx := context.Background()

		_, _, statusStore := mocksUp(t)
		model := New(nil, nil, statusStore)

		statusStore.EXPECT().SetString(ctx, "123status", "eyJzdGF0dXMiOiJzb21lU3RhdHVzIiwiY29tbWFuZCI6Ii9jb21tYW5kIDEyMyJ9").Return(nil)

		err := model.setUserInfo(ctx, 123, userInfo{
			Status:  "someStatus",
			Command: "/command 123",
		})

		assert.NoError(t, err)
	})

	t.Run("удаление статуса", func(t *testing.T) {
		ctx := context.Background()

		_, _, statusStore := mocksUp(t)
		model := New(nil, nil, statusStore)

		statusStore.EXPECT().Delete(ctx, "123status").Return(nil)

		err := model.setUserInfo(ctx, 123, userInfo{})

		assert.NoError(t, err)
	})
}
