package stores

import (
	"context"
	"fmt"
	"time"

	"github.com/soralabs/zen/cache"
	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

func NewFragmentStore(ctx context.Context, db *gorm.DB, fragmentTable db.FragmentTable) *FragmentStore {
	// Create cache with reasonable defaults for fragments
	cacheConfig := cache.Config{
		MaxSize:       1000,             // Store up to 1000 fragments
		TTL:           time.Minute * 30, // Cache fragments for 30 minutes
		CleanupPeriod: time.Minute,      // Clean up every minute
	}

	return &FragmentStore{
		Store: Store{
			db:  db,
			ctx: ctx,
		},
		fragmentTable: fragmentTable,
		cache:         cache.New(cacheConfig),
	}
}

func (f *FragmentStore) Create(fragment *db.Fragment) error {
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Create(fragment).Error
	if err != nil {
		return err
	}

	// Get full fragment with joins for cache
	fullFragment, err := f.GetByID(fragment.ID)
	if err != nil {
		return err
	}

	// Update cache
	cacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, fragment.ID))
	f.cache.Set(cacheKey, fullFragment)

	return nil
}

func (f *FragmentStore) Upsert(fragment *db.Fragment) error {
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Save(fragment).Error
	if err != nil {
		return err
	}

	// Refresh from DB with relationships
	fullFragment, err := f.GetByID(fragment.ID)
	if err != nil {
		return err
	}

	// Update cache
	cacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, fragment.ID))
	f.cache.Set(cacheKey, fullFragment)

	return nil
}

func (f *FragmentStore) GetByID(fragmentID id.ID) (*db.Fragment, error) {
	// Try cache first
	cacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, fragmentID))
	if cached, found := f.cache.Get(cacheKey); found {
		if fragment, ok := cached.(*db.Fragment); ok {
			return fragment, nil
		}
	}

	var fragment db.Fragment
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Joins("Actor").
		Joins("Session").
		Where(string(f.fragmentTable)+".id = ?", fragmentID).
		First(&fragment).Error
	if err != nil {
		return nil, err
	}

	if fragment.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	f.cache.Set(cacheKey, &fragment)
	return &fragment, nil
}

func (f *FragmentStore) GetBySession(sessionID id.ID, limit int) ([]db.Fragment, error) {
	var fragments []db.Fragment
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Joins("Actor").
		Joins("Session").
		Where(string(f.fragmentTable)+".session_id = ?", sessionID).
		Order(string(f.fragmentTable) + ".created_at DESC").
		Limit(limit).
		Find(&fragments).Error
	return fragments, err
}

func (f *FragmentStore) SearchSimilar(embedding pgvector.Vector, sessionID id.ID, limit int) ([]db.Fragment, error) {
	var fragments []db.Fragment
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Joins("Actor").
		Joins("Session").
		Select("*, ("+string(f.fragmentTable)+".embedding <=> ?) as similarity", embedding).
		Where(string(f.fragmentTable)+".session_id = ?", sessionID).
		Order("similarity").
		Limit(limit).
		Find(&fragments).Error
	return fragments, err
}

func (f *FragmentStore) DeleteByID(fragmentID id.ID) error {
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Delete(&db.Fragment{}, "id = ?", fragmentID).Error
	if err != nil {
		return err
	}

	// Remove from cache
	cacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, fragmentID))
	f.cache.Delete(cacheKey)

	return nil
}

func (f *FragmentStore) DeleteBySession(sessionID id.ID) error {
	return f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Delete(&db.Fragment{}, "session_id = ?", sessionID).Error
}

func (f *FragmentStore) GetRecentByManager(managerID id.ID, limit int) ([]db.Fragment, error) {
	var fragments []db.Fragment
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Joins("Actor").
		Joins("Session").
		Where(string(f.fragmentTable)+".manager_id = ?", managerID).
		Order(string(f.fragmentTable) + ".created_at DESC").
		Limit(limit).
		Find(&fragments).Error
	return fragments, err
}

