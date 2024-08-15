package service

import (
	pb "Booking_Service/genproto/booking"
	"Booking_Service/models"
	"Booking_Service/pkg/logger"
	"Booking_Service/storage"
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
)

type NotificationService struct {
	pb.UnimplementedNotificationsServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewNotificationService(s storage.IStorage) *NotificationService {
	return &NotificationService{
		storage: s,
		logger:  logger.NewLogger(),
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, req *pb.NewNotification) (*pb.ID, error) {
	s.logger.Info("CreateNotification is invoked", slog.Any("request", req))

	time := time.Now().Format(time.RFC3339)
	nf := models.NewNotification{
		UserID:    req.UserId,
		Title:     req.Title,
		Message:   req.Message,
		CreatedAt: time,
	}

	id, err := s.storage.Notification().Create(ctx, &nf)
	if err != nil {
		er := errors.Wrap(err, "failed to create notification")
		s.logger.Error(er.Error())
		return nil, er
	}

	resp := &pb.ID{Id: id}
	s.logger.Info("CreateNotification is completed", slog.Any("response", resp))
	return resp, nil
}

func (s *NotificationService) GetNotification(ctx context.Context, req *pb.ID) (*pb.Notification, error) {
	s.logger.Info("GetNotification is invoked", slog.Any("request", req))

	nf, err := s.storage.Notification().Get(ctx, req.Id)
	if err != nil {
		er := errors.Wrap(err, "failed to get notification")
		s.logger.Error(er.Error())
		return nil, er
	}

	resp := &pb.Notification{
		Id:        nf.Id,
		UserId:    nf.UserID,
		Title:     nf.Title,
		Message:   nf.Message,
		CreatedAt: nf.CreatedAt,
	}
	s.logger.Info("GetNotification is completed", slog.Any("response", resp))
	return resp, nil
}
