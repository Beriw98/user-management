package response

type UserIDResponse struct {
	ID string `json:"id"`
}

type UserResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}
