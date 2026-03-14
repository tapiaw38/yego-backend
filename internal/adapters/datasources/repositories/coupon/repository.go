package coupon

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"yego/internal/domain"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
)

// Repository defines the interface for coupon operations
type Repository interface {
	Create(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, apperrors.ApplicationError)
	GetByID(ctx context.Context, id string) (*domain.Coupon, apperrors.ApplicationError)
	GetByCode(ctx context.Context, code string) (*domain.Coupon, apperrors.ApplicationError)
	List(ctx context.Context) ([]*domain.Coupon, apperrors.ApplicationError)
	Update(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, apperrors.ApplicationError)
	Delete(ctx context.Context, id string) apperrors.ApplicationError
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, c *domain.Coupon) (*domain.Coupon, apperrors.ApplicationError) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Code = strings.ToUpper(strings.TrimSpace(c.Code))

	query := `
		INSERT INTO coupons (
			id, code, description, discount_type, discount_value,
			max_uses, current_uses, usage_limit_per_user, min_order_amount,
			valid_from, valid_until, active, icon_url, cover_url,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.Code, c.Description, string(c.DiscountType), c.DiscountValue,
		c.MaxUses, c.CurrentUses, c.UsageLimitPerUser, c.MinOrderAmount,
		c.ValidFrom, c.ValidUntil, c.Active, c.IconURL, c.CoverURL,
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, apperrors.NewApplicationError(mappings.CouponDuplicateCodeError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.CouponCreateError, err)
	}
	return c, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*domain.Coupon, apperrors.ApplicationError) {
	query := `
		SELECT id, code, description, discount_type, discount_value,
		       max_uses, current_uses, usage_limit_per_user, min_order_amount,
		       valid_from, valid_until, active, icon_url, cover_url,
		       created_at, updated_at
		FROM coupons WHERE id = $1
	`
	c, appErr := scanCoupon(r.db.QueryRowContext(ctx, query, id))
	if appErr != nil {
		return nil, appErr
	}
	return c, nil
}

func (r *repository) GetByCode(ctx context.Context, code string) (*domain.Coupon, apperrors.ApplicationError) {
	query := `
		SELECT id, code, description, discount_type, discount_value,
		       max_uses, current_uses, usage_limit_per_user, min_order_amount,
		       valid_from, valid_until, active, icon_url, cover_url,
		       created_at, updated_at
		FROM coupons WHERE code = $1
	`
	c, appErr := scanCoupon(r.db.QueryRowContext(ctx, query, strings.ToUpper(strings.TrimSpace(code))))
	if appErr != nil {
		return nil, appErr
	}
	return c, nil
}

func (r *repository) List(ctx context.Context) ([]*domain.Coupon, apperrors.ApplicationError) {
	query := `
		SELECT id, code, description, discount_type, discount_value,
		       max_uses, current_uses, usage_limit_per_user, min_order_amount,
		       valid_from, valid_until, active, icon_url, cover_url,
		       created_at, updated_at
		FROM coupons ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CouponListError, err)
	}
	defer rows.Close()

	var coupons []*domain.Coupon
	for rows.Next() {
		c, appErr := scanCouponRow(rows)
		if appErr != nil {
			return nil, appErr
		}
		coupons = append(coupons, c)
	}
	if err = rows.Err(); err != nil {
		return nil, apperrors.NewApplicationError(mappings.CouponListError, err)
	}
	return coupons, nil
}

