import Foundation
import UIKit
import LocalAuthentication
import Sweetjuice

public class BiometricPlugin: SweetJuicePlugin {
    private var container: UIViewController?

    public init() {}

    public func getDomain() -> String {
        return "biometric"
    }

    public func onAttach(container: UIViewController) {
        self.container = container
    }

    public func handleAction(action: String, jsonArgs: String) -> String {
        switch action {
        case "canAuthenticate":
            return canAuthenticate()
        case "authenticate":
            return authenticate(jsonArgs: jsonArgs)
        default:
            return "{\"error\":\"Unknown action\"}"
        }
    }

    private func canAuthenticate() -> String {
        let context = LAContext()
        var error: NSError?

        let can = context.canEvaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, error: &error)
        var status = "SUCCESS"

        if let error = error {
            switch LAError(_nsError: error).code {
            case .biometryNotAvailable: status = "NO_HARDWARE"
            case .biometryNotEnrolled: status = "NONE_ENROLLED"
            case .biometryLockout: status = "LOCKED_OUT"
            default: status = "UNSUPPORTED"
            }
        }

        return "{\"can_authenticate\":\(can), \"status\":\"\(status)\"}"
    }

    private func authenticate(jsonArgs: String) -> String {
        guard let data = jsonArgs.data(using: .utf8),
              let args = try? JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            return "{\"error\":\"Invalid JSON payload\"}"
        }

        let reason = args["description"] as? String ?? "Authenticate to continue"
        let context = LAContext()

        context.evaluatePolicy(.deviceOwnerAuthenticationWithBiometrics, localizedReason: reason) { success, error in
            var status = "SUCCESS"
            var errMsg: String?

            if !success {
                status = "FAILED"
                if let error = error {
                    errMsg = error.localizedDescription
                }
            }

            self.sendResult(success: success, status: status, error: errMsg)
        }

        return "{\"status\":\"started\"}"
    }

    private func sendResult(success: Bool, status: String, error: String?) {
        var result: [String: Any] = ["success": success, "status": status]
        if let error = error {
            result["error"] = error
        }

        if let data = try? JSONSerialization.data(withJSONObject: result),
           let json = String(data: data, encoding: .utf8) {
            let payload = "[\(json)]"
            SweetjuiceHandleNativeAction("biometric:result", payload)
        }
    }
}
