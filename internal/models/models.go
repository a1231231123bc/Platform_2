package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// ============================================================
// Enums
// ============================================================

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleManager  UserRole = "MANAGER"
	UserRoleObserver UserRole = "OBSERVER"
)

type ContractorType string

const (
	ContractorTypeIndividual ContractorType = "INDIVIDUAL"
	ContractorTypeIP         ContractorType = "IP"
	ContractorTypeLLC        ContractorType = "LLC"
)

type JobStatus string

const (
	JobStatusDraft      JobStatus = "DRAFT"
	JobStatusActive     JobStatus = "ACTIVE"
	JobStatusInProgress JobStatus = "IN_PROGRESS"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusCancelled  JobStatus = "CANCELLED"
)

type DispatchStatus string

const (
	DispatchStatusPending   DispatchStatus = "PENDING"
	DispatchStatusSending   DispatchStatus = "SENDING"
	DispatchStatusCompleted DispatchStatus = "COMPLETED"
	DispatchStatusPaused    DispatchStatus = "PAUSED"
	DispatchStatusCancelled DispatchStatus = "CANCELLED"
)

type OfferStatus string

const (
	OfferStatusSent     OfferStatus = "SENT"
	OfferStatusAccepted OfferStatus = "ACCEPTED"
	OfferStatusDeclined OfferStatus = "DECLINED"
	OfferStatusExpired  OfferStatus = "EXPIRED"
)

type AssignmentStatus string

const (
	AssignmentStatusAssigned  AssignmentStatus = "ASSIGNED"
	AssignmentStatusConfirmed AssignmentStatus = "CONFIRMED"
	AssignmentStatusInWork    AssignmentStatus = "IN_WORK"
	AssignmentStatusCompleted AssignmentStatus = "COMPLETED"
	AssignmentStatusCancelled AssignmentStatus = "CANCELLED"
)

type CommunicationDirection string

const (
	CommunicationDirectionOutgoing CommunicationDirection = "OUTGOING"
	CommunicationDirectionIncoming CommunicationDirection = "INCOMING"
)

type RatingAuthorType string

const (
	RatingAuthorTypeCustomer   RatingAuthorType = "CUSTOMER"
	RatingAuthorTypeContractor RatingAuthorType = "CONTRACTOR"
)

// ============================================================
// Models
// ============================================================

type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	INN          *string   `json:"inn"`
	ContactEmail *string   `json:"contactEmail"`
	ContactPhone *string   `json:"contactPhone"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Name           string    `json:"name"`
	Role           UserRole  `json:"role"`
	OrganizationID string    `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Contractor struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Phone             string         `json:"phone"`
	Email             *string        `json:"email"`
	Type              ContractorType `json:"type"`
	Regions           []string       `json:"regions"`
	EquipmentTypes    []string       `json:"equipmentTypes"`
	PriceExpectations *string        `json:"priceExpectations"`
	IsAvailable       bool           `json:"isAvailable"`
	ConsentGiven      bool           `json:"consentGiven"`
	ConsentDate       *time.Time     `json:"consentDate"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

type Job struct {
	ID              string           `json:"id"`
	Title           string           `json:"title"`
	Description     *string          `json:"description"`
	Region          string           `json:"region"`
	EquipmentType   *string          `json:"equipmentType"`
	Volume          *string          `json:"volume"`
	Deadline        *time.Time       `json:"deadline"`
	Price           *decimal.Decimal `json:"price"`
	Conditions      *string          `json:"conditions"`
	Status          JobStatus        `json:"status"`
	OrganizationID  string           `json:"organizationId"`
	CreatedByUserID string           `json:"createdByUserId"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

type JobWithCounts struct {
	Job
	DispatchCount   int `json:"dispatchCount"`
	OfferCount      int `json:"offerCount"`
	AssignmentCount int `json:"assignmentCount"`
}

type JobDispatch struct {
	ID        string         `json:"id"`
	JobID     string         `json:"jobId"`
	Wave      int            `json:"wave"`
	Status    DispatchStatus `json:"status"`
	Limit     int            `json:"limit"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type JobOffer struct {
	ID           string      `json:"id"`
	JobID        string      `json:"jobId"`
	DispatchID   string      `json:"dispatchId"`
	ContractorID string      `json:"contractorId"`
	Status       OfferStatus `json:"status"`
	Token        string      `json:"token"`
	SentAt       time.Time   `json:"sentAt"`
	RespondedAt  *time.Time  `json:"respondedAt"`
}

type JobOfferWithJob struct {
	JobOffer
	Job Job `json:"job"`
}

type Assignment struct {
	ID           string           `json:"id"`
	JobID        string           `json:"jobId"`
	ContractorID string           `json:"contractorId"`
	Status       AssignmentStatus `json:"status"`
	AssignedAt   time.Time        `json:"assignedAt"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

type AssignmentWithContractor struct {
	Assignment
	Contractor Contractor `json:"contractor"`
}

type AssignmentWithJob struct {
	Assignment
	Job Job `json:"job"`
}

type CommunicationLog struct {
	ID           string                 `json:"id"`
	JobID        *string                `json:"jobId"`
	ContractorID string                 `json:"contractorId"`
	Channel      string                 `json:"channel"`
	Message      string                 `json:"message"`
	Direction    CommunicationDirection `json:"direction"`
	CreatedAt    time.Time              `json:"createdAt"`
}

type ContractorBlacklist struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	ContractorID   string    `json:"contractorId"`
	Reason         *string   `json:"reason"`
	CreatedAt      time.Time `json:"createdAt"`
}

type ContractorBlacklistWithContractor struct {
	ContractorBlacklist
	Contractor Contractor `json:"contractor"`
}

type Rating struct {
	ID                 string           `json:"id"`
	JobID              string           `json:"jobId"`
	ContractorID       string           `json:"contractorId"`
	AuthorContractorID *string          `json:"authorContractorId"`
	OrganizationID     *string          `json:"organizationId"`
	Score              int              `json:"score"`
	Comment            *string          `json:"comment"`
	AuthorType         RatingAuthorType `json:"authorType"`
	CreatedAt          time.Time        `json:"createdAt"`
}

type RatingWithJob struct {
	Rating
	Job Job `json:"job"`
}

type RatingWithContractor struct {
	Rating
	Contractor Contractor `json:"contractor"`
}
