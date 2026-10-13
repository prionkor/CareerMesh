package models

import (
	"encoding/json"
	"time"

	"uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PendingUser struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	TokenHash    string    `db:"token_hash"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
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

func (employmentType EmploymentType) IsValid() bool {
	switch employmentType {
	case EmploymentTypeFullTime, EmploymentTypePartTime, EmploymentTypeContract,
		EmploymentTypeFreelance, EmploymentTypeInternship, EmploymentTypeTemporary,
		EmploymentTypeVolunteer:
		return true
	default:
		return false
	}
}

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

type JobDescription struct {
	ID               uuid.UUID              `json:"id"`
	UserID           uuid.UUID              `json:"user_id"`
	CurrentVersionID uuid.UUID              `json:"current_version_id"`
	CurrentVersion   *JobDescriptionVersion `json:"current_version,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type JobDescriptionVersion struct {
	ID               uuid.UUID       `json:"id"`
	JobDescriptionID uuid.UUID       `json:"job_description_id"`
	VersionNumber    int             `json:"version_number"`
	Title            string          `json:"title"`
	Company          string          `json:"company"`
	CompanyWebsite   string          `json:"company_website"`
	SourceURL        string          `json:"source_url"`
	Location         string          `json:"location"`
	EmploymentType   *EmploymentType `json:"employment_type"`
	OriginalText     string          `json:"original_text"`
	CreatedAt        time.Time       `json:"created_at"`
}

// Date is a calendar date serialized as YYYY-MM-DD.
type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(time.DateOnly))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	parsed, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

type JobApplicationStatus string

const (
	JobApplicationStatusDraft        JobApplicationStatus = "draft"
	JobApplicationStatusApplied      JobApplicationStatus = "applied"
	JobApplicationStatusScreening    JobApplicationStatus = "screening"
	JobApplicationStatusInterviewing JobApplicationStatus = "interviewing"
	JobApplicationStatusOffer        JobApplicationStatus = "offer"
	JobApplicationStatusAccepted     JobApplicationStatus = "accepted"
	JobApplicationStatusRejected     JobApplicationStatus = "rejected"
	JobApplicationStatusWithdrawn    JobApplicationStatus = "withdrawn"
)

func (status JobApplicationStatus) IsValid() bool {
	switch status {
	case JobApplicationStatusDraft, JobApplicationStatusApplied, JobApplicationStatusScreening,
		JobApplicationStatusInterviewing, JobApplicationStatusOffer, JobApplicationStatusAccepted,
		JobApplicationStatusRejected, JobApplicationStatusWithdrawn:
		return true
	default:
		return false
	}
}

type JobApplication struct {
	ID               uuid.UUID            `json:"id"`
	UserID           uuid.UUID            `json:"user_id"`
	JobDescriptionID uuid.UUID            `json:"job_description_id"`
	Status           JobApplicationStatus `json:"status"`
	AppliedAt        *Date                `json:"applied_at"`
	Notes            string               `json:"notes"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
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
