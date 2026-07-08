package com.wailspackage.devicestate;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.net.ConnectivityManager;
import android.net.Network;
import android.net.NetworkCapabilities;
import android.net.NetworkInfo;
import android.os.BatteryManager;
import android.os.Build;
import android.os.PowerManager;
import android.util.Log;

import androidx.annotation.NonNull;

import com.wailsplugin.WailsPlugin;

import org.json.JSONException;
import org.json.JSONObject;

import wailsmobile.Wailsmobile;

public class DeviceStatePlugin implements WailsPlugin {
    private static final String TAG = "DeviceStatePlugin";
    private Context mContext;
    private BroadcastReceiver mBatteryReceiver;
    private BroadcastReceiver mConnectivityReceiver;
    private ConnectivityManager.NetworkCallback mNetworkCallback;
    private boolean monitoring;

    @Override
    public String getDomain() {
        return "devicestate";
    }

    @Override
    public void onAttach(Context context) {
        this.mContext = context.getApplicationContext();
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if (mContext == null) {
            return "{\"error\":\"Context not attached\"}";
        }

        try {
            switch (action) {
                case "getState":
                    return buildStateJson();
                case "startMonitoring":
                    return startMonitoring();
                case "stopMonitoring":
                    return stopMonitoring();
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

    @Override
    public void onNewIntent(Intent intent) {
    }

    private String buildStateJson() throws JSONException {
        JSONObject state = new JSONObject();
        state.put("battery_level", getBatteryLevel());
        state.put("is_charging", isCharging());
        state.put("battery_status", getBatteryStatus());
        state.put("battery_health", getBatteryHealth());
        state.put("temperature", getBatteryTemperature());
        state.put("low_power_mode", isLowPowerMode());
        state.put("connectivity", getConnectivityInfo());
        state.put("timestamp", System.currentTimeMillis());
        return state.toString();
    }

    private int getBatteryLevel() {
        Intent intent = mContext.registerReceiver(null, new IntentFilter(Intent.ACTION_BATTERY_CHANGED));
        if (intent == null) {
            return -1;
        }
        int level = intent.getIntExtra(BatteryManager.EXTRA_LEVEL, -1);
        int scale = intent.getIntExtra(BatteryManager.EXTRA_SCALE, 100);
        if (scale <= 0) {
            scale = 100;
        }
        return (int) ((level / (float) scale) * 100);
    }

    private boolean isCharging() {
        Intent intent = mContext.registerReceiver(null, new IntentFilter(Intent.ACTION_BATTERY_CHANGED));
        if (intent == null) {
            return false;
        }
        int status = intent.getIntExtra(BatteryManager.EXTRA_STATUS, -1);
        return status == BatteryManager.BATTERY_STATUS_CHARGING || status == BatteryManager.BATTERY_STATUS_FULL;
    }

    private String getBatteryStatus() {
        Intent intent = mContext.registerReceiver(null, new IntentFilter(Intent.ACTION_BATTERY_CHANGED));
        if (intent == null) {
            return "unknown";
        }
        int status = intent.getIntExtra(BatteryManager.EXTRA_STATUS, -1);
        switch (status) {
            case BatteryManager.BATTERY_STATUS_CHARGING:
                return "charging";
            case BatteryManager.BATTERY_STATUS_DISCHARGING:
                return "discharging";
            case BatteryManager.BATTERY_STATUS_FULL:
                return "full";
            case BatteryManager.BATTERY_STATUS_NOT_CHARGING:
                return "not_charging";
            case BatteryManager.BATTERY_STATUS_UNKNOWN:
            default:
                return "unknown";
        }
    }

    private String getBatteryHealth() {
        Intent intent = mContext.registerReceiver(null, new IntentFilter(Intent.ACTION_BATTERY_CHANGED));
        if (intent == null) {
            return "unknown";
        }
        int health = intent.getIntExtra(BatteryManager.EXTRA_HEALTH, -1);
        switch (health) {
            case BatteryManager.BATTERY_HEALTH_COLD:
                return "cold";
            case BatteryManager.BATTERY_HEALTH_DEAD:
                return "dead";
            case BatteryManager.BATTERY_HEALTH_GOOD:
                return "good";
            case BatteryManager.BATTERY_HEALTH_OVERHEAT:
                return "overheat";
            case BatteryManager.BATTERY_HEALTH_OVER_VOLTAGE:
                return "over_voltage";
            case BatteryManager.BATTERY_HEALTH_UNSPECIFIED_FAILURE:
                return "failure";
            case BatteryManager.BATTERY_HEALTH_UNKNOWN:
            default:
                return "unknown";
        }
    }

    private double getBatteryTemperature() {
        Intent intent = mContext.registerReceiver(null, new IntentFilter(Intent.ACTION_BATTERY_CHANGED));
        if (intent == null) {
            return -1.0;
        }
        int temperature = intent.getIntExtra(BatteryManager.EXTRA_TEMPERATURE, -1);
        if (temperature < 0) {
            return -1.0;
        }
        return temperature / 10.0;
    }

    private boolean isLowPowerMode() {
        PowerManager pm = (PowerManager) mContext.getSystemService(Context.POWER_SERVICE);
        if (pm == null) {
            return false;
        }
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            return pm.isPowerSaveMode();
        }
        return false;
    }

    private JSONObject getConnectivityInfo() throws JSONException {
        JSONObject connectivity = new JSONObject();
        boolean isConnected = false;
        String networkType = "UNKNOWN";
        boolean isRoaming = false;
        boolean isUnmetered = false;

        ConnectivityManager cm = (ConnectivityManager) mContext.getSystemService(Context.CONNECTIVITY_SERVICE);
        if (cm != null) {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
                Network network = cm.getActiveNetwork();
                if (network != null) {
                    NetworkCapabilities caps = cm.getNetworkCapabilities(network);
                    if (caps != null) {
                        isConnected = caps.hasCapability(NetworkCapabilities.NET_CAPABILITY_INTERNET);
                        if (caps.hasTransport(NetworkCapabilities.TRANSPORT_WIFI)) {
                            networkType = "WIFI";
                        } else if (caps.hasTransport(NetworkCapabilities.TRANSPORT_CELLULAR)) {
                            networkType = "CELLULAR";
                        } else if (caps.hasTransport(NetworkCapabilities.TRANSPORT_VPN)) {
                            networkType = "VPN";
                        }
                        isUnmetered = caps.hasCapability(NetworkCapabilities.NET_CAPABILITY_NOT_METERED);
                        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.P) {
                            isRoaming = !caps.hasCapability(NetworkCapabilities.NET_CAPABILITY_NOT_ROAMING);
                        }
                    }
                }
            } else {
                NetworkInfo info = cm.getActiveNetworkInfo();
                if (info != null) {
                    isConnected = info.isConnected();
                    isRoaming = info.isRoaming();
                    int type = info.getType();
                    if (type == ConnectivityManager.TYPE_WIFI) {
                        networkType = "WIFI";
                    } else if (type == ConnectivityManager.TYPE_MOBILE) {
                        networkType = "CELLULAR";
                    } else {
                        networkType = info.getTypeName();
                    }
                }
            }
        }

        connectivity.put("is_connected", isConnected);
        connectivity.put("network_type", networkType);
        connectivity.put("is_roaming", isRoaming);
        connectivity.put("is_unmetered", isUnmetered);
        return connectivity;
    }

