package data

type BranchPolicies struct {
	TotalCount     int            `json:"total_count"`
	BranchPolicies []BranchPolicy `json:"branch_policies"`
}

type BranchPolicy struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateEnvironment struct {
	WaitTimer              int               `json:"wait_timer"`
	PreventSelfReview      bool              `json:"prevent_self_review"`
	Reviewers              []CreateReviewer  `json:"reviewers"`
	DeploymentBranchPolicy *DeploymentPolicy `json:"deployment_branch_policy"`
}

type CreateReviewer struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

type CreateDeploymentBranch struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type DeploymentApp struct {
	IntegrationID int    `json:"id"`
	Slug          string `json:"slug"`
}

type DeploymentPolicy struct {
	ProtectedBranches bool `json:"protected_branches"`
	CustomPolicies    bool `json:"custom_branch_policies"`
}

type DeploymentProtectionPolicy struct {
	TotalCount            int                             `json:"total_count"`
	CustomDeploymentRules []DeploymentProtectionPolicyApp `json:"custom_deployment_protection_rules"`
}

type DeploymentProtectionPolicyApp struct {
	PolicyID int           `json:"id"`
	Enabled  bool          `json:"enabled"`
	App      DeploymentApp `json:"app"`
}

type EnvResponse struct {
	TotalCount   int           `json:"total_count"`
	Environments []Environment `json:"environments"`
}

type Environment struct {
	Name             string            `json:"name"`
	AdminByPass      bool              `json:"can_admins_bypass"`
	ProtectionRules  []Rules           `json:"protection_rules"`
	DeploymentPolicy *DeploymentPolicy `json:"deployment_branch_policy"`
}

type ImportedEnvironment struct {
	RepositoryName    string
	RepositoryID      int
	EnvironmentName   string
	AdminBypass       string
	WaitTimer         int
	Reviewers         []Reviewers
	PreventSelfReview bool
	DeploymentPolicy  string
	Branches          []CreateDeploymentBranch
}

type Rules struct {
	Type              string      `json:"type"`
	WaitTimer         int         `json:"wait_timer"`
	Reviewers         []Reviewers `json:"reviewers"`
	PreventSelfReview bool        `json:"prevent_self_review"`
}

type Reviewers struct {
	Type     string   `json:"type"`
	Reviewer Reviewer `json:"reviewer"`
}

type Reviewer struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}
