package com.juicepackage.logger;

import android.content.Context;
import android.content.Intent;
import android.util.Log;
import com.juiceplugin.SweetJuicePlugin;
import org.json.JSONException;
import org.json.JSONObject;

/**
 * LoggerPlugin routes logs from Go to Android's Logcat.
 */
public class LoggerPlugin implements SweetJuicePlugin {
    @Override
    public String getDomain() { return "logger"; }

    @Override
    public void onAttach(Context context) {}

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if ("log".equals(action)) {
            try {
                JSONObject args = new JSONObject(jsonArgsPayload);
                String tag = args.optString("tag", "SweetJuiceGo");
                String level = args.optString("level", "INFO");
                String message = args.optString("message", "");

                switch (level) {
                    case "DEBUG":
                        Log.d(tag, message);
                        break;
                    case "WARN":
                        Log.w(tag, message);
                        break;
                    case "ERROR":
                        Log.e(tag, message);
                        break;
                    case "INFO":
                    default:
                        Log.i(tag, message);
                        break;
                }
                return "{\"status\":\"ok\"}";
            } catch (JSONException e) {
                return "{\"error\":\"Invalid JSON\"}";
            }
        }
        return "{\"error\":\"Unknown action\"}";
    }

    @Override public void onActivityResult(int req, int res, Intent d) {}
    @Override public void onRequestPermissionsResult(int req, String[] p, int[] res) {}
}
