package main

// User is used to encode coral user data.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Asset is used to encode coral asset data.
type Asset struct {
	ID    int    `json:"id"`
	Title string `json:"name"`
}

// Comment is used to encode coral comment data.
type Comment struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	AssetID  int    `json:"asset_id"`
	ParentID int    `json:"parent_id,omitempty"`
	Body     string `json:"body"`
}

// Job is used to encode data to be sent to sponge.
type Job struct {
	Data []byte
	Type string
}
