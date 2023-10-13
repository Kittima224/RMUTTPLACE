package dto

type UserReview struct {
	ID    uint
	Name  string
	Image string
}
type ReviewResponse struct {
	UaserID int
	Comment string
	Rating  float32
}
