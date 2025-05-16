CREATE TABLE IF NOT EXISTS notification_service.user_notification(
    user_id UUID NOT NULL,
    notification_id UUID NOT NULL,
    notification_type INT NOT NULL,
    data JSONB NOT NULL,
    timestamp BIGINT NOT NULL,
    PRIMARY KEY (user_id, notification_id)
);
