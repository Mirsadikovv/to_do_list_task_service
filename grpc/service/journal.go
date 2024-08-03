package service

import (
	"context"
	"go_schedule_service/config"
	"go_schedule_service/genproto/journal_service"

	"go_schedule_service/grpc/client"
	"go_schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type JournalService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.GrpcClientI
}

func NewJournalService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.GrpcClientI) *JournalService {
	return &JournalService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *JournalService) Create(ctx context.Context, req *journal_service.CreateJournal) (*journal_service.GetJournal, error) {

	f.log.Info("---CreateJournal--->>>", logger.Any("req", req))

	resp, err := f.strg.Journal().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreateJournal--->>>", logger.Error(err))
		return &journal_service.GetJournal{}, err
	}

	return resp, nil
}
func (f *JournalService) Update(ctx context.Context, req *journal_service.UpdateJournal) (*journal_service.GetJournal, error) {

	f.log.Info("---UpdateJournal--->>>", logger.Any("req", req))

	resp, err := f.strg.Journal().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdateJournal--->>>", logger.Error(err))
		return &journal_service.GetJournal{}, err
	}

	return resp, nil
}

func (f *JournalService) GetList(ctx context.Context, req *journal_service.GetListJournalRequest) (*journal_service.GetListJournalResponse, error) {
	f.log.Info("---GetListJournal--->>>", logger.Any("req", req))

	resp, err := f.strg.Journal().GetAll(ctx, req)
	if err != nil {
		f.log.Error("---GetListJournal--->>>", logger.Error(err))
		return &journal_service.GetListJournalResponse{}, err
	}

	return resp, nil
}

func (f *JournalService) GetByID(ctx context.Context, id *journal_service.JournalPrimaryKey) (*journal_service.GetJournal, error) {
	f.log.Info("---GetJournal--->>>", logger.Any("req", id))

	resp, err := f.strg.Journal().GetById(ctx, id)
	if err != nil {
		f.log.Error("---GetJournal--->>>", logger.Error(err))
		return &journal_service.GetJournal{}, err
	}

	return resp, nil
}

func (f *JournalService) Delete(ctx context.Context, req *journal_service.JournalPrimaryKey) (*emptypb.Empty, error) {

	f.log.Info("---DeleteJournal--->>>", logger.Any("req", req))

	_, err := f.strg.Journal().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeleteJournal--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