    private String startMonitoring() {
        if (monitoring) {
            return "{\"status\":\"already_monitoring\"}";
        }

        if (mContext == null) {
            return "{\"error\":\"Context not attached\"}";
        }

        mBatteryReceiver = new BroadcastReceiver() {
            @Override
            public void onReceive(Context context, Intent intent) {
                emitStateChanged();
            }
        };
        mContext.registerReceiver(mBatteryReceiver, new IntentFilter(Intent.ACTION_BATTERY_CHANGED));

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.N) {
            ConnectivityManager cm = (ConnectivityManager) mContext.getSystemService(Context.CONNECTIVITY_SERVICE);
            if (cm != null) {
                mNetworkCallback = new ConnectivityManager.NetworkCallback() {
                    @Override
                    public void onAvailable(@NonNull Network network) {
                        emitStateChanged();
                    }

                    @Override
                    public void onLost(@NonNull Network network) {
                        emitStateChanged();
                    }
                };
                cm.registerDefaultNetworkCallback(mNetworkCallback);
            }
        } else {
            mConnectivityReceiver = new BroadcastReceiver() {
                @Override
                public void onReceive(Context context, Intent intent) {
                    emitStateChanged();
                }
            };
            mContext.registerReceiver(mConnectivityReceiver, new IntentFilter(ConnectivityManager.CONNECTIVITY_ACTION));
        }

        monitoring = true;
        return "{\"status\":\"monitoring_started\"}";
    }

    private String stopMonitoring() {
        if (!monitoring) {
            return "{\"status\":\"not_monitoring\"}";
        }

        try {
            if (mBatteryReceiver != null) {
                mContext.unregisterReceiver(mBatteryReceiver);
                mBatteryReceiver = null;
            }
            if (mConnectivityReceiver != null) {
                mContext.unregisterReceiver(mConnectivityReceiver);
                mConnectivityReceiver = null;
            }
            if (mNetworkCallback != null) {
                ConnectivityManager cm = (ConnectivityManager) mContext.getSystemService(Context.CONNECTIVITY_SERVICE);
                if (cm != null) {
                    cm.unregisterNetworkCallback(mNetworkCallback);
                }
                mNetworkCallback = null;
            }
        } catch (IllegalArgumentException e) {
            Log.w(TAG, "Receiver was not registered", e);
        }

        monitoring = false;
        return "{\"status\":\"monitoring_stopped\"}";
    }

    private void emitStateChanged() {
        try {
            String stateJson = buildStateJson();
            String payload = "[" + stateJson + "]";
            Wailsmobile.handleNativeAction("devicestate:changed", payload);
        } catch (JSONException e) {
            Log.e(TAG, "Failed to emit device state change", e);
        }
    }
}
