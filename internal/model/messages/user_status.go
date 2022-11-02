package messages

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

type status string

const (
	statusNonExistentCategory status = "nonExistentCategory"
	statusNotAddedCategory    status = "notAddedCategory"
)

type userInfo struct {
	Status  status `json:"status"`  // статус юзера
	Command string `json:"command"` // "замороженная" команда
}

// getUserInfo получить информацию о статусе пользователя
func (m *Model) getUserInfo(ctx context.Context, userID int64) (userInfo, error) {
	decKey := strconv.FormatInt(userID, 10)

	res, err := m.statusStore.Get(ctx, decKey)
	if err != nil {
		return userInfo{}, errors.Wrap(err, "statusStore.Get")
	}

	fmt.Println("### redis val", res)
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
	decKey := strconv.FormatInt(userID, 10)

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

	if err = m.statusStore.Set(ctx, decKey, decVal); err != nil {
		return errors.Wrap(err, "statusStore.Set")
	}

	return nil
}
