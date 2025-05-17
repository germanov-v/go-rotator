package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/germanov-v/go-rotator/internal/model"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(connectionStr string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		return nil, err
	}

	return &PostgresRepo{db}, nil
}

func (r *PostgresRepo) AddBanner(ctx context.Context, slot model.SlotId, banner model.BannerId) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		"INSERT INTO slots(id) VALUES($1) ON CONFLICT(id) DO NOTHING",
		slot,
	); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		"INSERT INTO banners(id) VALUES($1) ON CONFLICT(id) DO NOTHING",
		banner,
	); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		"INSERT INTO slot_banners(slot_id,banner_id) VALUES($1,$2) ON CONFLICT(slot_id,banner_id) DO NOTHING",
		slot, banner,
	); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *PostgresRepo) AddGroup(ctx context.Context, group model.GroupId) error {

	if _, err := r.db.ExecContext(ctx,
		"INSERT INTO groups(id) VALUES($1) ON CONFLICT(id) DO NOTHING",
		group,
	); err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepo) RemoveBanner(ctx context.Context, slot model.SlotId, banner model.BannerId) error {
	res, err := r.db.ExecContext(ctx,
		"DELETE FROM slot_banners WHERE slot_id=$1 AND banner_id=$2",
		slot, banner,
	)
	if err != nil {
		return err
	}
	if c, _ := res.RowsAffected(); c == 0 {
		return errors.New("banner not found in slot")
	}
	return nil
}

func (r *PostgresRepo) IncrementDisplay(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO stats(slot_id,banner_id,group_id,impressions,clicks) VALUES($1,$2,$3,1,0) ON CONFLICT(slot_id,banner_id,group_id) DO UPDATE SET impressions = stats.impressions + 1",
		slot, banner, group,
	)
	return err
}

func (r *PostgresRepo) IncrementClick(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO stats(slot_id,banner_id,group_id,impressions,clicks) VALUES($1,$2,$3,0,1) ON CONFLICT(slot_id,banner_id,group_id) DO UPDATE SET clicks = stats.clicks + 1",
		slot, banner, group,
	)
	return err
}

func (r *PostgresRepo) ListBanners(ctx context.Context, slot model.SlotId) ([]model.Banner, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT banner_id FROM slot_banners WHERE slot_id=$1",
		slot,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []model.Banner
	for rows.Next() {
		var id model.BannerId
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		banners = append(banners, model.Banner{Id: id})
	}
	return banners, nil
}

func (r *PostgresRepo) GetStats(ctx context.Context, slot model.SlotId, banner model.BannerId, group model.GroupId) (*model.Stats, error) {
	row := r.db.QueryRowContext(ctx,
		"SELECT impressions,clicks FROM stats WHERE slot_id=$1 AND banner_id=$2 AND group_id=$3",
		slot, banner, group,
	)
	var impr, clicks int64
	if err := row.Scan(&impr, &clicks); err != nil {
		if err == sql.ErrNoRows {
			return &model.Stats{Slot: slot, Banner: banner, Group: group, CountDisplay: 0, Clicks: 0}, nil
		}
		return nil, err
	}
	return &model.Stats{Slot: slot, Banner: banner, Group: group, CountDisplay: impr, Clicks: clicks}, nil
}
