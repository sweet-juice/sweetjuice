package com.wailspackage.daemon;

import android.content.Context;
import android.content.Intent;
import android.os.Build;
import com.wailsplugin.WailsPlugin;
import org.json.JSONException;
import org.json.JSONObject;

/**
 * DaemonPlugin provides access to Android Foreground Services.
 */
public class DaemonPlugin implements WailsPlugin {
    private Context mContext;

    @Override
    public String getDomain() { return "daemon"; }

    @Override
    public void onAttach(Context context) {
        this.mContext = context.getApplicationContext();
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if (mContext == null) return "{\"error\":\"Context not attached\"}";

        try {
            switch (action) {
                case "start":
                    return start(jsonArgsPayload);
                case "stop":
                    return stop();
                default:
                    return "{\"error\":\"Unknown action\"}";
            }
        } catch (JSONException e) {
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    private String start(String jsonArgs) throws JSONException {
        JSONObject args = new JSONObject(jsonArgs);
        Intent intent = new Intent(mContext, WailsDaemonService.class);
        intent.setAction(WailsDaemonService.ACTION_START);
        intent.putExtra("title", args.optString("title", "App running"));
        intent.putExtra("message", args.optString("message", "Keeping connection alive"));
        intent.putExtra("channel_id", args.optString("channel_id", "daemon_channel"));
        intent.putExtra("channel_name", args.optString("channel_name", "Daemon Service"));
        intent.putExtra("importance", args.optString("importance", "LOW"));

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            mContext.startForegroundService(intent);
        } else {
            mContext.startService(intent);
        }
        return "{\"status\":\"started\"}";
    }

    private String stop() {
        Intent intent = new Intent(mContext, WailsDaemonService.class);
        intent.setAction(WailsDaemonService.ACTION_STOP);
        mContext.startService(intent);
        return "{\"status\":\"stopped\"}";
    }

    @Override public void onActivityResult(int req, int res, Intent d) {}
    @Override public void onRequestPermissionsResult(int req, String[] p, int[] res) {}
    @Override public void onNewIntent(Intent intent) {}
}
