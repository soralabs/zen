package stores

import (
	"context"
	"zen/internal/db"
	"zen/pkg/id"

	"gorm.io/gorm"
)

func NewActorStore(ctx context.Context, db *gorm.DB) *ActorStore {
	return &ActorStore{
		Store: Store{
			db:  db,
			ctx: ctx,
		},
	}
}

func (m *ActorStore) Create(actor *db.Actor) error {
	return m.db.WithContext(m.ctx).Create(actor).Error
}

func (m *ActorStore) Upsert(actor *db.Actor) error {
	return m.db.WithContext(m.ctx).Save(actor).Error
}

func (m *ActorStore) GetByID(id id.ID) (*db.Actor, error) {
	var actor db.Actor
	if err := m.db.WithContext(m.ctx).First(&actor, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &actor, nil
}

func (m *ActorStore) Update(actor *db.Actor) error {
	return m.db.WithContext(m.ctx).Save(actor).Error
}

func (m *ActorStore) DeleteByID(id id.ID) error {
	return m.db.WithContext(m.ctx).Delete(&db.Actor{}, "id = ?", id).Error
}

func (m *ActorStore) List(limit int) ([]db.Actor, error) {
	var actors []db.Actor
	err := m.db.WithContext(m.ctx).
		Limit(limit).
		Find(&actors).Error
	return actors, err
}

func (m *ActorStore) Search(query string, limit int) ([]db.Actor, error) {
	var actors []db.Actor
	err := m.db.WithContext(m.ctx).
		Where("name ILIKE ?", "%"+query+"%").
		Limit(limit).
		Find(&actors).Error
	return actors, err
}
