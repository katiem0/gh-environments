package data

import "time"

type RepoInfo struct {
	DatabaseId int       `json:"databaseId"`
	Name       string    `json:"name"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Visibility string    `json:"visibility"`
}

type ReposQuery struct {
	Organization struct {
		Repositories struct {
			TotalCount int
			Nodes      []RepoInfo
			PageInfo   struct {
				EndCursor   string
				HasNextPage bool
			}
		} `graphql:"repositories(first: 100, after: $endCursor)"`
	} `graphql:"organization(login: $owner)"`
}

type RepoSingleQuery struct {
	Repository RepoInfo `graphql:"repository(owner: $owner, name: $name)"`
}

type ScopedRepository struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ScopedResponse struct {
	TotalCount   int                `json:"total_count"`
	Repositories []ScopedRepository `json:"repositories"`
}
