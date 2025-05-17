package service

import (
	"context"
	"errors"
	"github.com/germanov-v/go-rotator/internal/model"
	"math"
	"testing"
)

// go test ./internal/service

type mockRepo struct {
	banners []model.Banner
	stats   map[string]*model.Stats
	//displays []model.StatTemp
	displays []struct {
		slot   model.SlotId
		group  model.GroupId
		banner model.BannerId
	}
}

func (r *mockRepo) AddBanner(ctx context.Context, slot model.SlotId, banner model.BannerId) error {
	return nil
}

func (r *mockRepo) AddGroup(ctx context.Context, group model.GroupId) error {
	return nil
}

func (r *mockRepo) RemoveBanner(ctx context.Context, slot model.SlotId, banner model.BannerId) error {
	return nil
}

func (m *mockRepo) IncrementClick(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) error {
	return nil
}

func (m *mockRepo) ListBanners(ctx context.Context, slot model.SlotId) ([]model.Banner, error) {
	return m.banners, nil
}

func (m *mockRepo) GetStats(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) (*model.Stats, error) {
	key := string(slot) + ":" + string(banner) + ":" + string(group)
	if st, ok := m.stats[key]; ok {
		return st, nil
	}
	return nil, errors.New("not found data")
}
func (m *mockRepo) IncrementDisplay(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) error {

	//key := string(slot) + ":" + string(banner) + ":" + string(group)
	//m.displays = append(m.displays, struct {
	//	slot   model.SlotId
	//	banner model.BannerId
	//	group  model.GroupId
	//}{slot, banner, group})
	//if st, ok := m.stats[key]; ok {
	//	st.CountDisplay++
	//} else {
	//	m.stats[key] = &model.Stats{Slot: slot, Group: group, CountDisplay: 1, Clicks: 0}
	//}
	//return nil
	item := struct {
		slot   model.SlotId
		group  model.GroupId
		banner model.BannerId
	}{slot: slot, group: group, banner: banner}
	m.displays = append(m.displays, item)
	key := string(slot) + ":" + string(banner) + ":" + string(group)
	if st, ok := m.stats[key]; ok {
		st.CountDisplay++
	} else {
		m.stats[key] = &model.Stats{Slot: slot, Banner: banner, Group: group, CountDisplay: 1, Clicks: 0}
	}
	return nil
}

func TestRotate_NoBanners(t *testing.T) {
	repo := &mockRepo{banners: nil, stats: make(map[string]*model.Stats)}
	svc := NewRotationService(repo)
	_, err := svc.Rotate(context.Background(), "slot1", "group1")
	if err == nil {
		t.Fatal("expected error =>  no baners found")
	}
}

func TestRotate_FirstDisplay(t *testing.T) {
	repo := &mockRepo{
		banners: []model.Banner{{Id: "a"}, {Id: "b"}},
		stats:   make(map[string]*model.Stats),
	}
	svc := NewRotationService(repo)

	// a
	banner, err := svc.Rotate(context.Background(), "slot1", "group1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// aa
	if banner != "a" {
		t.Errorf("expected 'a', got %s", banner)
	}

	// b проверяем
	banner2, _ := svc.Rotate(context.Background(), "slot1", "group1")
	if banner2 != "b" {
		t.Errorf("expected 'b', got %s", banner2)
	}
}

func TestRotate_UCB1Balance(t *testing.T) {
	repo := &mockRepo{
		//banners: []model.Banner{{Id: "a"}, {Id: "bb"}},
		banners: []model.Banner{{Id: "a"}, {Id: "b"}},
		stats:   make(map[string]*model.Stats),
	}
	repo.stats["slot1:a:group1"] = &model.Stats{Slot: "slot1", Group: "group1", CountDisplay: 10, Clicks: 2}
	repo.stats["slot1:b:group1"] = &model.Stats{Slot: "slot1", Group: "group1", CountDisplay: 10, Clicks: 5}

	svc := NewRotationService(repo)
	banner, err := svc.Rotate(context.Background(), "slot1", "group1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := float64(10 + 10)
	lnTotal := math.Log(total)
	sqrt := math.Sqrt(2 * lnTotal / 10)
	scoreA := 0.2 + sqrt
	scoreB := 0.5 + sqrt
	var expected model.BannerId
	if scoreA > scoreB {
		expected = model.BannerId("a")
	} else {
		expected = model.BannerId("b")
	}

	if banner != expected {
		t.Errorf("expected %s, got %s", expected, banner)
	}
}
