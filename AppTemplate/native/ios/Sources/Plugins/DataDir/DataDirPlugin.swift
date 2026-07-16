import Foundation
import UIKit

public class DataDirPlugin: SweetJuicePlugin {
    public init() {}

    public func getDomain() -> String {
        return "datadir"
    }

    public func onAttach(container: UIViewController) {}

    public func handleAction(action: String, jsonArgs: String) -> String {
        if action == "getDirs" {
            return getDirs()
        }
        return "{\"error\":\"Unknown action\"}"
    }

    private func getDirs() -> String {
        let fileManager = FileManager.default
        var dirs: [String: String] = [:]

        if let documentDir = fileManager.urls(for: .documentDirectory, in: .userDomainMask).first {
            dirs["documents"] = documentDir.path
        }

        if let cacheDir = fileManager.urls(for: .cachesDirectory, in: .userDomainMask).first {
            dirs["cache"] = cacheDir.path
        }

        if let libraryDir = fileManager.urls(for: .libraryDirectory, in: .userDomainMask).first {
            dirs["library"] = libraryDir.path
        }

        dirs["temp"] = NSTemporaryDirectory()

        if let data = try? JSONSerialization.data(withJSONObject: dirs),
           let json = String(data: data, encoding: .utf8) {
            return json
        }
        return "{}"
    }
}
