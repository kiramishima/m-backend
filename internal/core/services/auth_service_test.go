package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"kiramishima/m-backend/internal/core/domain"
	mock "kiramishima/m-backend/internal/mocks"
	"testing"
	"time"
)

func setup() {

}

func TestLogin(t *testing.T) {
	logger, _ := zap.NewProduction()
	slogger := logger.Sugar()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mock.NewMockAuthRepository(mockCtrl)
	repo.EXPECT().FindByCredentials(gomock.Any(), gomock.Any()).Return(&domain.User{
		ID:        "1",
		Email:     "gini@mail.com",
		Password:  "12356",
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}, nil)

	uc := NewAuthService(slogger, repo, 2)

	t.Run("OK", func(t *testing.T) {
		ctx := context.Background()
		data := &domain.AuthRequest{Email: "gini@mail.com", Password: "123456"}
		b, err := uc.FindByCredentials(ctx, data)
		t.Log("B", b)
		assert.NoError(t, err)
		assert.Equal(t, "gini@mail.com", b.Token)
	})
	t.Run("Not Found", func(t *testing.T) {
		ctx := context.Background()
		data := &domain.AuthRequest{Email: "gini@mail.com", Password: ""}
		b, err := uc.FindByCredentials(ctx, data)
		t.Log(b)
		t.Log(err)
		assert.Error(t, err)
	})
}
