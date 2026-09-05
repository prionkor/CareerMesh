package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Profile struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Headline  string    `json:"headline"`
	Location  string    `json:"location"`
	Phone     string    `json:"phone"`
	Website   string    `json:"website"`
	GitHub    string    `json:"github"`
	LinkedIn  string    `json:"linkedin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EmploymentType string

const (
	EmploymentTypeFullTime   EmploymentType = "full_time"
	EmploymentTypePartTime   EmploymentType = "part_time"
	EmploymentTypeContract   EmploymentType = "contract"
	EmploymentTypeFreelance  EmploymentType = "freelance"
	EmploymentTypeInternship EmploymentType = "internship"
	EmploymentTypeTemporary  EmploymentType = "temporary"
	EmploymentTypeVolunteer  EmploymentType = "volunteer"
)

type Experience struct {
	ID             uuid.UUID      `json:"id"`
	UserID         uuid.UUID      `json:"user_id"`
	Company        string         `json:"company"`
	Website        string         `json:"website"`
	Title          string         `json:"title"`
	EmploymentType EmploymentType `json:"employment_type"`
	Location       string         `json:"location"`
	StartDate      time.Time      `json:"start_date"`
	EndDate        *time.Time     `json:"end_date"`
	Description    string         `json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type ExperienceBullet struct {
	ID           uuid.UUID `json:"id"`
	ExperienceID uuid.UUID `json:"experience_id"`
	Content      string    `json:"content"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Project struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	URL           string     `json:"url"`
	RepositoryURL string     `json:"repository_url"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Skill struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExperienceSkill struct {
	ExperienceID uuid.UUID `json:"experience_id"`
	SkillID      uuid.UUID `json:"skill_id"`
}

type ProjectSkill struct {
	ProjectID uuid.UUID `json:"project_id"`
	SkillID   uuid.UUID `json:"skill_id"`
}

type Education struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Institution string     `json:"institution"`
	Degree      string     `json:"degree"`
	Field       string     `json:"field"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Certification struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	Name         string     `json:"name"`
	Issuer       string     `json:"issuer"`
	IssueDate    *time.Time `json:"issue_date"`
	ExpiryDate   *time.Time `json:"expiry_date"`
	CredentialID string     `json:"credential_id"`
	URL          string     `json:"url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type LanguageProficiency string

const (
	LanguageProficiencyBeginner          LanguageProficiency = "beginner"
	LanguageProficiencyElementary        LanguageProficiency = "elementary"
	LanguageProficiencyIntermediate      LanguageProficiency = "intermediate"
	LanguageProficiencyUpperIntermediate LanguageProficiency = "upper_intermediate"
	LanguageProficiencyAdvanced          LanguageProficiency = "advanced"
	LanguageProficiencyNative            LanguageProficiency = "native"
)

type Language struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"user_id"`
	Name        string              `json:"name"`
	Proficiency LanguageProficiency `json:"proficiency"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
