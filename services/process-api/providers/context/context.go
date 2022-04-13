package context

import (
	"github.com/google/uuid"
)

type Context struct {
	UUID      uuid.UUID `json:"uuid,omitempty"`
	UserID    uint      `json:"user_id"`
	UserRoles []string  `json:"user_roles"`
	TenantRef string    `json:"tenant_ref"`
}
