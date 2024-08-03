package service

import (
	"context"
	"go_schedule_service/config"
	"go_schedule_service/genproto/group_service"
	"go_schedule_service/genproto/schedule_service"

	"go_schedule_service/grpc/client"
	"go_schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ScheduleService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.GrpcClientI
}

func NewScheduleService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.GrpcClientI) *ScheduleService {
	return &ScheduleService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *ScheduleService) Create(ctx context.Context, req *schedule_service.CreateSchedule) (*schedule_service.GetSchedule, error) {

	f.log.Info("---CreateSchedule--->>>", logger.Any("req", req))
	_, err := f.services.Group().Check(ctx, &group_service.GroupPrimaryKey{Id: req.GroupId})
	if err != nil {
		f.log.Error("error while check group")
		return nil, err
	}
	resp, err := f.strg.Schedule().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreateSchedule--->>>", logger.Error(err))
		return &schedule_service.GetSchedule{}, err
	}

	return resp, nil
}
func (f *ScheduleService) Update(ctx context.Context, req *schedule_service.UpdateSchedule) (*schedule_service.GetSchedule, error) {

	f.log.Info("---UpdateSchedule--->>>", logger.Any("req", req))

	resp, err := f.strg.Schedule().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdateSchedule--->>>", logger.Error(err))
		return &schedule_service.GetSchedule{}, err
	}

	return resp, nil
}

func (f *ScheduleService) GetList(ctx context.Context, req *schedule_service.GetListScheduleRequest) (*schedule_service.GetListScheduleResponse, error) {
	f.log.Info("---GetListSchedule--->>>", logger.Any("req", req))

	resp, err := f.strg.Schedule().GetAll(ctx, req)
	if err != nil {
		f.log.Error("---GetListSchedule--->>>", logger.Error(err))
		return &schedule_service.GetListScheduleResponse{}, err
	}

	return resp, nil
}

func (f *ScheduleService) GetByID(ctx context.Context, id *schedule_service.SchedulePrimaryKey) (*schedule_service.GetSchedule, error) {
	f.log.Info("---GetSchedule--->>>", logger.Any("req", id))

	resp, err := f.strg.Schedule().GetById(ctx, id)
	if err != nil {
		f.log.Error("---GetSchedule--->>>", logger.Error(err))
		return &schedule_service.GetSchedule{}, err
	}

	return resp, nil
}

func (f *ScheduleService) Delete(ctx context.Context, req *schedule_service.SchedulePrimaryKey) (*emptypb.Empty, error) {

	f.log.Info("---DeleteSchedule--->>>", logger.Any("req", req))

	_, err := f.strg.Schedule().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeleteSchedule--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (f *ScheduleService) GetStudentSchedule(ctx context.Context, id *schedule_service.SchedulePrimaryKey) (*schedule_service.GetStudentSchedules, error) {
	f.log.Info("---GetStudentSchedule--->>>", logger.Any("req", id))

	resp, err := f.strg.Schedule().GetStudentSchedule(ctx, id)
	if err != nil {
		f.log.Error("---GetStudentSchedule--->>>", logger.Error(err))
		return &schedule_service.GetStudentSchedules{}, err
	}

	return resp, nil
}

func (f *ScheduleService) GetForWeek(ctx context.Context, req *schedule_service.GetWeekScheduleRequest) (*schedule_service.GetWeekScheduleResponse, error) {
	f.log.Info("---GetWeekSchedule--->>>", logger.Any("req", req))

	resp, err := f.strg.Schedule().GetForWeek(ctx, req)
	if err != nil {
		f.log.Error("---GetWeekSchedule--->>>", logger.Error(err))
		return nil, err
	}

	return resp, nil
}
