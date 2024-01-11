package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/unrolled/render"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"kiramishima/m-backend/internal/core/domain"
	mock "kiramishima/m-backend/internal/mocks"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignInHandler(t *testing.T) {
	testCases := map[string]struct {
		ID            any
		buildStubs    func(uc *mock.MockAuthService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		"OK": {
			ID: 1,
			buildStubs: func(uc *mock.MockAuthService) {
				uc.EXPECT().
					FindByCredentials(gomock.Any(), &domain.AuthRequest{Email: "gini@mail.com", Password: "123456"}).
					Times(1).
					Return(&domain.User{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusAccepted, recorder.Code)
			},
		},
		/*"Invalid URL Param": {
			ID: "ID",
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		"Not Found": {
			ID: 0,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, repoErr.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		"Unexpected Error": {
			ID: 1,
			buildStubs: func(uc *mock.MockUserUsecase) {
				uc.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},*/
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockAuthService(ctrl)
			tc.buildStubs(uc)

			recorder := httptest.NewRecorder()

			url := "/v1/auth/sign-in"
			teacher := domain.AuthRequest{
				Email:    "gini@mail.com",
				Password: "123456",
			}
			// marshall data to json (like json_encode)
			marshalled, err := json.Marshal(teacher)
			if err != nil {
				log.Fatalf("impossible to marshall form: %s", err)
			}

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshalled))
			assert.NoError(t, err)

			router := chi.NewRouter()
			logger, _ := zap.NewProduction()
			slogger := logger.Sugar()
			r := render.New()
			NewAuthHandlers(router, slogger, uc, r)
			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
