package com.juicepackage.workmanager;

import android.content.Context;
import android.content.Intent;
import androidx.work.Constraints;
import androidx.work.Data;
import androidx.work.NetworkType;
import androidx.work.OneTimeWorkRequest;
import androidx.work.PeriodicWorkRequest;
import androidx.work.WorkManager;

import com.juiceplugin.SweetJuicePlugin;

import org.json.JSONException;
import org.json.JSONObject;
import java.util.concurrent.TimeUnit;

/**
 * WorkManagerPlugin (WorkManager) allows Go to schedule background tasks.
 */
public class WorkManagerPlugin implements SweetJuicePlugin {
    private Context mContext;

    @Override
    public String getDomain() { return "workmanager"; }

    @Override
    public void onAttach(Context context) { 
        this.mContext = context; 
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        try {
            JSONObject args = new JSONObject(jsonArgsPayload);
            String taskKey = args.optString("task_key", "default_task");
            
            Constraints constraints = parseConstraints(args.optJSONObject("constraints"));

            if ("enqueueOneTime".equals(action)) {
                Data inputData = new Data.Builder().putString("task_key", taskKey).build();

                OneTimeWorkRequest request = new OneTimeWorkRequest.Builder(SweetJuiceBackgroundWorker.class)
                        .setInputData(inputData)
                        .setConstraints(constraints)
                        .build();

                WorkManager.getInstance(mContext).enqueue(request);
                return "{\"status\":\"enqueued\",\"id\":\"" + request.getId().toString() + "\"}";
            }

            if ("enqueuePeriodic".equals(action)) {
                long intervalMinutes = args.optLong("interval_minutes", 15);
                Data inputData = new Data.Builder().putString("task_key", taskKey).build();

                PeriodicWorkRequest request = new PeriodicWorkRequest.Builder(
                        SweetJuiceBackgroundWorker.class, intervalMinutes, TimeUnit.MINUTES)
                        .setInputData(inputData)
                        .setConstraints(constraints)
                        .build();

                WorkManager.getInstance(mContext).enqueue(request);
                return "{\"status\":\"periodic_enqueued\",\"id\":\"" + request.getId().toString() + "\"}";
            }

            if ("cancelAll".equals(action)) {
                WorkManager.getInstance(mContext).cancelAllWork();
                return "{\"status\":\"cancelled_all\"}";
            }

        } catch (JSONException e) {
            return "{\"error\":\"Invalid JSON payload\"}";
        }

        return "{\"error\":\"Unknown action\"}";
    }

    private Constraints parseConstraints(JSONObject json) {
        Constraints.Builder builder = new Constraints.Builder();
        if (json == null) return builder.build();

        String net = json.optString("network_type", "NOT_REQUIRED");
        switch (net) {
            case "CONNECTED": builder.setRequiredNetworkType(NetworkType.CONNECTED); break;
            case "UNMETERED": builder.setRequiredNetworkType(NetworkType.UNMETERED); break;
            case "NOT_ROAMING": builder.setRequiredNetworkType(NetworkType.NOT_ROAMING); break;
            default: builder.setRequiredNetworkType(NetworkType.NOT_REQUIRED);
        }

        builder.setRequiresCharging(json.optBoolean("requires_charging", false));
        builder.setRequiresDeviceIdle(json.optBoolean("requires_device_idle", false));
        builder.setRequiresBatteryNotLow(json.optBoolean("requires_battery_not_low", false));
        builder.setRequiresStorageNotLow(json.optBoolean("requires_storage_not_low", false));

        return builder.build();
    }

    @Override public void onRequestPermissionsResult(int rc, String[] p, int[] g) {}
    @Override public void onActivityResult(int r, int rc, Intent d) {}
}