func (r *repository) Update(ctx context.Context, c *domain.Coupon) (*domain.Coupon, apperrors.ApplicationError) {
	c.UpdatedAt = time.Now()
	c.Code = strings.ToUpper(strings.TrimSpace(c.Code))

	query := `
		UPDATE coupons SET
			code = $2, description = $3, discount_type = $4, discount_value = $5,
			max_uses = $6, usage_limit_per_user = $7, min_order_amount = $8,
			valid_from = $9, valid_until = $10, active = $11,
			icon_url = $12, cover_url = $13, updated_at = $14
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query,
		c.ID, c.Code, c.Description, string(c.DiscountType), c.DiscountValue,
		c.MaxUses, c.UsageLimitPerUser, c.MinOrderAmount,
		c.ValidFrom, c.ValidUntil, c.Active,
		c.IconURL, c.CoverURL, c.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, apperrors.NewApplicationError(mappings.CouponDuplicateCodeError, err)
		}
		return nil, apperrors.NewApplicationError(mappings.CouponUpdateError, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, apperrors.NewApplicationError(mappings.CouponNotFoundError, nil)
	}
	return c, nil
}

func (r *repository) Delete(ctx context.Context, id string) apperrors.ApplicationError {
	res, err := r.db.ExecContext(ctx, `DELETE FROM coupons WHERE id = $1`, id)
	if err != nil {
		return apperrors.NewApplicationError(mappings.CouponDeleteError, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperrors.NewApplicationError(mappings.CouponNotFoundError, nil)
	}
	return nil
}

// scanCoupon scans a single row from QueryRowContext
type rowScanner interface {
	Scan(dest ...any) error
}

func scanCoupon(row *sql.Row) (*domain.Coupon, apperrors.ApplicationError) {
	var c domain.Coupon
	var description, iconURL, coverURL sql.NullString
	var maxUses sql.NullInt64
	var minOrderAmount sql.NullFloat64
	var validFrom, validUntil sql.NullTime
	var discountType string

	err := row.Scan(
		&c.ID, &c.Code, &description, &discountType, &c.DiscountValue,
		&maxUses, &c.CurrentUses, &c.UsageLimitPerUser, &minOrderAmount,
		&validFrom, &validUntil, &c.Active, &iconURL, &coverURL,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewApplicationError(mappings.CouponNotFoundError, err)
	}
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CouponGetError, err)
	}

	c.DiscountType = domain.DiscountType(discountType)
	if description.Valid {
		c.Description = &description.String
	}
	if maxUses.Valid {
		v := int(maxUses.Int64)
		c.MaxUses = &v
	}
	if minOrderAmount.Valid {
		c.MinOrderAmount = &minOrderAmount.Float64
	}
	if validFrom.Valid {
		c.ValidFrom = &validFrom.Time
	}
	if validUntil.Valid {
		c.ValidUntil = &validUntil.Time
	}
	if iconURL.Valid {
		c.IconURL = &iconURL.String
	}
	if coverURL.Valid {
		c.CoverURL = &coverURL.String
	}
	return &c, nil
}

func scanCouponRow(rows *sql.Rows) (*domain.Coupon, apperrors.ApplicationError) {
	var c domain.Coupon
	var description, iconURL, coverURL sql.NullString
	var maxUses sql.NullInt64
	var minOrderAmount sql.NullFloat64
	var validFrom, validUntil sql.NullTime
	var discountType string

	err := rows.Scan(
		&c.ID, &c.Code, &description, &discountType, &c.DiscountValue,
		&maxUses, &c.CurrentUses, &c.UsageLimitPerUser, &minOrderAmount,
		&validFrom, &validUntil, &c.Active, &iconURL, &coverURL,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CouponListError, err)
	}

	c.DiscountType = domain.DiscountType(discountType)
	if description.Valid {
		c.Description = &description.String
	}
	if maxUses.Valid {
		v := int(maxUses.Int64)
		c.MaxUses = &v
	}
	if minOrderAmount.Valid {
		c.MinOrderAmount = &minOrderAmount.Float64
	}
	if validFrom.Valid {
		c.ValidFrom = &validFrom.Time
	}
	if validUntil.Valid {
		c.ValidUntil = &validUntil.Time
	}
	if iconURL.Valid {
		c.IconURL = &iconURL.String
	}
	if coverURL.Valid {
		c.CoverURL = &coverURL.String
	}
	return &c, nil
}
