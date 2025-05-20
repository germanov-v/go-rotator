package repository

import (
	"context"
	"github.com/germanov-v/go-rotator/internal/model"
)

type Repository interface {
	AddBanner(ctx context.Context, slot model.SlotId, banner model.BannerId) error
	AddGroup(ctx context.Context, group model.GroupId) error
	AddSlot(ctx context.Context, slot model.SlotId) error
	RemoveBanner(ctx context.Context, slot model.SlotId, banner model.BannerId) error
	IncrementDisplay(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) error
	IncrementClick(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) error
	ListBanners(ctx context.Context, slot model.SlotId) ([]model.Banner, error)
	GetStats(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) (*model.Stats, error)
}
