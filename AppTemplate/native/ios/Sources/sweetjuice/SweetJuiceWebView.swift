import SwiftUI
import WebKit
import Sweetjuice

struct SweetJuiceWebView: UIViewRepresentable {
    static var hasStartedGo = false

    func makeUIView(context: Context) -> WKWebView {
        let config = WKWebViewConfiguration()

        // Add custom scheme handler for sweetjuice://
        let handler = SweetJuiceSchemeHandler()
        config.setURLSchemeHandler(handler, forURLScheme: "sweetjuice")

        // Setup UserContentController for JS -> Go bridge
        let userContentController = WKUserContentController()
        userContentController.add(context.coordinator, name: "SweetJuiceBind")
        config.userContentController = userContentController

        let webView = WKWebView(frame: .zero, configuration: config)
        webView.navigationDelegate = context.coordinator

        #if DEBUG
        if #available(iOS 16.4, *) {
            webView.isInspectable = true
        }
        #endif

        if !SweetJuiceWebView.hasStartedGo {
            // Start Go Application
            _ = SweetjuiceStartApplication()

            // Register native call handler for Go-to-Native calls
            SweetjuiceSetNativeCallHandler(context.coordinator)

            SweetJuiceWebView.hasStartedGo = true
        }

        // Attach plugins to a view controller if possible
        DispatchQueue.main.async {
            if let windowScene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
               let rootVC = windowScene.windows.first?.rootViewController {
                for plugin in SweetJuiceManager.shared.getPlugins() {
                    plugin.onAttach(container: rootVC)
                }
            }
        }

        // Start event polling
        context.coordinator.startPolling(webView: webView)

        let url = URL(string: "sweetjuice://index.html")!
        webView.load(URLRequest(url: url))

        return webView
    }

    func updateUIView(_ uiView: WKWebView, context: Context) {}

    func makeCoordinator() -> Coordinator {
        Coordinator()
    }

    class Coordinator: NSObject, WKNavigationDelegate, WKScriptMessageHandler, SweetjuiceNativeCallHandlerProtocol {
        private var isPolling = true

        func onNativeCall(_ method: String?, args: String?) -> String {
            return SweetJuiceManager.shared.handleNativeCall(method: method ?? "", args: args ?? "")
        }

        func userContentController(_ userContentController: WKUserContentController, didReceive message: WKScriptMessage) {
            guard message.name == "SweetJuiceBind",
                  let body = message.body as? [String: String],
                  let methodKey = body["methodKey"],
                  let jsonArgs = body["jsonArgs"] else { return }

            let result = SweetjuiceHandleMessageFromFrontend(methodKey, jsonArgs)
            // Send result back to JS if needed
            // For now, we can use evaluateJavaScript to return results if we implement a promise-based call in JS
        }

        func startPolling(webView: WKWebView) {
            Thread.detachNewThread {
                while self.isPolling {
                    let eventJson = SweetjuicePollNativeEvent()
                    if !eventJson.isEmpty {
                        DispatchQueue.main.async {
                            let script = "if(window.SweetJuiceBind && window.SweetJuiceBind.dispatch) { window.SweetJuiceBind.dispatch(\(eventJson)); }"
                            webView.evaluateJavaScript(script)
                        }
                    }
                    Thread.sleep(forTimeInterval: 0.1)
                }
            }
        }

        func webView(_ webView: WKWebView, didFinish navigation: WKNavigation!) {
            let js = """
            if (!window.SweetJuiceEvents) {
              window.SweetJuiceEvents = {
                listeners: {},
                on: function(name, cb) {
                  if(!this.listeners[name]) this.listeners[name] = [];
                  this.listeners[name].push(cb);
                },
                dispatch: function(obj) {
                  var name = obj.name; var data = obj.data;
                  if(this.listeners[name]) {
                    this.listeners[name].forEach(function(cb) { try { cb(data); } catch(e) { console.error(e); } });
                  }
                }
              };
              window.SweetJuiceBind = {
                callGo: function(methodKey, jsonArgs) {
                  window.webkit.messageHandlers.SweetJuiceBind.postMessage({methodKey: methodKey, jsonArgs: jsonArgs});
                },
                on: window.SweetJuiceEvents.on.bind(window.SweetJuiceEvents),
                dispatch: window.SweetJuiceEvents.dispatch.bind(window.SweetJuiceEvents)
              };
            }
            """
            webView.evaluateJavaScript(js)
        }

        deinit {
            isPolling = false
        }
    }
}

class SweetJuiceSchemeHandler: NSObject, WKURLSchemeHandler {
    func webView(_ webView: WKWebView, start urlSchemeTask: WKURLSchemeTask) {
        guard let url = urlSchemeTask.request.url else { return }

        var assetKey = url.absoluteString.replacingOccurrences(of: "sweetjuice://", with: "")
        if assetKey.contains("?") { assetKey = String(assetKey.split(separator: "?")[0]) }
        if assetKey.contains("#") { assetKey = String(assetKey.split(separator: "#")[0]) }
        if assetKey.isEmpty || assetKey == "/" { assetKey = "index.html" }

        let fileData = SweetjuiceRequestAssetBytes(assetKey)
        let mimeType = SweetjuiceRequestAssetMime(assetKey)

        let finalMime = mimeType.isEmpty ? "text/plain" : mimeType
        let finalData = fileData ?? Data()

        let response = HTTPURLResponse(url: url, statusCode: 200, httpVersion: nil, headerFields: ["Content-Type": finalMime])!

        urlSchemeTask.didReceive(response)
        urlSchemeTask.didReceive(finalData)
        urlSchemeTask.didFinish()
    }

    func webView(_ webView: WKWebView, stop urlSchemeTask: WKURLSchemeTask) {}
}
