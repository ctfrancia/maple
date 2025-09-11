package types

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type Player struct {
	ID         *string `json:"id,omitempty"`
	First      *string `json:"first,omitempty"`
	Last       *string `json:"last,omitempty"`
	Gender     *Gender `json:"gender,omitempty"`
	FIDETitle  *string `json:"fide_title,omitempty"`
	LocalTitle *string `json:"local_title,omitempty"`
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	DOB        *string `json:"dob,omitempty"`
}
