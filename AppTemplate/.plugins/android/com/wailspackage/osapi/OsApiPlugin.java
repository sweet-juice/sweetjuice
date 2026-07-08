package com.wailspackage.osapi;

import android.content.Context;
import android.content.Intent;
import android.os.Build;
import com.wailsplugin.WailsPlugin;
import org.json.JSONException;
import org.json.JSONObject;

/**
 * OsApiPlugin provides information about the Android system version.
 */
public class OsApiPlugin implements WailsPlugin {
    @Override
    public String getDomain() { return "osapi"; }

    @Override
    public void onAttach(Context context) {}

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        try {
            switch (action) {
                case "getInfo":
                    return getInfo();
                default:
                    return "{\"error\":\"Unknown action\"}";
            }
        } catch (JSONException e) {
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    private String getInfo() throws JSONException {
        JSONObject info = new JSONObject();
        info.put("sdk_int", Build.VERSION.SDK_INT);
        info.put("release", Build.VERSION.RELEASE);
        info.put("codename", Build.VERSION.CODENAME);
        info.put("manufacturer", Build.MANUFACTURER);
        info.put("model", Build.MODEL);
        return info.toString();
    }

    @Override public void onActivityResult(int req, int res, Intent d) {}
    @Override public void onRequestPermissionsResult(int req, String[] p, int[] res) {}
}
