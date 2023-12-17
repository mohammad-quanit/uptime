package monitor

import (
	"context"
	"time"

	"encore.app/site"
)

type Check struct {
	SiteID    int64     `json:"site_id"`
	Up        bool      `json:"up"`
	CheckedAt time.Time `json:"checked_at"`
}

type ListChecked struct {
	Checked []*Check `json:"checked"`
}

// Check checks a single site.
//
//encore:api public method=POST path=/check/:siteID
func (s *Service) AddCheck(ctx context.Context, siteID int) error {
	site, err := site.Get(ctx, siteID)
	if err != nil {
		return err
	}

	response, err := Ping(ctx, site.URL)
	if err != nil {
		return err
	}

	check := Check{
		SiteID:    int64(site.ID),
		Up:        response.Up,
		CheckedAt: time.Now(),
	}

	if err := s.db.Create(&check).Error; err != nil {
		return err
	}

	return nil
}

// Get All Checked list.
//
//encore:api public method=GET path=/check
func (s *Service) GetCheck(ctx context.Context) (*ListChecked, error) {
	var checked []*Check

	if err := s.db.Find(&checked).Error; err != nil {
		return nil, err
	}

	return &ListChecked{Checked: checked}, nil
}
