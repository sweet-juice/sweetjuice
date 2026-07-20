import SwiftUI

@main
struct GenericAppApp: App {
    init() {
        // Register default SweetJuice plugins
        let manager = SweetJuiceManager.shared
        manager.registerPlugin(NotificationPlugin())
        manager.registerPlugin(PermissionsPlugin())
        manager.registerPlugin(DeviceStatePlugin())
        manager.registerPlugin(OsApiPlugin())
        manager.registerPlugin(LoggerPlugin())
        manager.registerPlugin(BiometricPlugin())
        manager.registerPlugin(FilePickerPlugin())
        manager.registerPlugin(DataDirPlugin())
    }

    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
