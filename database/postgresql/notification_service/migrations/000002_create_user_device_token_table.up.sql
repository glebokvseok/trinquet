CREATE TABLE IF NOT EXISTS notification_service.user_device_token(
    user_id UUID NOT NULL,
    device_token VARCHAR(256) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, device_token)
);
