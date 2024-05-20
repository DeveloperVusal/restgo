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
	// When a user account already activated
	ErrorAccountAlreadyActivate = 5
	// When time activation out
	ErrorAccountActivateTimeout = 6
	// When don't match of confirm code
	ErrorAccountInvalidCode = 6
	// When not found a user account
	ErrorAccountNotFound = 7
)
