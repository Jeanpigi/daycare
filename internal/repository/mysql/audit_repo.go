package mysqlrepo

import (
	"context"
	"database/sql"
)

type AuditRepo struct{ db *sql.DB }

func NewAuditRepo(db *sql.DB) *AuditRepo { return &AuditRepo{db: db} }

func (r *AuditRepo) Insert(ctx context.Context, actorID uint64, action, entityType string, entityID *uint64, beforeJSON, afterJSON []byte) error {
	q := `INSERT INTO audit_log (actor_user_id, action, entity_type, entity_id, before_json, after_json)
	      VALUES (?, ?, ?, ?, ?, ?)`
	var eID any = nil
	if entityID != nil {
		eID = *entityID
	}
	_, err := r.db.ExecContext(ctx, q, actorID, action, entityType, eID, nullBytes(beforeJSON), nullBytes(afterJSON))
	return err
}

func nullBytes(b []byte) any {
	if len(b) == 0 {
		return nil
	}
	return b
}
