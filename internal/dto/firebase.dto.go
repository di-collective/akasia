package dto

type FirebaseUser struct {
	UID           string                 `bson:"uid,omitempty" json:"uid,omitempty"`
	Email         string                 `bson:"email,omitempty" json:"email,omitempty"`
	EmailVerified bool                   `bson:"email_verified,omitempty" json:"email_verified,omitempty"`
	PhoneNumber   string                 `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	DisplayName   string                 `bson:"display_name,omitempty" json:"display_name,omitempty"`
	PhotoURL      string                 `bson:"photo_url,omitempty" json:"photo_url,omitempty"`
	Disabled      bool                   `bson:"disabled,omitempty" json:"disabled,omitempty"`
	Password      string                 `json:"password,omitempty"`
	CustomClaims  map[string]interface{} `json:"custom_claims,omitempty"`
}

type FirebaseClaims struct {
	Role   string `json:"x-hasura-default-role,omitempty"`
	UserID string `json:"x-hasura-user-id,omitempty"`
}
