import Foundation
import UIKit
import AVFoundation
import UserNotifications
import Sweetjuice

public class PermissionsPlugin: SweetJuicePlugin {
    private var container: UIViewController?

    public init() {}

    public func getDomain() -> String {
        return "permissions"
    }

    public func onAttach(container: UIViewController) {
        self.container = container
    }

    public func handleAction(action: String, jsonArgs: String) -> String {
        guard let data = jsonArgs.data(using: .utf8),
              let args = try? JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            return "{\"error\":\"Invalid JSON payload\"}"
        }

        let permission = args["permission"] as? String ?? ""

        if action == "check" {
            let status = checkPermission(permission)
            return "{\"status\":\"\(status)\"}"
        }

        if action == "request" {
            requestPermission(permission)
            return "{\"status\":\"requested\"}"
        }

        return "{\"error\":\"Unknown action\"}"
    }

    private func checkPermission(_ permission: String) -> String {
        switch permission {
        case "camera":
            let status = AVCaptureDevice.authorizationStatus(for: .video)
            return status == .authorized ? "granted" : "denied"
        case "notifications":
            var status = "unknown"
            let semaphore = DispatchSemaphore(value: 0)
            UNUserNotificationCenter.current().getNotificationSettings { settings in
                status = settings.authorizationStatus == .authorized ? "granted" : "denied"
                semaphore.signal()
            }
            _ = semaphore.wait(timeout: .now() + 2)
            return status
        default:
            return "denied"
        }
    }

    private func requestPermission(_ permission: String) {
        switch permission {
        case "camera":
            AVCaptureDevice.requestAccess(for: .video) { granted in
                self.emitResult(permission: "camera", granted: granted)
            }
        case "notifications":
            UNUserNotificationCenter.current().requestAuthorization(options: [.alert, .sound, .badge]) { granted, _ in
                self.emitResult(permission: "notifications", granted: granted)
            }
        default:
            break
        }
    }

    private func emitResult(permission: String, granted: Bool) {
        let payload = "[{\"permission\":\"\(permission)\", \"granted\":\(granted)}]"
        // This helper is available because we import Sweetjuice
        SweetjuiceHandleNativeAction("permissions:result", payload)
    }
}
