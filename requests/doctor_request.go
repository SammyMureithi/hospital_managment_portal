package request

type DoctorRequest struct {
	FirstName string `json:"first_name" bson:"first_name" validate:"required,min=3,max=20"`
	LastName  string `json:"last_name" bson:"last_name" validate:"required,min=3,max=20"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
	Phone     string `json:"phone" bson:"phone" validate:"required,len=10"`
}
