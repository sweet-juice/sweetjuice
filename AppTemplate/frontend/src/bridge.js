/**
 * bridge.js - SweetJuice Native Bridge
 */

const isNative = !!(window.SweetJuiceBind || (window.webkit && window.webkit.messageHandlers && window.webkit.messageHandlers.SweetJuiceBind));

export const SweetJuice = {
    /**
     * CallGo invokes a method bound in the Go backend.
     */
    CallGo: async function(methodName, ...args) {
        if (!isNative) {
            console.warn(`[Browser Mode] CallGo: ${methodName}`, args);
            return { error: "Native bridge not found" };
        }

        const jsonArgs = JSON.stringify(args);

        // Android (SweetJuiceBind.callGo returns string)
        if (window.SweetJuiceBind && window.SweetJuiceBind.callGo) {
            const rawResponse = window.SweetJuiceBind.callGo(methodName, jsonArgs);
            try { return JSON.parse(rawResponse); } catch(e) { return rawResponse; }
        }

        // iOS (webkit.messageHandlers.SweetJuiceBind.postMessage)
        if (window.webkit && window.webkit.messageHandlers && window.webkit.messageHandlers.SweetJuiceBind) {
            window.webkit.messageHandlers.SweetJuiceBind.postMessage({methodKey: methodName, jsonArgs: jsonArgs});
            return { status: "sent" };
        }

        return { error: "Unknown bridge type" };
    },

    /**
     * on registers a listener for events emitted from the Go backend.
     */
    on: function(name, cb) {
        if (!window.SweetJuiceEvents) {
            window.SweetJuiceEvents = {
                listeners: {},
                on: function(name, cb) {
                    if(!this.listeners[name]) this.listeners[name] = [];
                    this.listeners[name].push(cb);
                },
                dispatch: function(obj) {
                    const name = obj.name;
                    const data = obj.data;
                    if(this.listeners[name]) {
                        this.listeners[name].forEach(function(cb) {
                            try { cb(data); } catch(e) { console.error(e); }
                        });
                    }
                }
            };
        }
        window.SweetJuiceEvents.on(name, cb);
    }
};

window.SweetJuice = SweetJuice;
// Compatibility with old templates
window.SweetJuice = SweetJuice;
