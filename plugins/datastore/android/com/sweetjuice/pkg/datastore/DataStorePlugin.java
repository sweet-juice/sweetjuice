package com.sweetjuice.pkg.datastore;

import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;

import com.sweetjuice.plugin.SweetJuicePlugin;

import org.json.JSONException;
import org.json.JSONObject;

import java.util.Map;

/**
 * DataStorePlugin provides a simple key-value store using SharedPreferences.
 */
public class DataStorePlugin implements SweetJuicePlugin {
    private static final String PREFS_NAME = "WailsDataStore";
    private Context mContext;
    private SharedPreferences mPrefs;

    @Override
    public String getDomain() {
        return "datastore";
    }

    @Override
    public void onAttach(Context context) {
        this.mContext = context.getApplicationContext();
        this.mPrefs = mContext.getSharedPreferences(PREFS_NAME, Context.MODE_PRIVATE);
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if (mPrefs == null) return "{\"error\":\"Prefs not initialized\"}";

        try {
            JSONObject args = new JSONObject(jsonArgsPayload);
            switch (action) {
                case "set":
                    return set(args);
                case "get":
                    return get(args);
                case "delete":
                    return delete(args);
                case "clear":
                    return clear();
                case "getAll":
                    return getAll();
                default:
                    return "{\"error\":\"Unknown action\"}";
            }
        } catch (JSONException e) {
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    private String set(JSONObject args) throws JSONException {
        String key = args.getString("key");
        String value = args.getString("value");
        mPrefs.edit().putString(key, value).apply();
        return "{\"status\":\"ok\"}";
    }

    private String get(JSONObject args) throws JSONException {
        String key = args.getString("key");
        String defaultValue = args.optString("default", null);
        String value = mPrefs.getString(key, defaultValue);
        
        JSONObject resp = new JSONObject();
        resp.put("value", value);
        return resp.toString();
    }

    private String delete(JSONObject args) throws JSONException {
        String key = args.getString("key");
        mPrefs.edit().remove(key).apply();
        return "{\"status\":\"ok\"}";
    }

    private String clear() {
        mPrefs.edit().clear().apply();
        return "{\"status\":\"ok\"}";
    }

    private String getAll() throws JSONException {
        Map<String, ?> allEntries = mPrefs.getAll();
        JSONObject resp = new JSONObject();
        for (Map.Entry<String, ?> entry : allEntries.entrySet()) {
            resp.put(entry.getKey(), entry.getValue().toString());
        }
        return resp.toString();
    }

    @Override public void onActivityResult(int requestCode, int resultCode, Intent data) {}
    @Override public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {}
    @Override public void onNewIntent(Intent intent) {}
}
