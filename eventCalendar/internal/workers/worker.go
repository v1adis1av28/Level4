package workers

import (
	"context"
	"eventCalendar/internal/models"
	"eventCalendar/internal/storage"
	"time"

	"github.com/quay/zlog"
)

func StartNotificationJobWorker(ch <-chan models.NotificationJob) {
	for job := range ch {
		delay := time.Until(job.ReminderTime)
		if delay > 0 {
			go func(job models.NotificationJob) {
				time.Sleep(delay)
				sendNotification(job)
			}(job)
		} else {
			zlog.Warn(context.Background()).Msgf("Notification for event %d is overdue or due immediately. Sending now.", job.ID)
			go sendNotification(job)
		}
	}
}

func sendNotification(job models.NotificationJob) {
	zlog.Info(context.Background()).Msgf("Send notification on id: %v", job.ID)
}

func StartArchivalWorker(storage *storage.Storage, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ctx := context.Background()
	for {
		select {
		case <-ticker.C:
			cutoffTime := time.Now().AddDate(0, 0, -30)

			err := storage.ArchiveOldEvents(ctx, cutoffTime)
			if err != nil {
				zlog.Error(ctx).Msgf("Archival worker failed: %v", err)
			}
		}
	}
}
