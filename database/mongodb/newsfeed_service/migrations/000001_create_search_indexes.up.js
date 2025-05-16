db.newsfeed_posts.createIndex({ user_id: 1, timestamp: -1, is_deleted: 1 })

db.post_likes.createIndex({ post_id: 1, user_id: 1 })

db.post_comments.createIndex({ post_id: 1, timestamp: -1 })
