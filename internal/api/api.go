package api

import (
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/api/health"
	queue "github.com/System-Analysis-and-Design-2023-SUT/Server/internal/api/queue"
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/api/server"
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/helper"
	queuerepo "github.com/System-Analysis-and-Design-2023-SUT/Server/internal/repository/queue"
	queueservice "github.com/System-Analysis-and-Design-2023-SUT/Server/internal/services/queue"
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/settings"
	models "github.com/System-Analysis-and-Design-2023-SUT/Server/models/queue"
	"github.com/pkg/errors"
)

func NewAPIServer(settings *settings.Settings, helper *helper.Helper, q *models.Queue, s *models.Subscriber) (*server.Server, error) {
	queueRepo, err := queuerepo.NewRepository(settings, helper, q, s)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize user repository")
	}

	queueService, err := queueservice.NewService(queueRepo)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize user service")
	}

	queueModule, err := queue.NewQueueModule(queueRepo, queueService, settings, helper)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize users module")
	}

	healthModule, err := health.NewHealth(settings.Global.Environment)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize health module")
	}

	srv, err := server.NewServer(queueModule, healthModule, settings)
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize api server object")
	}

	return srv, nil
}
