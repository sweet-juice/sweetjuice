import Foundation
import UIKit
import Network
import Sweetjuice

@MainActor
public class DeviceStatePlugin: SweetJuicePlugin {
    private var container: UIViewController?
    private let monitor = NWPathMonitor()
    private var isMonitoring = false

    public init() {}

    public func getDomain() -> String {
        return "devicestate"
    }

    public func onAttach(container: UIViewController) {
        self.container = container
        UIDevice.current.isBatteryMonitoringEnabled = true
    }

    public func handleAction(action: String, jsonArgs: String) -> String {
        switch action {
        case "getState":
            return buildStateJson()
        case "startMonitoring":
            return startMonitoring()
        case "stopMonitoring":
            return stopMonitoring()
        default:
            return "{\"error\":\"Unknown action\"}"
        }
    }

    private func buildStateJson() -> String {
        let batteryLevel = Int(UIDevice.current.batteryLevel * 100)
        let batteryStatus = getBatteryStatus()
        let isLowPowerMode = ProcessInfo.processInfo.isLowPowerModeEnabled

        var state: [String: Any] = [
            "battery_level": batteryLevel,
            "is_charging": batteryStatus == "charging" || batteryStatus == "full",
            "battery_status": batteryStatus,
            "low_power_mode": isLowPowerMode,
            "timestamp": Int64(Date().timeIntervalSince1970 * 1000)
        ]

        // Connectivity
        let path = monitor.currentPath
        let connectivity: [String: Any] = [
            "is_connected": path.status == .satisfied,
            "network_type": getNetworkType(path: path),
            "is_roaming": false,
            "is_unmetered": !path.isExpensive
        ]
        state["connectivity"] = connectivity

        if let data = try? JSONSerialization.data(withJSONObject: state),
           let json = String(data: data, encoding: .utf8) {
            return json
        }
        return "{}"
    }

    private func getBatteryStatus() -> String {
        switch UIDevice.current.batteryState {
        case .charging: return "charging"
        case .full: return "full"
        case .unplugged: return "discharging"
        case .unknown: return "unknown"
        @unknown default: return "unknown"
        }
    }

    private func getNetworkType(path: NWPath) -> String {
        if path.usesInterfaceType(.wifi) { return "WIFI" }
        if path.usesInterfaceType(.cellular) { return "CELLULAR" }
        if path.usesInterfaceType(.wiredEthernet) { return "ETHERNET" }
        return "UNKNOWN"
    }

    private func startMonitoring() -> String {
        if isMonitoring { return "{\"status\":\"already_monitoring\"}" }

        monitor.pathUpdateHandler = { _ in
            Task { @MainActor in
                self.emitStateChanged()
            }
        }
        let queue = DispatchQueue(label: "DeviceStateMonitor")
        monitor.start(queue: queue)

        NotificationCenter.default.addObserver(
            forName: UIDevice.batteryLevelDidChangeNotification,
            object: nil, queue: .main) { _ in
                Task { @MainActor in
                    self.emitStateChanged()
                }
        }

        isMonitoring = true
        return "{\"status\":\"monitoring_started\"}"
    }

    private func stopMonitoring() -> String {
        if !isMonitoring { return "{\"status\":\"not_monitoring\"}" }
        monitor.cancel()
        NotificationCenter.default.removeObserver(self, name: UIDevice.batteryLevelDidChangeNotification, object: nil)
        isMonitoring = false
        return "{\"status\":\"monitoring_stopped\"}"
    }

    private func emitStateChanged() {
        let stateJson = buildStateJson()
        let payload = "[\(stateJson)]"
        SweetjuiceHandleNativeAction("devicestate:changed", payload)
    }
}
