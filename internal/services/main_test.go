package services_test

import (
	"os"
	"testing"

	mockSheet "github.com/SmmTouch-com/instagram-notification-service/pkg/googlesheets/mocks"
	mockMailer "github.com/SmmTouch-com/instagram-notification-service/pkg/mail/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ServicesTestSuite struct {
	suite.Suite

	mocks *mocks
}

type mocks struct {
	sheet *mockSheet.MockSheet
	mailer *mockMailer.MockMailer
}

func TestServicesSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(ServicesTestSuite))
}

func (s *ServicesTestSuite) SetupSuite() {
	s.initMocks()
}

func (s *ServicesTestSuite) initMocks() {
	mockCtl := gomock.NewController(s.T())
	defer mockCtl.Finish()

	s.mocks = &mocks{
		sheet: mockSheet.NewMockSheet(mockCtl),
		mailer: mockMailer.NewMockMailer(mockCtl),
	}
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
