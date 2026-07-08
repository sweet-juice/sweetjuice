package com.wailspackage.permissions;

import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;

import com.wailsplugin.WailsPlugin;

import org.json.JSONException;
import org.json.JSONObject;

import wailsmobile.Wailsmobile;

/**
 * PermissionsPlugin handles system runtime permissions.
 * It is background-aware for status checks but requires an active Activity for request prompts.
 */
public class PermissionsPlugin implements WailsPlugin {
    private Context mContext;
    private AppCompatActivity mActivity;
    private static final int PERMISSION_REQ_CODE = 9911;

    @Override
    public String getDomain() { return "permissions"; }

    @Override
    public void onAttach(Context context) { 
        this.mContext = context; 
        if (context instanceof AppCompatActivity) {
            this.mActivity = (AppCompatActivity) context;
        }
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if ("check".equals(action)) {
            String perm = parsePermissionFromJson(jsonArgsPayload);
            // Context is sufficient for checking permissions
            int result = ContextCompat.checkSelfPermission(mContext, perm);
            return result == PackageManager.PERMISSION_GRANTED ? "{\"status\":\"granted\"}" : "{\"status\":\"denied\"}";
        }

        if ("request".equals(action)) {
            // Cannot request UI prompts from the background
            if (mActivity == null) {
                return "{\"error\":\"No active UI to request permissions\"}";
            }
            String perm = parsePermissionFromJson(jsonArgsPayload);
            ActivityCompat.requestPermissions(mActivity, new String[]{perm}, PERMISSION_REQ_CODE);
            return "{\"status\":\"requested\"}";
        }

        return "{\"error\":\"Unknown action\"}";
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {
        if (requestCode == PERMISSION_REQ_CODE && grantResults.length > 0) {
            boolean granted = grantResults[0] == PackageManager.PERMISSION_GRANTED;
            String permission = permissions[0];

            JSONObject result = new JSONObject();
            try {
                result.put("permission", permission);
                result.put("granted", granted);
                String payload = "[" + result.toString() + "]";
                Wailsmobile.handleNativeAction("permissions:result", payload);
            } catch (JSONException e) {
                Log.e("PermissionsPlugin", "Error creating result JSON", e);
            }
        }
    }

    @Override public void onActivityResult(int r, int rc, Intent d) {}
    @Override public void onNewIntent(Intent intent) {}

    private String parsePermissionFromJson(String json) {
        try {
            JSONObject obj = new JSONObject(json);
            return obj.optString("permission", android.Manifest.permission.CAMERA);
        } catch (JSONException e) {
            return android.Manifest.permission.CAMERA;
        }
    }
}
