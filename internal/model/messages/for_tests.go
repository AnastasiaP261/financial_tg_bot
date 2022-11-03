package messages

import (
	"testing"

	"github.com/golang/mock/gomock"
	mocks "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages/_mocks"
)

func mocksUp(t *testing.T) (*mocks.MockMessageSender, *mocks.MockPurchasesModel, *mocks.MockStatusStore) {
	ctrl := gomock.NewController(t)
	return mocks.NewMockMessageSender(ctrl), mocks.NewMockPurchasesModel(ctrl), mocks.NewMockStatusStore(ctrl)
}
