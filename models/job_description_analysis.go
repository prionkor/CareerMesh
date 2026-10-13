package models

import (
	"fmt"
	"math"
	"time"

	"uuid"
)

type JobDescriptionAnalysisStatus string

const (
	JobDescriptionAnalysisStatusPending   JobDescriptionAnalysisStatus = "pending"
	JobDescriptionAnalysisStatusCompleted JobDescriptionAnalysisStatus = "completed"
	JobDescriptionAnalysisStatusFailed    JobDescriptionAnalysisStatus = "failed"
)

func (status JobDescriptionAnalysisStatus) IsValid() bool {
	switch status {
	case JobDescriptionAnalysisStatusPending, JobDescriptionAnalysisStatusCompleted, JobDescriptionAnalysisStatusFailed:
		return true
	default:
		return false
	}
}

const JobDescriptionAnalysisSchemaVersion = "1"

type JobDescriptionAnalysis struct {
	ID                      uuid.UUID                     `json:"id"`
	JobDescriptionVersionID uuid.UUID                     `json:"job_description_version_id"`
	SchemaVersion           string                        `json:"schema_version"`
	PromptVersion           string                        `json:"prompt_version"`
	Provider                string                        `json:"provider"`
	Model                   string                        `json:"model"`
	Status                  JobDescriptionAnalysisStatus  `json:"status"`
	Result                  *JobDescriptionAnalysisResult `json:"result"`
	Error                   *string                       `json:"error"`
	CreatedAt               time.Time                     `json:"created_at"`
	CompletedAt             *time.Time                    `json:"completed_at"`
}

type JobDescriptionAnalysisResult struct {
	SchemaVersion    string                                 `json:"schema_version"`
	Requirements     []JobDescriptionAnalysisRequirement    `json:"requirements"`
	Responsibilities []JobDescriptionAnalysisResponsibility `json:"responsibilities"`
}

type JobDescriptionAnalysisRequirementCategory string

const (
	JobDescriptionAnalysisRequirementCategorySkill         JobDescriptionAnalysisRequirementCategory = "skill"
	JobDescriptionAnalysisRequirementCategoryExperience    JobDescriptionAnalysisRequirementCategory = "experience"
	JobDescriptionAnalysisRequirementCategoryEducation     JobDescriptionAnalysisRequirementCategory = "education"
	JobDescriptionAnalysisRequirementCategoryCertification JobDescriptionAnalysisRequirementCategory = "certification"
	JobDescriptionAnalysisRequirementCategoryLanguage      JobDescriptionAnalysisRequirementCategory = "language"
	JobDescriptionAnalysisRequirementCategoryOther         JobDescriptionAnalysisRequirementCategory = "other"
)

type JobDescriptionAnalysisRequirementPriority string

const (
	JobDescriptionAnalysisRequirementPriorityRequired    JobDescriptionAnalysisRequirementPriority = "required"
	JobDescriptionAnalysisRequirementPriorityPreferred   JobDescriptionAnalysisRequirementPriority = "preferred"
	JobDescriptionAnalysisRequirementPriorityUnspecified JobDescriptionAnalysisRequirementPriority = "unspecified"
)

type JobDescriptionAnalysisAssertion string

const (
	JobDescriptionAnalysisAssertionExplicit JobDescriptionAnalysisAssertion = "explicit"
	JobDescriptionAnalysisAssertionInferred JobDescriptionAnalysisAssertion = "inferred"
)

type JobDescriptionAnalysisRequirement struct {
	Name        string                                    `json:"name"`
	Description string                                    `json:"description"`
	Category    JobDescriptionAnalysisRequirementCategory `json:"category"`
	Priority    JobDescriptionAnalysisRequirementPriority `json:"priority"`
	Assertion   JobDescriptionAnalysisAssertion           `json:"assertion"`
	Weight      float64                                   `json:"weight"`
	Evidence    []JobDescriptionAnalysisEvidence          `json:"evidence"`
}

type JobDescriptionAnalysisResponsibilityCategory string

const (
	JobDescriptionAnalysisResponsibilityCategoryTechnical     JobDescriptionAnalysisResponsibilityCategory = "technical"
	JobDescriptionAnalysisResponsibilityCategoryOperational   JobDescriptionAnalysisResponsibilityCategory = "operational"
	JobDescriptionAnalysisResponsibilityCategoryCollaboration JobDescriptionAnalysisResponsibilityCategory = "collaboration"
	JobDescriptionAnalysisResponsibilityCategoryLeadership    JobDescriptionAnalysisResponsibilityCategory = "leadership"
	JobDescriptionAnalysisResponsibilityCategoryOther         JobDescriptionAnalysisResponsibilityCategory = "other"
)

