package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/riabininkf/goragames-assignment/internal/http"
	"github.com/riabininkf/goragames-assignment/internal/logger"
	"github.com/riabininkf/goragames-assignment/internal/repository"
	"github.com/riabininkf/goragames-assignment/internal/repository/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
)

var (
	ErrCodeNotFound           = errors.New("code not found")
	ErrUsagesSpent            = errors.New("all usages are spent")
	ErrCandidateAlreadyExists = errors.New("candidate already exists")
)

func NewRegisterCandidateV1(
	tx repository.Tx,
	log logger.Logger,
	codes repository.Codes,
	candidates repository.Candidates,
) http.Handler {
	return &registerCandidateV1{
		tx:         tx,
		log:        log,
		codes:      codes,
		candidates: candidates,
	}
}

type (
	registerCandidateV1 struct {
		tx         repository.Tx
		log        logger.Logger
		codes      repository.Codes
		candidates repository.Candidates
	}

	registerCandidateV1Request struct {
		Email string `json:"email"`
	}
)

func (h *registerCandidateV1) Method() string {
	return http.MethodPost
}

func (h *registerCandidateV1) Path() string {
	return "/api/v1/invitations/:code"
}

func (h *registerCandidateV1) Handle(ctx context.Context, req *gin.Context) http.Response {
	// todo since there are only 3 possible codes, we can save them in memory

	var codeName string
	if codeName = req.Param("code"); codeName == "" {
		h.log.Warn("code is empty")
		return http.NewBadRequest("code is required")
	}

	var reqBody registerCandidateV1Request
	if err := req.BindJSON(&reqBody); err != nil && !errors.Is(err, io.EOF) {
		h.log.Error("can't unmarshal request body to json", logger.Error(err))
		fmt.Println(err.Error())
		return http.InternalServerError
	}

	if reqBody.Email == "" {
		h.log.Warn("email is empty")
		return http.NewBadRequest("email is required")
	}

	// todo validate email

	//todo add rate limit by email and/or IP address

	var (
		err       error
		candidate *entity.Candidate
	)
	if candidate, err = h.candidates.GetByEmailAndCode(ctx, reqBody.Email, codeName); err != nil {
		h.log.Error("can't get candidate by email and code", logger.Error(err))
		return http.InternalServerError
	}

	if candidate.IsExist() {
		h.log.Warn("candidate with provided email and code is already exist")
		return http.NewBadRequest("candidate is already exist")
	}

	// todo add translation for errors

	candidate = &entity.Candidate{Email: reqBody.Email, Code: codeName}

	//todo find a way to pass response error from callback
	if err = h.tx.Do(ctx, h.register(candidate)); err != nil {
		switch {
		case errors.Is(err, ErrCodeNotFound):
			h.log.Warn("code not found")
			return http.NewBadNotFound("code not found")
		case errors.Is(err, ErrUsagesSpent):
			h.log.Warn("all usages are spent")
			return http.NewBadRequest("registration for this code is closed")
		case errors.Is(err, ErrCandidateAlreadyExists):
			h.log.Warn("candidate with provided email and code is already exist")
			return http.NewBadRequest("candidate is already exist")
		default:
			h.log.Error("can't register code", logger.Error(err))
			return http.InternalServerError
		}
	}

	// todo trigger event
	return http.NewEmptyResponse()
}

func (h *registerCandidateV1) register(candidate *entity.Candidate) func(ctx mongo.SessionContext) (interface{}, error) {
	return func(ctx mongo.SessionContext) (interface{}, error) {
		var (
			err  error
			code *entity.Code
		)
		if code, err = h.codes.GetByNameWithLock(ctx, candidate.Code); err != nil {
			return nil, fmt.Errorf("can't get code by name: %w", err)
		}

		if !code.IsExist() {
			return nil, ErrCodeNotFound
		}

		if code.Usages == 0 {
			return nil, ErrUsagesSpent
		}

		if err = h.codes.DecrementUsagesByName(ctx, code.Name); err != nil {
			return nil, fmt.Errorf("can't decrement usages by name: %w", err)
		}

		if err = h.candidates.Add(ctx, candidate); err != nil {
			if errors.Is(err, repository.ErrDuplicatedRow) {
				return nil, ErrCandidateAlreadyExists
			}

			return nil, fmt.Errorf("can't add candidate: %w", err)
		}

		return nil, nil
	}
}
