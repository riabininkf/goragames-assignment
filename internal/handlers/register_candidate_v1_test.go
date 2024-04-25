package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/riabininkf/goragames-assignment/internal/handlers"
	"github.com/riabininkf/goragames-assignment/internal/http"
	"github.com/riabininkf/goragames-assignment/internal/repository"
	"github.com/riabininkf/goragames-assignment/internal/repository/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"io"
	netHttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterCandidateV1_Method(t *testing.T) {
	handler := handlers.NewRegisterCandidateV1(nil, nil, nil, nil)
	assert.Equal(t, http.MethodPost, handler.Method())
}

func TestRegisterCandidateV1_Path(t *testing.T) {
	handler := handlers.NewRegisterCandidateV1(nil, nil, nil, nil)
	assert.Equal(t, "/api/v1/invitations/:code", handler.Path())
}

func TestRegisterCandidateV1_Handle(t *testing.T) {
	testCases := []struct {
		name                         string
		params                       []gin.Param
		body                         map[string]interface{}
		onGetCandidateByEmailAndCode func() (*entity.Candidate, error)
		onGetCodeByNameWithLock      func() (*entity.Code, error)
		onDecrementCodeUsagesByName  func() error
		onAddCandidate               func() error
		onDo                         func() error
		expResponse                  http.Response
	}{
		{
			name:        "code is missing",
			expResponse: http.NewBadRequest("code is required"),
		},
		{
			name:        "can't unmarshal request body to json",
			params:      []gin.Param{{Key: "code", Value: "some_code"}},
			body:        map[string]interface{}{"email": []int{}},
			expResponse: http.InternalServerError,
		},
		{
			name:        "email is empty",
			params:      []gin.Param{{Key: "code", Value: "some_code"}},
			body:        map[string]interface{}{},
			expResponse: http.NewBadRequest("email is required"),
		},
		{
			name:                         "can't get candidate by email and code",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.InternalServerError,
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, errors.New("test error") },
		},
		{
			name:                         "candidate is already exist",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.NewBadRequest("candidate is already exist"),
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return &entity.Candidate{}, nil },
		},
		{
			name:                         "can't get code by name",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.InternalServerError,
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, nil },
			onGetCodeByNameWithLock:      func() (*entity.Code, error) { return nil, errors.New("test error") },
			onDo:                         func() error { return errors.New("test error") },
		},
		{
			name:                         "can't get code by name",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.InternalServerError,
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, nil },
			onGetCodeByNameWithLock:      func() (*entity.Code, error) { return nil, errors.New("test error") },
			onDo:                         func() error { return errors.New("test error") },
		},
		{
			name:                         "code not found",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.NewBadNotFound("code not found"),
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, nil },
			onGetCodeByNameWithLock:      func() (*entity.Code, error) { return nil, nil },
			onDo:                         func() error { return handlers.ErrCodeNotFound },
		},
		{
			name:                         "all usages are spent",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.NewBadRequest("registration for this code is closed"),
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, nil },
			onGetCodeByNameWithLock:      func() (*entity.Code, error) { return &entity.Code{}, nil },
			onDo:                         func() error { return handlers.ErrUsagesSpent },
		},
		{
			name:                         "can't decrement usages by name",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.InternalServerError,
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, nil },
			onGetCodeByNameWithLock: func() (*entity.Code, error) {
				return &entity.Code{Usages: 1, Name: "some_code"}, nil
			},
			onDecrementCodeUsagesByName: func() error { return errors.New("test error") },
			onDo:                        func() error { return errors.New("test error") },
		},
		{
			name:                         "can't add candidate",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.InternalServerError,
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, nil },
			onGetCodeByNameWithLock: func() (*entity.Code, error) {
				return &entity.Code{Usages: 1, Name: "some_code"}, nil
			},
			onDecrementCodeUsagesByName: func() error { return nil },
			onAddCandidate:              func() error { return errors.New("test error") },
			onDo:                        func() error { return errors.New("test error") },
		},
		{
			name:                         "candidate already exists (concurrency)",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.NewBadRequest("candidate is already exist"),
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, nil },
			onGetCodeByNameWithLock: func() (*entity.Code, error) {
				return &entity.Code{Usages: 1, Name: "some_code"}, nil
			},
			onDecrementCodeUsagesByName: func() error { return nil },
			onAddCandidate:              func() error { return repository.ErrDuplicatedRow },
			onDo:                        func() error { return handlers.ErrCandidateAlreadyExists },
		},
		{
			name:                         "positive case",
			params:                       []gin.Param{{Key: "code", Value: "some_code"}},
			body:                         map[string]interface{}{"email": "someEmail"},
			expResponse:                  http.NewEmptyResponse(),
			onGetCandidateByEmailAndCode: func() (*entity.Candidate, error) { return nil, nil },
			onGetCodeByNameWithLock: func() (*entity.Code, error) {
				return &entity.Code{Usages: 1, Name: "some_code"}, nil
			},
			onDecrementCodeUsagesByName: func() error { return nil },
			onAddCandidate:              func() error { return nil },
			onDo:                        func() error { return nil },
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			body, err := json.Marshal(testCase.body)
			assert.NoError(t, err)

			c.Request = &netHttp.Request{Body: io.NopCloser(bytes.NewBuffer(body))}
			c.Params = testCase.params

			log := zap.NewNop()

			candidates := repository.NewCandidatesMock(t)

			if testCase.onGetCandidateByEmailAndCode != nil {
				candidates.On("GetByEmailAndCode", ctx, testCase.body["email"], c.Param("code")).
					Return(testCase.onGetCandidateByEmailAndCode())
			}

			sessCtx := mongo.NewSessionContext(ctx, nil)

			tx := repository.NewTxMock(t)

			if testCase.onDo != nil {
				tx.On("Do", ctx, mock.Anything).Run(func(args mock.Arguments) {
					_, _ = args[1].(func(ctx mongo.SessionContext) (interface{}, error))(sessCtx)
				}).Return(testCase.onDo())
			}

			codes := repository.NewCodesMock(t)
			if testCase.onGetCodeByNameWithLock != nil {
				result, err2 := testCase.onGetCodeByNameWithLock()
				codes.On("GetByNameWithLock", sessCtx, c.Param("code")).
					Return(result, err2)

			}

			if testCase.onDecrementCodeUsagesByName != nil {
				codes.On("DecrementUsagesByName", sessCtx, c.Param("code")).
					Return(testCase.onDecrementCodeUsagesByName())
			}

			if testCase.onAddCandidate != nil {
				candidates.On("Add", sessCtx, &entity.Candidate{
					Email: testCase.body["email"].(string),
					Code:  c.Param("code"),
				}).Return(testCase.onAddCandidate())
			}

			//if testCase.onDo != nil {
			//	tx.On("Do", ctx, mock.Anything).Run(func(args mock.Arguments) {
			//		_, _ = args[1].(func(ctx mongo.SessionContext) (interface{}, error))(sessCtx)
			//	}).Return(txErr)
			//}

			handler := handlers.NewRegisterCandidateV1(tx, log, codes, candidates)

			assert.Equal(t, testCase.expResponse, handler.Handle(ctx, c))

		})
	}
}
