/**
 * juice.js - The Core Bridge Contract
 * This file translates native calls to the frontend.
 * IT SHOULD NOT BE MODIFIED BY THE USER, UNLESS YOU KNOW WHAT YOU ARE DOING.
 */
window.SweetJuice = {
    /**
     * CallGo invokes a method bound in the Go backend.
     */
    CallGo: async function(methodName, ...args) {
        if (!window.SweetJuiceBind) throw new Error('Native bridge not found.');
        const jsonArgs = JSON.stringify(args);
        const rawResponse = window.SweetJuiceBind.callGo(methodName, jsonArgs);
        const parsed = JSON.parse(rawResponse);
        if (parsed.error) throw new Error(parsed.error);
        return parsed.result;
    },
    /**
     * on registers a listener for events emitted from the Go backend.
     */
    on: function(name, cb) {
        const attachListener = () => {
            if (window.SweetJuiceBind && window.SweetJuiceBind.on) {
                window.SweetJuiceBind.on(name, cb);
            } else if (window.SweetJuiceEvents && window.SweetJuiceEvents.on) {
                window.SweetJuiceEvents.on(name, cb);
            } else {
                setTimeout(attachListener, 100);
            }
        };
        attachListener();
    }
};
