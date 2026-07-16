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
            let status = SweetjuiceStartApplication()
            print("[Go] Start Status: \(status)")

            // Register native call handler for Go-to-Native calls
            SweetjuiceSetNativeCallHandler(context.coordinator)

            SweetJuiceWebView.hasStartedGo = true
        }

        // Attach plugins to a view controller
        if let windowScene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
           let rootVC = windowScene.windows.first?.rootViewController {
            for plugin in SweetJuiceManager.shared.getPlugins() {
                plugin.onAttach(container: rootVC)
            }
        }

        // Start event polling
        context.coordinator.startPolling(webView: webView)

        // Load application entry point
        let url = URL(string: "sweetjuice://localhost/index.html")!
        webView.load(URLRequest(url: url))

        return webView
    }

    func updateUIView(_ uiView: WKWebView, context: Context) {}

    func makeCoordinator() -> Coordinator {
        Coordinator()
    }

    // Coordinator class with Swift 6 concurrency safety
    @MainActor
    class Coordinator: NSObject, WKNavigationDelegate, WKScriptMessageHandler, SweetjuiceNativeCallHandlerProtocol {
        private var isPolling = true

        // nonisolated to satisfy SweetjuiceNativeCallHandlerProtocol (called from Go thread)
        nonisolated func onNativeCall(_ method: String?, args: String?) -> String {
            let methodKey = method ?? ""
            let jsonArgs = args ?? ""

            var response = ""
            // Synchronously hop to main thread to call plugins
            DispatchQueue.main.sync {
                response = SweetJuiceManager.shared.handleNativeCall(method: methodKey, args: jsonArgs)
            }
            return response
        }

        func userContentController(_ userContentController: WKUserContentController, didReceive message: WKScriptMessage) {
            guard message.name == "SweetJuiceBind",
                  let body = message.body as? [String: String],
                  let methodKey = body["methodKey"],
                  let jsonArgs = body["jsonArgs"] else { return }

            _ = SweetjuiceHandleMessageFromFrontend(methodKey, jsonArgs)
        }

        func startPolling(webView: WKWebView) {
            Task.detached { [weak self] in
                while true {
                    guard let self = self else { break }

                    let shouldContinue = await self.getPollingStatus()
                    if !shouldContinue { break }

                    let eventJson = SweetjuicePollNativeEvent()
                    if !eventJson.isEmpty {
                        await MainActor.run {
                            let script = "if(window.SweetJuiceBind && window.SweetJuiceBind.dispatch) { window.SweetJuiceBind.dispatch(\(eventJson)); }"
                            webView.evaluateJavaScript(script)
                        }
                    }

                    try? await Task.sleep(nanoseconds: 100_000_000) // 100ms
                }
            }
        }

        func getPollingStatus() -> Bool {
            return isPolling
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
              window.SweetJuice = window.SweetJuiceBind;
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

        // Robust path extraction from sweetjuice://localhost/...
        var assetKey = url.path
        if assetKey.hasPrefix("/") { assetKey = String(assetKey.dropFirst()) }
        if assetKey.isEmpty { assetKey = "index.html" }

        // Fetch asset from Go
        let fileData = SweetjuiceRequestAssetBytes(assetKey)
        let mimeType = SweetjuiceRequestAssetMime(assetKey)

        // Ensure non-empty MIME type for iOS strictness
        let finalMime = mimeType.isEmpty ? "application/octet-stream" : mimeType
        let finalData = fileData ?? Data()

        let response = HTTPURLResponse(
            url: url,
            statusCode: fileData != nil ? 200 : 404,
            httpVersion: nil,
            headerFields: [
                "Content-Type": finalMime,
                "Access-Control-Allow-Origin": "*"
            ]
        )!

        urlSchemeTask.didReceive(response)
        urlSchemeTask.didReceive(finalData)
        urlSchemeTask.didFinish()
    }

    func webView(_ webView: WKWebView, stop urlSchemeTask: WKURLSchemeTask) {}
}
