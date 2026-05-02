-- Create a new notification
-- name: CreateNotification :one

INSERT INTO notifications (
	user_id,
	message
) VALUES (
	$1,
	$2
)
RETURNING id, user_id, message, notification_type, created_at, updated_at;

-- get a notification by id
-- name: GetNotificationByID :one

SELECT id, user_id, message, notification_type, created_at, updated_at
FROM notifications
WHERE id = $1;

-- get all notifications
-- name: GetAllNotifications :many

SELECT id, user_id, message, notification_type, created_at, updated_at
FROM notifications
ORDER BY created_at DESC;

-- get all notifications by user_id
-- name: GetAllNotificationsByUserID :many

SELECT id, user_id, message, notification_type, created_at, updated_at
FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC;

-- update message/type
-- name: UpdateNotification :one

UPDATE notifications
SET
	message = $2,
	updated_at = now()
WHERE id = $1
RETURNING id, user_id, message, notification_type, created_at, updated_at;

-- delete by id
-- name: DeleteNotification :exec

DELETE FROM notifications
WHERE id = $1;