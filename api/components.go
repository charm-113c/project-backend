package api

/*
* This file defines the different structs (and possibly interfaces)
* needed to implement the interfaces defined in api.go and its related files.
* Note that these compoenents follow the structures defined in the documentations.
 */

// NecessaryUserData contains the minimum data needed to open
// an account
type NecessaryUserData struct {
	Username string
	email    string
	password string // can be null
}

// PublicProfile struct contains user data that is publicly available
type PublicProfile struct {
	Username      string
	ID            string
	Avatar        any
	ProfilePic    []string // Placeholder type, modify as necessary
	Bio           string
	Prestige      int
	Followers     []string // Slice of usernames
	Following     []string
	ExternalLinks []string
	CreatedEvents []string // Slice of event IDs
}

// Profile struct contains both PublicProfile data and private data
// available only to current user
type Profile struct {
	PublicData      PublicProfile
	email           string
	password        string
	FavouriteCats   []string
	Blakclist       []string
	SafeArea        any
	ProfileUpgrades []string
	FollowedEvents  []string
	JoinedEvents    []string
}
