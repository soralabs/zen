package stores

import (
	"context"

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"

	"gorm.io/gorm"
)

func NewSessionStore(ctx context.Context, db *gorm.DB) *SessionStore {
	return &SessionStore{
		Store: Store{
			db:  db,
			ctx: ctx,
		},
	}
}

func (cs *SessionStore) Create(session *db.Session) error {
	return cs.db.WithContext(cs.ctx).Create(session).Error
}

func (cs *SessionStore) Upsert(session *db.Session) error {
	return cs.db.WithContext(cs.ctx).Save(session).Error
}

func (cs *SessionStore) GetByID(id id.ID) (*db.Session, error) {
	var session db.Session
	if err := cs.db.WithContext(cs.ctx).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
