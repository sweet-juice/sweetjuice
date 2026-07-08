package com.wailspackage.datadir;

import android.content.Context;
import android.content.Intent;
import com.wailsplugin.WailsPlugin;
import org.json.JSONException;
import org.json.JSONObject;
import java.io.File;

/**
 * DataDirPlugin provides access to standard Android application directories.
 */
public class DataDirPlugin implements WailsPlugin {
    private Context mContext;

    @Override
    public String getDomain() { return "datadir"; }

    @Override
    public void onAttach(Context context) {
        this.mContext = context.getApplicationContext();
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if (mContext == null) return "{\"error\":\"Context not attached\"}";

        try {
            switch (action) {
                case "getDirs":
                    return getDirs();
                default:
                    return "{\"error\":\"Unknown action\"}";
            }
        } catch (JSONException e) {
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    private String getDirs() throws JSONException {
        JSONObject dirs = new JSONObject();
        
        // Private internal storage
        dirs.put("files", mContext.getFilesDir().getAbsolutePath());
        dirs.put("cache", mContext.getCacheDir().getAbsolutePath());
        
        // External storage (scoped to app)
        File extFiles = mContext.getExternalFilesDir(null);
        if (extFiles != null) {
            dirs.put("external_files", extFiles.getAbsolutePath());
        }
        
        File extCache = mContext.getExternalCacheDir();
        if (extCache != null) {
            dirs.put("external_cache", extCache.getAbsolutePath());
        }

        return dirs.toString();
    }

    @Override public void onActivityResult(int req, int res, Intent d) {}
    @Override public void onRequestPermissionsResult(int req, String[] p, int[] res) {}
    @Override public void onNewIntent(Intent intent) {}
}
