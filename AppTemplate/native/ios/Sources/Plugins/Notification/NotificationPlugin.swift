import Foundation
import UIKit
import UserNotifications

public class NotificationPlugin: SweetJuicePlugin {
    private var container: UIViewController?

    public init() {}

    public func getDomain() -> String {
        return "notification"
    }

    public func onAttach(container: UIViewController) {
        self.container = container
    }

    public func handleAction(action: String, jsonArgs: String) -> String {
        guard let data = jsonArgs.data(using: .utf8),
              let args = try? JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            return "{\"error\":\"Invalid JSON payload\"}"
        }

        if action == "post" {
            let id = args["id"] as? Int ?? Int(Date().timeIntervalSince1970)
            let title = args["title"] as? String ?? ""
            let body = args["body"] as? String ?? ""

            postNotification(id: id, title: title, body: body)
            return "{\"status\":\"posted\", \"id\":\(id)}"
        }

        if action == "cancel" {
            let id = args["id"] as? Int ?? -1
            UNUserNotificationCenter.current().removePendingNotificationRequests(withIdentifiers: ["\(id)"])
            return "{\"status\":\"cancelled\"}"
        }

        return "{\"error\":\"Unknown action\"}"
    }

    private func postNotification(id: Int, title: String, body: String) {
        let content = UNMutableNotificationContent()
        content.title = title
        content.body = body
        content.sound = .default

        let trigger = UNTimeIntervalNotificationTrigger(timeInterval: 1, repeats: false)
        let request = UNNotificationRequest(identifier: "\(id)", content: content, trigger: trigger)

        UNUserNotificationCenter.current().add(request) { error in
            if let error = error {
                print("Error posting notification: \(error.localizedDescription)")
            }
        }
    }
}