func (f *FragmentStore) GetByActor(actorID id.ID, limit int) ([]db.Fragment, error) {
	var fragments []db.Fragment
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Joins("Actor").
		Joins("Session").
		Where(string(f.fragmentTable)+".actor_id = ?", actorID).
		Order(string(f.fragmentTable) + ".created_at DESC").
		Limit(limit).
		Find(&fragments).Error
	return fragments, err
}

func (f *FragmentStore) UpdateContent(fragmentID id.ID, content string) error {
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Where("id = ?", fragmentID).
		Update("content", content).Error
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, fragmentID))
	f.cache.Delete(cacheKey)

	return nil
}

func (f *FragmentStore) UpdateEmbedding(fragmentID id.ID, embedding []float32) error {
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Where("id = ?", fragmentID).
		Update("embedding", embedding).Error
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, fragmentID))
	f.cache.Delete(cacheKey)

	return nil
}

func (f *FragmentStore) UpdateMetadata(fragmentID id.ID, metadata map[string]interface{}) error {
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Where("id = ?", fragmentID).
		Update("metadata", metadata).Error
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, fragmentID))
	f.cache.Delete(cacheKey)

	return nil
}

func (f *FragmentStore) UpdateID(oldID id.ID, newID id.ID) error {
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Where("id = ?", oldID).
		Update("id", newID).Error
	if err != nil {
		return err
	}

	// Invalidate cache for both old and new IDs
	oldCacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, oldID))
	newCacheKey := cache.CacheKey(fmt.Sprintf("fragment:%s:%s", f.fragmentTable, newID))
	f.cache.Delete(oldCacheKey)
	f.cache.Delete(newCacheKey)

	return nil
}

func (f *FragmentStore) SearchByFilter(filter FragmentFilter) ([]db.Fragment, error) {
	query := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Joins("Actor").
		Joins("Session")

	// Apply basic filters
	if filter.ActorID != nil {
		query = query.Where(string(f.fragmentTable)+".actor_id = ?", *filter.ActorID)
	}
	if filter.SessionID != nil {
		query = query.Where(string(f.fragmentTable)+".session_id = ?", *filter.SessionID)
	}

	// Apply metadata filters
	for _, condition := range filter.Metadata {
		switch condition.Operator {
		case MetadataOpEquals:
			query = query.Where(string(f.fragmentTable)+".metadata->>? = ?", condition.Key, toString(condition.Value))
		case MetadataOpNotEquals:
			query = query.Where(string(f.fragmentTable)+".metadata->>? != ?", condition.Key, toString(condition.Value))
		case MetadataOpContains:
			query = query.Where(string(f.fragmentTable)+".metadata ?? ?", condition.Key)
		case MetadataOpIn:
			if values, ok := condition.Value.([]interface{}); ok {
				valueStrings := make([]string, len(values))
				for i, v := range values {
					valueStrings[i] = toString(v)
				}
				query = query.Where(string(f.fragmentTable)+".metadata->>? IN (?)", condition.Key, valueStrings)
			}
		}
	}

	// Apply time range filters
	if filter.StartTime != nil {
		query = query.Where(string(f.fragmentTable)+".created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where(string(f.fragmentTable)+".created_at <= ?", *filter.EndTime)
	}

	// Apply embedding similarity if provided
	if filter.Embedding != nil {
		query = query.Select("*, ("+string(f.fragmentTable)+".embedding <=> ?) as similarity", *filter.Embedding).
			Order("similarity")
	} else {
		query = query.Order(string(f.fragmentTable) + ".created_at DESC")
	}

	// Apply limit
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	var fragments []db.Fragment
	err := query.Find(&fragments).Error
	return fragments, err
}

func (f *FragmentStore) GetRecentSessionsByActor(actorID id.ID, limit int) ([]id.ID, error) {
	var sessionIDs []id.ID
	err := f.db.WithContext(f.ctx).
		Table(string(f.fragmentTable)).
		Select("DISTINCT "+string(f.fragmentTable)+".session_id").
		Where(string(f.fragmentTable)+".actor_id = ?", actorID).
		Order("MAX(" + string(f.fragmentTable) + ".created_at) DESC").
		Group(string(f.fragmentTable) + ".session_id").
		Limit(limit).
		Find(&sessionIDs).Error
	return sessionIDs, err
}
