package service

import (
	"context"
	"go_schedule_service/config"
	"go_schedule_service/genproto/perfomance_service"
	"go_schedule_service/genproto/student_service"

	"go_schedule_service/grpc/client"
	"go_schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PerfomanceService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.GrpcClientI
}

func NewPerfomanceService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.GrpcClientI) *PerfomanceService {
	return &PerfomanceService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *PerfomanceService) Create(ctx context.Context, req *perfomance_service.CreatePerfomance) (*perfomance_service.GetPerfomance, error) {
	f.log.Info("---CreatePerfomance--->>>", logger.Any("req", req))

	_, err := f.services.Student().Check(ctx, &student_service.StudentPrimaryKey{Id: req.StudentId})
	if err != nil {
		f.log.Error("error while check student")
		return nil, err
	}
	resp, err := f.strg.Perfomance().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreatePerfomance--->>>", logger.Error(err))
		return &perfomance_service.GetPerfomance{}, err
	}

	return resp, nil
}
func (f *PerfomanceService) Update(ctx context.Context, req *perfomance_service.UpdatePerfomance) (*perfomance_service.GetPerfomance, error) {

	f.log.Info("---UpdatePerfomance--->>>", logger.Any("req", req))

	resp, err := f.strg.Perfomance().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdatePerfomance--->>>", logger.Error(err))
		return &perfomance_service.GetPerfomance{}, err
	}

	return resp, nil
}

func (f *PerfomanceService) GetList(ctx context.Context, req *perfomance_service.GetListPerfomanceRequest) (*perfomance_service.GetListPerfomanceResponse, error) {
	f.log.Info("---GetListPerfomance--->>>", logger.Any("req", req))

	resp, err := f.strg.Perfomance().GetAll(ctx, req)
	if err != nil {
		f.log.Error("---GetListPerfomance--->>>", logger.Error(err))
		return &perfomance_service.GetListPerfomanceResponse{}, err
	}

	return resp, nil
}

func (f *PerfomanceService) GetByID(ctx context.Context, id *perfomance_service.PerfomancePrimaryKey) (*perfomance_service.GetPerfomance, error) {
	f.log.Info("---GetPerfomance--->>>", logger.Any("req", id))

	resp, err := f.strg.Perfomance().GetById(ctx, id)
	if err != nil {
		f.log.Error("---GetPerfomance--->>>", logger.Error(err))
		return &perfomance_service.GetPerfomance{}, err
	}

	return resp, nil
}

func (f *PerfomanceService) Delete(ctx context.Context, req *perfomance_service.PerfomancePrimaryKey) (*emptypb.Empty, error) {

	f.log.Info("---DeletePerfomance--->>>", logger.Any("req", req))

	_, err := f.strg.Perfomance().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeletePerfomance--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
