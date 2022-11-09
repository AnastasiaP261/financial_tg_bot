package messages

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type status string

const (
	keySuffix = "status" // чтобы не перепутать значения, которые могут лежать в редисе с таким же ключом и не относиться к статусам

	statusNonExistentCategory status = "msgNonExistentCategory"
)

type userInfo struct {
	Status  status `json:"status"`  // статус юзера
	Command string `json:"command"` // "замороженная" команда
}

func createKey(userID int64) string {
	return strconv.FormatInt(userID, 10) + keySuffix
}

// getUserInfo получить информацию о статусе пользователя
func (m *Model) getUserInfo(ctx context.Context, userID int64) (userInfo, error) {
	decKey := createKey(userID)

	res, err := m.statusStore.GetString(ctx, decKey)
	if err != nil {
		return userInfo{}, errors.Wrap(err, "statusStore.GetString")
	}

	rawJson, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		return userInfo{}, errors.Wrap(err, "base64.StdEncoding.DecodeString")
	}

	var info userInfo
	if err = json.Unmarshal(rawJson, &info); err != nil {
		return userInfo{}, errors.Wrap(err, "unmarshalling err")
	}

	return info, nil
}

// setUserInfo установить новый статус пользователю. Для удаления статуса отправить пустой userInfo
func (m *Model) setUserInfo(ctx context.Context, userID int64, info userInfo) error {
	decKey := createKey(userID)

	if info.Status == "" {
		if err := m.statusStore.Delete(ctx, decKey); err != nil {
			return errors.Wrap(err, "statusStore.Delete")
		}
		return nil
	}

	bytes, err := json.Marshal(info)
	if err != nil {
		return errors.Wrap(err, "marshalling err")
	}

	decVal := base64.StdEncoding.EncodeToString(bytes)

	if err = m.statusStore.SetString(ctx, decKey, decVal); err != nil {
		return errors.Wrap(err, "statusStore.SetString")
	}

	return nil
}
