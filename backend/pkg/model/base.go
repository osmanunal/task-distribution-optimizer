package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	UUID      uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