type JobDescriptionAnalysisResponsibility struct {
	Name        string                                       `json:"name"`
	Description string                                       `json:"description"`
	Category    JobDescriptionAnalysisResponsibilityCategory `json:"category"`
	Assertion   JobDescriptionAnalysisAssertion              `json:"assertion"`
	Weight      float64                                      `json:"weight"`
	Evidence    []JobDescriptionAnalysisEvidence             `json:"evidence"`
}

type JobDescriptionAnalysisEvidence struct {
	Text string `json:"text"`
}

func (result JobDescriptionAnalysisResult) Validate() error {
	if result.SchemaVersion != JobDescriptionAnalysisSchemaVersion {
		return fmt.Errorf("unsupported job description analysis schema version %q", result.SchemaVersion)
	}
	if result.Requirements == nil {
		return fmt.Errorf("job description analysis requirements must be an array")
	}
	if result.Responsibilities == nil {
		return fmt.Errorf("job description analysis responsibilities must be an array")
	}
	for index, requirement := range result.Requirements {
		if !validRequirementCategory(requirement.Category) {
			return fmt.Errorf("requirement %d has invalid category %q", index, requirement.Category)
		}
		if !validRequirementPriority(requirement.Priority) {
			return fmt.Errorf("requirement %d has invalid priority %q", index, requirement.Priority)
		}
		if !validAnalysisAssertion(requirement.Assertion) {
			return fmt.Errorf("requirement %d has invalid assertion %q", index, requirement.Assertion)
		}
		if !validAnalysisWeight(requirement.Weight) {
			return fmt.Errorf("requirement %d has invalid weight %v", index, requirement.Weight)
		}
		if requirement.Evidence == nil {
			return fmt.Errorf("requirement %d evidence must be an array", index)
		}
	}
	for index, responsibility := range result.Responsibilities {
		if !validResponsibilityCategory(responsibility.Category) {
			return fmt.Errorf("responsibility %d has invalid category %q", index, responsibility.Category)
		}
		if !validAnalysisAssertion(responsibility.Assertion) {
			return fmt.Errorf("responsibility %d has invalid assertion %q", index, responsibility.Assertion)
		}
		if !validAnalysisWeight(responsibility.Weight) {
			return fmt.Errorf("responsibility %d has invalid weight %v", index, responsibility.Weight)
		}
		if responsibility.Evidence == nil {
			return fmt.Errorf("responsibility %d evidence must be an array", index)
		}
	}
	return nil
}

func validRequirementCategory(category JobDescriptionAnalysisRequirementCategory) bool {
	switch category {
	case JobDescriptionAnalysisRequirementCategorySkill, JobDescriptionAnalysisRequirementCategoryExperience,
		JobDescriptionAnalysisRequirementCategoryEducation, JobDescriptionAnalysisRequirementCategoryCertification,
		JobDescriptionAnalysisRequirementCategoryLanguage, JobDescriptionAnalysisRequirementCategoryOther:
		return true
	default:
		return false
	}
}

func validRequirementPriority(priority JobDescriptionAnalysisRequirementPriority) bool {
	switch priority {
	case JobDescriptionAnalysisRequirementPriorityRequired, JobDescriptionAnalysisRequirementPriorityPreferred,
		JobDescriptionAnalysisRequirementPriorityUnspecified:
		return true
	default:
		return false
	}
}

func validResponsibilityCategory(category JobDescriptionAnalysisResponsibilityCategory) bool {
	switch category {
	case JobDescriptionAnalysisResponsibilityCategoryTechnical, JobDescriptionAnalysisResponsibilityCategoryOperational,
		JobDescriptionAnalysisResponsibilityCategoryCollaboration, JobDescriptionAnalysisResponsibilityCategoryLeadership,
		JobDescriptionAnalysisResponsibilityCategoryOther:
		return true
	default:
		return false
	}
}

func validAnalysisAssertion(assertion JobDescriptionAnalysisAssertion) bool {
	return assertion == JobDescriptionAnalysisAssertionExplicit || assertion == JobDescriptionAnalysisAssertionInferred
}

func validAnalysisWeight(weight float64) bool {
	return !math.IsNaN(weight) && !math.IsInf(weight, 0) && weight >= 0 && weight <= 1
}
