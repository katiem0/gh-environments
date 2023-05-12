package data

import "time"

type CreateEnvSecret struct {
	EncryptedValue string `json:"encrypted_value"`
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

type ImportedSecret struct {
	RepositoryID    int
	RepositoryName  string
	EnvironmentName string
	Name            string `json:"name"`
	Value           string `json:"value"`
}

type PublicKey struct {
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
}

type Secret struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
