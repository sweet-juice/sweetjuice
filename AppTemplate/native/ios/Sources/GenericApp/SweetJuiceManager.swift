import Foundation
import UIKit
import Sweetjuice

@MainActor
public class SweetJuiceManager: Sendable {
    public static let shared = SweetJuiceManager()

    private var plugins: [String: SweetJuicePlugin] = [:]

    private init() {}

    public func registerPlugin(_ plugin: SweetJuicePlugin) {
        plugins[plugin.getDomain()] = plugin
    }

    public func getPlugins() -> [SweetJuicePlugin] {
        return Array(plugins.values)
    }

    public func handleNativeCall(method: String, args: String) -> String {
        if method.contains(":") {
            let parts = method.components(separatedBy: ":")
            if parts.count >= 2 {
                let domain = parts[0]
                let action = parts[1]

                if let plugin = plugins[domain] {
                    return plugin.handleAction(action: action, jsonArgs: args)
                }
            }
        }
        return "{\"error\":\"Plugin domain not found\"}"
    }
}
