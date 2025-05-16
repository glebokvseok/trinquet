db.newsfeed_posts.dropIndex({ user_id: 1, timestamp: -1, is_deleted: 1 })

db.post_likes.dropIndex({ post_id: 1, user_id: 1 })

db.post_comments.dropIndex({ post_id: 1, timestamp: -1 })
