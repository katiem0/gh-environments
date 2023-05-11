package data

import "time"

type CreateEnvSecret struct {
	EncryptedValue string `json:"encrypted_vale"`
	KeyID          string `json:"key_id"`
}

type EnvPublicKey struct {
	KeyID string `json:"key_id"`
	Key   string `json:"string"`
}

type EnvSecret struct {
	TotalCount int      `json:"total_count"`
	Secrets    []Secret `json:"secrets"`
}

type Secret struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
