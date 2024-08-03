package service

import (
	"context"
	"go_schedule_service/config"
	"go_schedule_service/genproto/lesson_service"

	"go_schedule_service/grpc/client"
	"go_schedule_service/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LessonService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.GrpcClientI
}

func NewLessonService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.GrpcClientI) *LessonService {
	return &LessonService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (f *LessonService) Create(ctx context.Context, req *lesson_service.CreateLesson) (*lesson_service.GetLesson, error) {

	f.log.Info("---CreateLesson--->>>", logger.Any("req", req))

	resp, err := f.strg.Lesson().Create(ctx, req)
	if err != nil {
		f.log.Error("---CreateLesson--->>>", logger.Error(err))
		return &lesson_service.GetLesson{}, err
	}

	return resp, nil
}
func (f *LessonService) Update(ctx context.Context, req *lesson_service.UpdateLesson) (*lesson_service.GetLesson, error) {

	f.log.Info("---UpdateLesson--->>>", logger.Any("req", req))

	resp, err := f.strg.Lesson().Update(ctx, req)
	if err != nil {
		f.log.Error("---UpdateLesson--->>>", logger.Error(err))
		return &lesson_service.GetLesson{}, err
	}

	return resp, nil
}

func (f *LessonService) GetList(ctx context.Context, req *lesson_service.GetListLessonRequest) (*lesson_service.GetListLessonResponse, error) {
	f.log.Info("---GetListLesson--->>>", logger.Any("req", req))

	resp, err := f.strg.Lesson().GetAll(ctx, req)
	if err != nil {
		f.log.Error("---GetListLesson--->>>", logger.Error(err))
		return &lesson_service.GetListLessonResponse{}, err
	}

	return resp, nil
}

func (f *LessonService) GetByID(ctx context.Context, id *lesson_service.LessonPrimaryKey) (*lesson_service.GetLesson, error) {
	f.log.Info("---GetLesson--->>>", logger.Any("req", id))

	resp, err := f.strg.Lesson().GetById(ctx, id)
	if err != nil {
		f.log.Error("---GetLesson--->>>", logger.Error(err))
		return &lesson_service.GetLesson{}, err
	}

	return resp, nil
}

func (f *LessonService) Delete(ctx context.Context, req *lesson_service.LessonPrimaryKey) (*emptypb.Empty, error) {

	f.log.Info("---DeleteLesson--->>>", logger.Any("req", req))

	_, err := f.strg.Lesson().Delete(ctx, req)
	if err != nil {
		f.log.Error("---DeleteLesson--->>>", logger.Error(err))
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
