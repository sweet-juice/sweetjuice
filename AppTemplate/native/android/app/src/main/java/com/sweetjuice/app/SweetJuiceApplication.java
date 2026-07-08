package com.sweetjuice.app;

import android.app.Application;
import com.sweetjuice.plugin.SweetJuicePlugin;
import com.sweetjuice.pkg.biometric.BiometricPlugin;
import com.sweetjuice.pkg.daemon.DaemonPlugin;
import com.sweetjuice.pkg.devicestate.DeviceStatePlugin;
import com.sweetjuice.pkg.filepicker.FilePickerPlugin;
import com.sweetjuice.pkg.logger.LoggerPlugin;
import com.sweetjuice.pkg.notifications.NotificationPlugin;
import com.sweetjuice.pkg.permissions.PermissionsPlugin;
import com.sweetjuice.pkg.workmanager.WorkManagerPlugin;
import java.util.HashMap;
import java.util.Map;
import sweetjuice.Sweetjuice;

public class SweetJuiceApplication extends Application {
    private final Map<String, SweetJuicePlugin> mPlugins = new HashMap<>();

    @Override
    public void onCreate() {
        super.onCreate();
        
        // Start Go engine
        Sweetjuice.startApplication();

        // Initialize and Register global plugins
        registerPlugin(new PermissionsPlugin());
        registerPlugin(new WorkManagerPlugin());
        registerPlugin(new NotificationPlugin());
        registerPlugin(new LoggerPlugin());
        registerPlugin(new DeviceStatePlugin());
        registerPlugin(new BiometricPlugin());
        registerPlugin(new FilePickerPlugin());
        registerPlugin(new DaemonPlugin());

        // Register the global handler for Go-to-Native calls
        Sweetjuice.setNativeCallHandler(new sweetjuice.NativeCallHandler() {
            @Override
            public String onNativeCall(String method, String args) {
                if (method.contains(":")) {
                    String[] parts = method.split(":", 2);
                    String domain = parts[0];
                    String action = parts[1];

                    SweetJuicePlugin plugin = mPlugins.get(domain);
                    if (plugin != null) {
                        return plugin.handleAction(action, args);
                    }
                }
                return "{\"error\":\"Plugin domain not found\"}";
            }
        });
    }

    private void registerPlugin(SweetJuicePlugin plugin) {
        // Prime the plugin with the Application context for background tasks
        plugin.onAttach(this);
        mPlugins.put(plugin.getDomain(), plugin);
    }

    public Map<String, SweetJuicePlugin> getPlugins() {
        return mPlugins;
    }
}
