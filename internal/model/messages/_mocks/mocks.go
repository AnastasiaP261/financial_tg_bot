// Code generated by MockGen. DO NOT EDIT.
// Source: internal/model/messages/incoming_msg.go

// Package mock_messages is a generated GoMock package.
package mock_messages

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	purchases "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

// MockMessageSender is a mock of MessageSender interface.
type MockMessageSender struct {
	ctrl     *gomock.Controller
	recorder *MockMessageSenderMockRecorder
}

// MockMessageSenderMockRecorder is the mock recorder for MockMessageSender.
type MockMessageSenderMockRecorder struct {
	mock *MockMessageSender
}

// NewMockMessageSender creates a new mock instance.
func NewMockMessageSender(ctrl *gomock.Controller) *MockMessageSender {
	mock := &MockMessageSender{ctrl: ctrl}
	mock.recorder = &MockMessageSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageSender) EXPECT() *MockMessageSenderMockRecorder {
	return m.recorder
}

// SendImage mocks base method.
func (m *MockMessageSender) SendImage(img []byte, chatID int64, userName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendImage", img, chatID, userName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendImage indicates an expected call of SendImage.
func (mr *MockMessageSenderMockRecorder) SendImage(img, chatID, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendImage", reflect.TypeOf((*MockMessageSender)(nil).SendImage), img, chatID, userName)
}

// SendMessage mocks base method.
func (m *MockMessageSender) SendMessage(text string, userID int64, userName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", text, userID, userName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockMessageSenderMockRecorder) SendMessage(text, userID, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockMessageSender)(nil).SendMessage), text, userID, userName)
}

// MockPurchasesModel is a mock of PurchasesModel interface.
type MockPurchasesModel struct {
	ctrl     *gomock.Controller
	recorder *MockPurchasesModelMockRecorder
}

// MockPurchasesModelMockRecorder is the mock recorder for MockPurchasesModel.
type MockPurchasesModelMockRecorder struct {
	mock *MockPurchasesModel
}

// NewMockPurchasesModel creates a new mock instance.
func NewMockPurchasesModel(ctrl *gomock.Controller) *MockPurchasesModel {
	mock := &MockPurchasesModel{ctrl: ctrl}
	mock.recorder = &MockPurchasesModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPurchasesModel) EXPECT() *MockPurchasesModelMockRecorder {
	return m.recorder
}

// AddCategory mocks base method.
func (m *MockPurchasesModel) AddCategory(userID int64, category string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCategory", userID, category)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCategory indicates an expected call of AddCategory.
func (mr *MockPurchasesModelMockRecorder) AddCategory(userID, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCategory", reflect.TypeOf((*MockPurchasesModel)(nil).AddCategory), userID, category)
}

// AddPurchase mocks base method.
func (m *MockPurchasesModel) AddPurchase(userID int64, rawSum, category, rawDate string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPurchase", userID, rawSum, category, rawDate)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPurchase indicates an expected call of AddPurchase.
func (mr *MockPurchasesModelMockRecorder) AddPurchase(userID, rawSum, category, rawDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPurchase", reflect.TypeOf((*MockPurchasesModel)(nil).AddPurchase), userID, rawSum, category, rawDate)
}

// Report mocks base method.
func (m *MockPurchasesModel) Report(period purchases.Period, userID int64) (string, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Report", period, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Report indicates an expected call of Report.
func (mr *MockPurchasesModelMockRecorder) Report(period, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Report", reflect.TypeOf((*MockPurchasesModel)(nil).Report), period, userID)
}

// ToPeriod mocks base method.
func (m *MockPurchasesModel) ToPeriod(str string) (purchases.Period, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToPeriod", str)
	ret0, _ := ret[0].(purchases.Period)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToPeriod indicates an expected call of ToPeriod.
func (mr *MockPurchasesModelMockRecorder) ToPeriod(str interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToPeriod", reflect.TypeOf((*MockPurchasesModel)(nil).ToPeriod), str)
}
