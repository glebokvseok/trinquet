package collections

const (
	PostsCollection             = "posts"
	NewsfeedPostsCollectionName = "newsfeed_posts"
	MediasCollectionName        = "medias"
	PostLikesCollectionName     = "post_likes"
	PostCommentsCollectionName  = "post_comments"
)

type NewsfeedPost struct {
	UserID    string `bson:"user_id"`
	PostID    string `bson:"post_id"`
	IsDeleted bool   `bson:"is_deleted"`
	Timestamp int64  `bson:"timestamp"`
}

type Post struct {
	ID        string   `bson:"_id"`
	AuthorID  string   `bson:"author_id"`
	Text      string   `bson:"text"`
	MediaIDs  []string `bson:"media_ids"`
	Timestamp int64    `bson:"timestamp"`
}

type Media struct {
	ID       string `bson:"_id"`
	MimeType string `bson:"mime_type"`
}

type PostLike struct {
	PostID    string `bson:"post_id"`
	UserID    string `bson:"user_id"`
	IsActive  bool   `bson:"is_active"`
	Timestamp int64  `bson:"timestamp"`
}

type PostComment struct {
	ID        string             `bson:"_id"`
	PostID    string             `bson:"post_id"`
	AuthorID  string             `bson:"author_id"`
	Text      string             `bson:"text"`
	IsDeleted bool               `bson:"is_deleted"`
	Timestamp int64              `bson:"timestamp"`
	Replies   []PostCommentReply `bson:"replies"`
}

type PostCommentReply struct {
	AuthorID  string `bson:"author_id"`
	Text      string `bson:"text"`
	IsDeleted bool   `bson:"is_deleted"`
	Timestamp int64  `bson:"timestamp"`
}

type AggregatedPost struct {
	ID        string           `bson:"id"`
	AuthorID  string           `bson:"author_id"`
	Text      string           `bson:"text"`
	Timestamp int64            `bson:"timestamp"`
	IsLiked   bool             `bson:"is_liked"`
	Medias    []map[string]any `bson:"medias"`
}
