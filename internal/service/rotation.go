package service

import (
	"context"
	"github.com/germanov-v/go-rotator/internal/model"
	"github.com/germanov-v/go-rotator/internal/repository"
	"github.com/pkg/errors"
	"math"
)

type RotationService struct {
	repo repository.Repository
}

func NewRotationService(repo repository.Repository) *RotationService {
	return &RotationService{repo: repo}
}

func (s *RotationService) Rotate(ctx context.Context, slot model.SlotId, group model.GroupId) (model.BannerId, error) {
	banners, err := s.repo.ListBanners(ctx, slot)

	if err != nil {
		return "", err
	}

	if len(banners) == 0 {
		return "", errors.New("no banners found")
	}

	var stats []model.StatTemp
	var totalDisplay int64

	for _, banner := range banners {
		st, err := s.repo.GetStats(ctx, slot, banner.Id, group)
		if err != nil {
			st = &model.Stats{Slot: slot, Group: group, CountDisplay: 0, Clicks: 0}
		}
		stats = append(stats, model.StatTemp{Id: banner.Id, CountDisplays: st.CountDisplay, CountClicks: st.Clicks})
		totalDisplay += st.CountDisplay
	}

	// план минимум: показываем, которые не показывались
	//for key, st := range stats {
	for _, st := range stats {
		if st.CountDisplays == 0 {
			_ = s.repo.IncrementDisplay(ctx, slot, st.Id, group)
			return st.Id, nil
		}
	}

	// ucb1
	var best model.BannerId
	var bestScore float64 = -1
	// логариф
	//totalLn := math.Log(totalDisplay)
	totalLn := math.Log(float64(totalDisplay))

	for _, st := range stats {
		// ctr Click-Through Rate
		ctr := float64(st.CountClicks) / float64(st.CountDisplays)

		// коэфффциент с использованием интегрвала доверительного
		sqrt := math.Sqrt(2 * totalLn / float64(st.CountDisplays))

		// балл
		score := ctr + bestScore*sqrt

		if score > bestScore {
			bestScore = score
			best = st.Id
		}
	}

	if err := s.repo.IncrementDisplay(ctx, slot, best, group); err != nil {
		// return "", err
		//	return best, err
		return best, err
	}

	return best, nil
}
