import Foundation
import UIKit

public protocol SweetJuicePlugin {
    func getDomain() -> String
    func onAttach(container: UIViewController)
    func handleAction(action: String, jsonArgs: String) -> String
}
