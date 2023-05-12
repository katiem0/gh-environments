package data

import "time"

type CreateVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EnvVariables struct {
	TotalCount int        `json:"total_count"`
	Variables  []Variable `json:"variables"`
}

type ImportedVariable struct {
	RepositoryID    int
	RepositoryName  string
	EnvironmentName string
	Name            string `json:"name"`
	Value           string `json:"value"`
}

type Variable struct {
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
