package handlers

import (
	"github.com/riabininkf/goragames-assignment/internal/container"
	"github.com/riabininkf/goragames-assignment/internal/http"
	"github.com/riabininkf/goragames-assignment/internal/logger"
	"github.com/riabininkf/goragames-assignment/internal/repository"
	"github.com/sarulabs/di/v2"
)

const DefRegisterCandidateV1Name = "handler.register_candidate_v1"

func init() {
	container.Add(di.Def{
		Name: DefRegisterCandidateV1Name,
		Tags: []di.Tag{{Name: http.TagHandler}},
		Build: func(ctn di.Container) (interface{}, error) {
			var tx repository.Tx
			if err := container.Fill(repository.DefTxName, &tx); err != nil {
				return nil, err
			}

			var candidates repository.Candidates
			if err := container.Fill(repository.DefCandidatesName, &candidates); err != nil {
				return nil, err
			}

			var codes repository.Codes
			if err := container.Fill(repository.DefCodesName, &codes); err != nil {
				return nil, err
			}

			var log logger.Logger
			if err := container.Fill(logger.DefName, &log); err != nil {
				return nil, err
			}

			return NewRegisterCandidateV1(tx, log, codes, candidates), nil
		},
	})
}
