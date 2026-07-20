import Foundation
import UIKit

@MainActor
public class LoggerPlugin: SweetJuicePlugin {
    public init() {}

    public func getDomain() -> String {
        return "logger"
    }

    public func onAttach(container: UIViewController) {}

    public func handleAction(action: String, jsonArgs: String) -> String {
        guard let data = jsonArgs.data(using: .utf8),
              let args = try? JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            return "{\"error\":\"Invalid JSON payload\"}"
        }

        let message = args["message"] as? String ?? ""
        let tag = args["tag"] as? String ?? "SweetJuice"

        switch action {
        case "debug":
            print("[\(tag)] DEBUG: \(message)")
        case "info":
            print("[\(tag)] INFO: \(message)")
        case "error":
            print("[\(tag)] ERROR: \(message)")
        default:
            return "{\"error\":\"Unknown action\"}"
        }

        return "{\"status\":\"ok\"}"
    }
}
