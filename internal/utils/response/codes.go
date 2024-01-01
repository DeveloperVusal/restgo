package response

const (
	// When all right
	ErrorEmpty = 0
	// When don't activated a user account
	ErrorAccountActivate = 1
	// When don't match: password and confirm password
	ErrorAccountConfirmPassword = 2
	// When during registration, if a user account exist
	ErrorAccountExists = 3
	// When don't creating a user account
	ErrorAccountNotCreated = 4
)
