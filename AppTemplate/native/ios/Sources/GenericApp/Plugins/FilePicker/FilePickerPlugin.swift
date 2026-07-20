import Foundation
import UIKit
import UniformTypeIdentifiers
import Sweetjuice

@MainActor
public class FilePickerPlugin: NSObject, SweetJuicePlugin, UIDocumentPickerDelegate {
    private var container: UIViewController?

    public override init() {
        super.init()
    }

    public func getDomain() -> String {
        return "filepicker"
    }

    public func onAttach(container: UIViewController) {
        self.container = container
    }

    public func handleAction(action: String, jsonArgs: String) -> String {
        guard let data = jsonArgs.data(using: .utf8),
              let args = try? JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            return "{\"error\":\"Invalid JSON payload\"}"
        }

        if action == "pickFile" {
            let multiple = args["multiple"] as? Bool ?? false

            DispatchQueue.main.async {
                let picker = UIDocumentPickerViewController(forOpeningContentTypes: [.data, .content])
                picker.delegate = self
                picker.allowsMultipleSelection = multiple
                self.container?.present(picker, animated: true)
            }
            return "{\"status\":\"started\"}"
        }

        return "{\"error\":\"Unknown action\"}"
    }

    public nonisolated func documentPicker(_ controller: UIDocumentPickerViewController, didPickDocumentsAt urls: [URL]) {
        let uris = urls.map { $0.absoluteString }
        let result: [String: Any] = [
            "uris": uris,
            "multiple": uris.count > 1
        ]

        if let data = try? JSONSerialization.data(withJSONObject: result),
           let json = String(data: data, encoding: .utf8) {
            let payload = "[\(json)]"
            SweetjuiceHandleNativeAction("filepicker:result", payload)
        }
    }

    public nonisolated func documentPickerWasCancelled(_ controller: UIDocumentPickerViewController) {
        SweetjuiceHandleNativeAction("filepicker:result", "[{\"error\":\"cancelled\"}]")
    }
}
