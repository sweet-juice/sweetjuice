package com.sweetjuice.pkg.filepicker;

import android.app.Activity;
import android.content.Context;
import android.content.Intent;
import android.net.Uri;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;

import com.sweetjuice.plugin.SweetJuicePlugin;

import org.json.JSONException;
import org.json.JSONObject;

import sweetjuice.Sweetjuice;

/**
 * FilePickerPlugin allows selecting files or media from the device.
 */
public class FilePickerPlugin implements SweetJuicePlugin {
    private static final String TAG = "FilePickerPlugin";
    private static final int PICK_FILE_REQUEST = 4422;
    private Context mContext;
    private AppCompatActivity mActivity;

    @Override
    public String getDomain() {
        return "filepicker";
    }

    @Override
    public void onAttach(Context context) {
        this.mContext = context.getApplicationContext();
        if (context instanceof AppCompatActivity) {
            this.mActivity = (AppCompatActivity) context;
        }
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if (mActivity == null) return "{\"error\":\"No active UI to show file picker\"}";

        try {
            JSONObject args = new JSONObject(jsonArgsPayload);
            switch (action) {
                case "pickFile":
                    return pickFile(args);
                default:
                    return "{\"error\":\"Unknown action\"}";
            }
        } catch (JSONException e) {
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    private String pickFile(JSONObject args) throws JSONException {
        String mimeType = args.optString("mime_type", "*/*");
        boolean multiple = args.optBoolean("multiple", false);

        Intent intent = new Intent(Intent.ACTION_GET_CONTENT);
        intent.setType(mimeType);
        intent.putExtra(Intent.EXTRA_ALLOW_MULTIPLE, multiple);
        intent.addCategory(Intent.CATEGORY_OPENABLE);

        mActivity.startActivityForResult(Intent.createChooser(intent, "Select File"), PICK_FILE_REQUEST);
        return "{\"status\":\"started\"}";
    }

    @Override
    public void onActivityResult(int requestCode, int resultCode, Intent data) {
        if (requestCode == PICK_FILE_REQUEST) {
            if (resultCode == Activity.RESULT_OK && data != null) {
                try {
                    JSONObject result = new JSONObject();
                    if (data.getClipData() != null) {
                        int count = data.getClipData().getItemCount();
                        org.json.JSONArray uris = new org.json.JSONArray();
                        for (int i = 0; i < count; i++) {
                            uris.put(data.getClipData().getItemAt(i).getUri().toString());
                        }
                        result.put("uris", uris);
                        result.put("multiple", true);
                    } else if (data.getData() != null) {
                        result.put("uri", data.getData().toString());
                        result.put("multiple", false);
                    }
                    
                    String payload = "[" + result.toString() + "]";
                    Sweetjuice.handleNativeAction("filepicker:result", payload);
                } catch (JSONException e) {
                    Log.e(TAG, "Failed to build result JSON", e);
                }
            } else {
                Sweetjuice.handleNativeAction("filepicker:result", "[{\"error\":\"cancelled\"}]");
            }
        }
    }

    @Override public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {}
    @Override public void onNewIntent(Intent intent) {}
}
