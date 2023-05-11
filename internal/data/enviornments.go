package data

type BranchPolicies struct {
	TotalCount     int            `json:"total_count"`
	BranchPolicies []BranchPolicy `json:"branch_policies"`
}

type BranchPolicy struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DeploymentPolicy struct {
	ProtectedBranches bool `json:"protected_branches"`
	CustomPolicies    bool `json:"custom_branch_policies"`
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

type Rules struct {
	Type      string      `json:"type"`
	WaitTimer int         `json:"wait_timer"`
	Reviewers []Reviewers `json:"reviewers"`
}

type Reviewers struct {
	Type     string   `json:"type"`
	Reviewer Reviewer `json:"reviewer"`
}

type Reviewer struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}
