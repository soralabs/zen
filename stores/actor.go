package stores

import (
	"context"

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"

	"gorm.io/gorm"
)

// NewActorStore returns a new ActorStore initialized with the provided context and DB connection
func NewActorStore(ctx context.Context, db *gorm.DB) *ActorStore {
	return &ActorStore{
		Store: Store{
			db:  db,
			ctx: ctx,
		},
	}
}

// Create inserts a new Actor record into the database
func (m *ActorStore) Create(actor *db.Actor) error {
	return m.db.WithContext(m.ctx).Create(actor).Error
}

// Upsert creates or updates an Actor record based on its primary key
func (m *ActorStore) Upsert(actor *db.Actor) error {
	return m.db.WithContext(m.ctx).Save(actor).Error
}

// GetByID retrieves a single Actor by its ID
func (m *ActorStore) GetByID(id id.ID) (*db.Actor, error) {
	var actor db.Actor
	if err := m.db.WithContext(m.ctx).First(&actor, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &actor, nil
}

// Update modifies an existing Actor record in the database
func (m *ActorStore) Update(actor *db.Actor) error {
	return m.db.WithContext(m.ctx).Save(actor).Error
}

// DeleteByID removes an Actor record from the database by ID
func (m *ActorStore) DeleteByID(id id.ID) error {
	return m.db.WithContext(m.ctx).Delete(&db.Actor{}, "id = ?", id).Error
}

// List returns a slice of Actors up to the specified limit
func (m *ActorStore) List(limit int) ([]db.Actor, error) {
	var actors []db.Actor
	err := m.db.WithContext(m.ctx).
		Limit(limit).
		Find(&actors).Error
	return actors, err
}

// Search returns a slice of Actors whose names match the given query, up to the limit
func (m *ActorStore) Search(query string, limit int) ([]db.Actor, error) {
	var actors []db.Actor
	err := m.db.WithContext(m.ctx).
		Where("name ILIKE ?", "%"+query+"%").
		Limit(limit).
		Find(&actors).Error
	return actors, err
}
