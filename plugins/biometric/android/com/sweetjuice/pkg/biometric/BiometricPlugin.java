package com.sweetjuice.pkg.biometric;

import android.content.Context;
import android.content.Intent;
import android.os.Build;
import android.util.Log;

import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;
import androidx.biometric.BiometricManager;
import androidx.biometric.BiometricPrompt;
import androidx.core.content.ContextCompat;

import com.sweetjuice.plugin.SweetJuicePlugin;

import org.json.JSONException;
import org.json.JSONObject;

import java.util.concurrent.Executor;

import sweetjuice.Sweetjuice;

/**
 * BiometricPlugin provides access to biometric authentication (Fingerprint, Face, etc.)
 */
public class BiometricPlugin implements SweetJuicePlugin {
    private static final String TAG = "BiometricPlugin";
    private Context mContext;
    private AppCompatActivity mActivity;

    @Override
    public String getDomain() {
        return "biometric";
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
        if (mContext == null) {
            return "{\"error\":\"Context not attached\"}";
        }

        try {
            switch (action) {
                case "canAuthenticate":
                    return canAuthenticate();
                case "authenticate":
                    return authenticate(jsonArgsPayload);
                default:
                    return "{\"error\":\"Unknown action\"}";
            }
        } catch (JSONException e) {
            return "{\"error\":\"" + e.getMessage() + "\"}";
        }
    }

    @Override
    public void onActivityResult(int requestCode, int resultCode, Intent data) {
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {
    }

    private String canAuthenticate() throws JSONException {
        BiometricManager biometricManager = BiometricManager.from(mContext);
        int result = biometricManager.canAuthenticate(BiometricManager.Authenticators.BIOMETRIC_STRONG | BiometricManager.Authenticators.DEVICE_CREDENTIAL);
        
        JSONObject response = new JSONObject();
        boolean can = false;
        String status = "UNKNOWN";

        switch (result) {
            case BiometricManager.BIOMETRIC_SUCCESS:
                can = true;
                status = "SUCCESS";
                break;
            case BiometricManager.BIOMETRIC_ERROR_NO_HARDWARE:
                status = "NO_HARDWARE";
                break;
            case BiometricManager.BIOMETRIC_ERROR_HW_UNAVAILABLE:
                status = "HW_UNAVAILABLE";
                break;
            case BiometricManager.BIOMETRIC_ERROR_NONE_ENROLLED:
                status = "NONE_ENROLLED";
                break;
            case BiometricManager.BIOMETRIC_ERROR_SECURITY_UPDATE_REQUIRED:
                status = "SECURITY_UPDATE_REQUIRED";
                break;
            case BiometricManager.BIOMETRIC_ERROR_UNSUPPORTED:
                status = "UNSUPPORTED";
                break;
            case BiometricManager.BIOMETRIC_STATUS_UNKNOWN:
                status = "UNKNOWN";
                break;
        }

        response.put("can_authenticate", can);
        response.put("status", status);
        return response.toString();
    }

    private String authenticate(String jsonArgs) throws JSONException {
        if (mActivity == null) {
            return "{\"error\":\"No active UI to show biometric prompt\"}";
        }

        JSONObject args = new JSONObject(jsonArgs);
        String title = args.optString("title", "Biometric Authentication");
        String subtitle = args.optString("subtitle", "");
        String description = args.optString("description", "");
        String negativeButtonText = args.optString("negative_button_text", "Cancel");

        mActivity.runOnUiThread(() -> {
            Executor executor = ContextCompat.getMainExecutor(mActivity);
            BiometricPrompt biometricPrompt = new BiometricPrompt(mActivity, executor, new BiometricPrompt.AuthenticationCallback() {
                @Override
                public void onAuthenticationError(int errorCode, @NonNull CharSequence errString) {
                    super.onAuthenticationError(errorCode, errString);
                    sendAuthResult(false, "ERROR", errString.toString(), errorCode);
                }

                @Override
                public void onAuthenticationSucceeded(@NonNull BiometricPrompt.AuthenticationResult result) {
                    super.onAuthenticationSucceeded(result);
                    sendAuthResult(true, "SUCCESS", null, 0);
                }

                @Override
                public void onAuthenticationFailed() {
                    super.onAuthenticationFailed();
                    sendAuthResult(false, "FAILED", "Authentication failed", 0);
                }
            });

            BiometricPrompt.PromptInfo promptInfo = new BiometricPrompt.PromptInfo.Builder()
                    .setTitle(title)
                    .setSubtitle(subtitle)
                    .setDescription(description)
                    .setNegativeButtonText(negativeButtonText)
                    .setAllowedAuthenticators(BiometricManager.Authenticators.BIOMETRIC_STRONG)
                    .build();

            biometricPrompt.authenticate(promptInfo);
        });

        return "{\"status\":\"started\"}";
    }

    private void sendAuthResult(boolean success, String status, String error, int errorCode) {
        try {
            JSONObject result = new JSONObject();
            result.put("success", success);
            result.put("status", status);
            if (error != null) {
                result.put("error", error);
                result.put("error_code", errorCode);
            }
            String payload = "[" + result.toString() + "]";
            Sweetjuice.handleNativeAction("biometric:result", payload);
        } catch (JSONException e) {
            Log.e(TAG, "Failed to send biometric result", e);
        }
    }
}
