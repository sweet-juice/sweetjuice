package com.wailspackage.daemon;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.os.Build;
import android.os.IBinder;
import androidx.annotation.Nullable;
import androidx.core.app.NotificationCompat;

/**
 * WailsDaemonService keeps the Go process alive and high-priority.
 */
public class WailsDaemonService extends Service {
    public static final String ACTION_START = "START";
    public static final String ACTION_STOP = "STOP";
    
    private static final int NOTIFICATION_ID = 9988;

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        if (intent != null) {
            String action = intent.getAction();
            if (ACTION_START.equals(action)) {
                startServiceWithNotification(intent);
            } else if (ACTION_STOP.equals(action)) {
                stopForeground(true);
                stopSelf();
            }
        }
        return START_STICKY;
    }

    private void startServiceWithNotification(Intent intent) {
        String title = intent.getStringExtra("title");
        String message = intent.getStringExtra("message");
        String channelId = intent.getStringExtra("channel_id");
        String channelName = intent.getStringExtra("channel_name");
        String importanceStr = intent.getStringExtra("importance");

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            int importance = NotificationManager.IMPORTANCE_LOW;
            if ("HIGH".equals(importanceStr)) importance = NotificationManager.IMPORTANCE_HIGH;
            else if ("DEFAULT".equals(importanceStr)) importance = NotificationManager.IMPORTANCE_DEFAULT;

            NotificationChannel channel = new NotificationChannel(channelId, channelName, importance);
            NotificationManager manager = getSystemService(NotificationManager.class);
            if (manager != null) {
                manager.createNotificationChannel(channel);
            }
        }

        Notification notification = new NotificationCompat.Builder(this, channelId)
                .setContentTitle(title)
                .setContentText(message)
                .setSmallIcon(getApplicationInfo().icon)
                .setPriority(NotificationCompat.PRIORITY_LOW)
                .build();

        startForeground(NOTIFICATION_ID, notification);
    }

    @Nullable
    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }
}
