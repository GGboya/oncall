package model

import "time"

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	AvatarTemplate string `json:"avatar_template"`
	FlairName      string `json:"flair_name"`
	TrustLevel     int    `json:"trust_level"`
	AssignIcon     string `json:"assign_icon"`
	AssignPath     string `json:"assign_path"`
	IsAdmin        bool   `json:"admin,omitempty"`
	IsModerator    bool   `json:"moderator,omitempty"`
}

type Poster struct {
	Extras      string `json:"extras,omitempty"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

type Topic struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	FancyTitle         string    `json:"fancy_title"`
	Slug               string    `json:"slug"`
	PostsCount         int       `json:"posts_count"`
	ReplyCount         int       `json:"reply_count"`
	HighestPostNumber  int       `json:"highest_post_number"`
	ImageUrl           string    `json:"image_url"`
	CreatedAt          time.Time `json:"created_at"`
	LastPostedAt       time.Time `json:"last_posted_at"`
	Bumped             bool      `json:"bumped"`
	BumpedAt           string    `json:"bumped_at"`
	Archetype          string    `json:"archetype"`
	Unseen             bool      `json:"unseen"`
	LastReadPostNumber int       `json:"last_read_post_number"`
	Unread             int       `json:"unread"`
	NewPosts           int       `json:"new_posts"`
	UnreadPosts        int       `json:"unread_posts"`
	Pinned             bool      `json:"pinned"`
	Visible            bool      `json:"visible"`
	Closed             bool      `json:"closed"`
	Archived           bool      `json:"archived"`
	NotificationLevel  int       `json:"notification_level"`
	Bookmarked         bool      `json:"bookmarked"`
	LikeCount          int       `json:"like_count"`
	HasSummary         bool      `json:"has_summary"`
	LastPosterUsername string    `json:"last_poster_username"`
	CategoryID         int       `json:"category_id"`
	PinnedGlobally     bool      `json:"pinned_globally"`
	FeaturedLink       string    `json:"featured_link"`
	HasAcceptedAnswer  bool      `json:"has_accepted_answer"`
	CanVote            bool      `json:"can_vote"`
	Posters            []Poster  `json:"posters"`
	Tags               []string  `json:"tags"`
}

type Response struct {
	Users     []User `json:"users"`
	TopicList struct {
		CanCreateTopic bool     `json:"can_create_topic"`
		MoreTopicsUrl  string   `json:"more_topics_url"`
		PerPage        int      `json:"per_page"`
		TopTags        []string `json:"top_tags"`
		Topics         []Topic  `json:"topics"`
	} `json:"topic_list"`
}

// Post 包含帖子的详细信息，我们只关注 cooked 字段
type Post struct {
	Cooked string `json:"cooked"`
	Name   string `json:"name"`
}

// PostStream 包含帖子流的信息
type PostStream struct {
	Posts []Post `json:"posts"`
}

type PostResponse struct {
	PostStream PostStream `json:"post_stream"`
}
