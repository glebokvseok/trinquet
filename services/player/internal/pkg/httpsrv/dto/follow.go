package dto

type FollowInfo struct {
	FollowersCount int64 `json:"followers_count"`
	FollowingCount int64 `json:"following_count"`
	Following      bool  `json:"following"`
	FollowingBack  bool  `json:"following_back"`
}
