package domain

type Profile struct {
	ID    uint
	Name  string
	Email string
}

type ProfileFeed struct {
	Profile
	Pictures []Picture
}
