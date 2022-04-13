package user

import (
	"errors"

	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
	"gorm.io/gorm"
)

type RepositoryPostgres struct {
	db *gorm.DB
}

func NewRepositoryPostgres(db *gorm.DB) *RepositoryPostgres {
	r := &RepositoryPostgres{
		db: db,
	}

	r.db.AutoMigrate(&User{}) //nolint

	return r
}

func (r *RepositoryPostgres) Create(u *User) error {
	return r.db.FirstOrCreate(u).Error
}

func (r *RepositoryPostgres) Find(id entity.ID) (*User, error) {
	u := &User{}
	err := r.db.First(u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (r *RepositoryPostgres) ListAllWithApiEnabled() ([]*User, error) {
	var users []*User
	err := r.db.Find(&users, "api_enabled = true").Error

	return users, err
}

func (r *RepositoryPostgres) UpdateFaceitApiUsage(u *User) error {
	return r.db.Model(&u).Where("id = ?", u.ID).Update("ApiEnabled", u.ApiEnabled).Error
}
