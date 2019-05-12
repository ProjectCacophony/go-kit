package interfaces

type PermissionInterface interface {
	Name() string
	Match(
		state State,
		userID string,
		channelID string,
		dm bool,
	) bool
}
