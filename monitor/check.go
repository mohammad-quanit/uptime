package monitor

import (
	"context"
	"time"

	"encore.app/site"
	"encore.dev/cron"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var _ = cron.NewJob("check-all", cron.JobConfig{
	Title:    "Check all sites",
	Endpoint: CheckAll,
	Every:    1 * cron.Hour,
})

type IDBHandler interface {
	Create(interface{}) error
	Find(interface{}) error
}

func (s *Service) Create(value interface{}) error {
	return s.db.Create(value).Error
}

func (s *Service) Find(value interface{}) error {
	return s.db.Find(value).Error
}

type Check struct {
	gorm.Model
	SiteID    int64     `json:"site_id"`
	Up        bool      `json:"up"`
	CheckedAt time.Time `json:"checked_at"`
}

type SListChecked struct {
	Checked []*Check `json:"checked"`
}

// Get All Checked list.
//
//encore:api public method=GET path=/check
func (s *Service) GetChecks(ctx context.Context) (*SListChecked, error) {
	var checked []*Check

	if err := s.Find(&checked); err != nil {
		return nil, err
	}

	return &SListChecked{Checked: checked}, nil
}

// Check checks a single site.
//
//encore:api public method=POST path=/check/:siteID
func (s *Service) AddCheck(ctx context.Context, siteID int) error {
	site, err := site.Get(ctx, siteID)
	if err != nil {
		return err
	}
	dbService := &Service{db: s.db}

	return check(ctx, site, dbService)
}

// CheckAll checks all sites.
//
//encore:api public method=POST path=/checkall
func (s *Service) CheckAll(ctx context.Context) error {
	resp, err := site.List(ctx)
	if err != nil {
		return err
	}

	// Check up to 8 sites concurrently.
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(8)
	for _, site := range resp.Sites {
		site := site
		eg.Go(func() error {
			return check(ctx, site, &Service{db: s.db})
		})
	}
	return eg.Wait()
}

// Refactored Function
func check(ctx context.Context, site *site.Site, s IDBHandler) error {
	response, err := Ping(ctx, site.URL)
	if err != nil {
		return err
	}

	check := Check{
		SiteID:    int64(site.ID),
		Up:        response.Up,
		CheckedAt: time.Now(),
	}

	if err := s.Create(&check); err != nil {
		return err
	}

	return nil
}
