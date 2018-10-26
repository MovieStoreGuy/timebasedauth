package types

type Perception interface {

	// SetSessionCount will create that many sessions simultaneously to gather more timing data to see if we are progressing further
	SetSessionCount(activeSessions int) Perception

	// SetUser will configure the user to try login as.
	SetUser(username string) Perception

	// SetPasswordLength will stop the ensure that it will not go own forever and stop after it has reached the max length
	SetPasswordLength(length int) Perception

	// AttackHost takes the given URL and will try to login using timing based authentication attacks
	// and return the string that was able to login, other return an error about its failure.
	AttackHost(url string) (string, error)
}