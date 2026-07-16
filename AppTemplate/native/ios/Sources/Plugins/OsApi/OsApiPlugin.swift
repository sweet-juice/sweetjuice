import Foundation
import UIKit

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
        let info: [String: Any] = [
            "name": device.name,
            "systemName": device.systemName,
            "systemVersion": device.systemVersion,
            "model": device.model,
            "localizedModel": device.localizedModel,
            "identifierForVendor": device.identifierForVendor?.uuidString ?? "",
            "isPhysicalDevice": TARGET_OS_SIMULATOR == 0
        ]

        if let data = try? JSONSerialization.data(withJSONObject: info),
           let json = String(data: data, encoding: .utf8) {
            return json
        }
        return "{}"
    }
}
