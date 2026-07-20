import Foundation
import UIKit

@MainActor
public class OsApiPlugin: SweetJuicePlugin {
    public init() {}

    public func getDomain() -> String {
        return "osapi"
    }

    public func onAttach(container: UIViewController) {}

    public func handleAction(action: String, jsonArgs: String) -> String {
        switch action {
        case "getInfo":
            return getInfo()
        default:
            return "{\"error\":\"Unknown action\"}"
        }
    }

    private func getInfo() -> String {
        let device = UIDevice.current

        #if targetEnvironment(simulator)
        let isPhysicalDevice = false
        #else
        let isPhysicalDevice = true
        #endif

        let info: [String: Any] = [
            "name": device.name,
            "system_name": device.systemName,
            "system_version": device.systemVersion,
            "model": device.model,
            "localized_model": device.localizedModel,
            "identifier_for_vendor": device.identifierForVendor?.uuidString ?? "",
            "is_physical_device": isPhysicalDevice
        ]

        if let data = try? JSONSerialization.data(withJSONObject: info),
           let json = String(data: data, encoding: .utf8) {
            return json
        }
        return "{}"
    }
}
