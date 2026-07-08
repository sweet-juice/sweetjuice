package com.wailspackage.deeplinking;

import android.content.Context;
import android.content.Intent;
import android.net.Uri;
import android.util.Log;

import com.wailsplugin.WailsPlugin;

import wailsmobile.Wailsmobile;

/**
 * DeepLinkingPlugin handles incoming deep links and forwards them to the Go layer.
 */
public class DeepLinkingPlugin implements WailsPlugin {
    private static final String TAG = "DeepLinkingPlugin";
    private Context mContext;

    @Override
    public String getDomain() {
        return "deeplinking";
    }

    @Override
    public void onAttach(Context context) {
        this.mContext = context.getApplicationContext();
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        return "{\"error\":\"Unknown action\"}";
    }

    @Override
    public void onActivityResult(int requestCode, int resultCode, Intent data) {
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {
    }

    @Override
    public void onNewIntent(Intent intent) {
        handleIntent(intent);
    }

    private void handleIntent(Intent intent) {
        if (intent == null) return;
        
        String action = intent.getAction();
        Uri data = intent.getData();

        if (Intent.ACTION_VIEW.equals(action) && data != null) {
            String url = data.toString();
            Log.d(TAG, "Received deep link: " + url);
            
            // Forward the URL to the Go layer
            String payload = "[\"" + url + "\"]";
            Wailsmobile.handleNativeAction("deeplinking:received", payload);
        }
    }
}
