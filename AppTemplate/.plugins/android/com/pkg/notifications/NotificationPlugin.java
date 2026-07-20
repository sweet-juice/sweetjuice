package com.juicepackage.notifications;

import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.content.Context;
import android.content.Intent;
import android.os.Build;
import androidx.core.app.NotificationCompat;

import com.juiceplugin.SweetJuicePlugin;

import org.json.JSONException;
import org.json.JSONObject;

/**
 * NotificationPlugin allows Go to post system notifications.
 */
public class NotificationPlugin implements SweetJuicePlugin {
    private Context mContext;

    @Override
    public String getDomain() { return "notification"; }

    @Override
    public void onAttach(Context context) { 
        this.mContext = context; 
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        try {
            JSONObject args = new JSONObject(jsonArgsPayload);

            if ("post".equals(action)) {
                int id = args.optInt("id", (int) System.currentTimeMillis());
                String title = args.optString("title", "");
                String body = args.optString("body", "");
                String channelId = args.optString("channel_id", "default");
                String channelName = args.optString("channel_name", "Default");
                String importance = args.optString("importance", "DEFAULT");

                postNotification(id, title, body, channelId, channelName, importance);
                return "{\"status\":\"posted\", \"id\":" + id + "}";
            }

            if ("cancel".equals(action)) {
                int id = args.optInt("id", -1);
                NotificationManager manager = (NotificationManager) mContext.getSystemService(Context.NOTIFICATION_SERVICE);
                manager.cancel(id);
                return "{\"status\":\"cancelled\"}";
            }

        } catch (JSONException e) {
            return "{\"error\":\"Invalid JSON payload\"}";
        }

        return "{\"error\":\"Unknown action\"}";
    }

    private void postNotification(int id, String title, String body, String channelId, String channelName, String importanceStr) {
        NotificationManager manager = (NotificationManager) mContext.getSystemService(Context.NOTIFICATION_SERVICE);

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            int importance = NotificationManager.IMPORTANCE_DEFAULT;
            switch (importanceStr) {
                case "HIGH": importance = NotificationManager.IMPORTANCE_HIGH; break;
                case "LOW": importance = NotificationManager.IMPORTANCE_LOW; break;
                case "MIN": importance = NotificationManager.IMPORTANCE_MIN; break;
            }
            NotificationChannel channel = new NotificationChannel(channelId, channelName, importance);
            manager.createNotificationChannel(channel);
        }

        NotificationCompat.Builder builder = new NotificationCompat.Builder(mContext, channelId)
                .setSmallIcon(mContext.getApplicationInfo().icon)
                .setContentTitle(title)
                .setContentText(body)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT)
                .setAutoCancel(true);

        manager.notify(id, builder.build());
    }

    @Override public void onRequestPermissionsResult(int rc, String[] p, int[] g) {}
    @Override public void onActivityResult(int r, int rc, Intent d) {}
}
