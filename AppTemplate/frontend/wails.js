/**
 * wails.js - The Core Bridge Contract
 * This file translates native calls to the frontend.
 * IT SHOULD NOT BE MODIFIED BY THE USER, UNLESS YOU KNOW WHAT YOU ARE DOING.
 */
window.Wails = {
    /**
     * CallGo invokes a method bound in the Go backend.
     */
    CallGo: async function(methodName, ...args) {
        if (!window.WailsBind) throw new Error('Native bridge not found.');
        const jsonArgs = JSON.stringify(args);
        const rawResponse = window.WailsBind.callGo(methodName, jsonArgs);
        const parsed = JSON.parse(rawResponse);
        if (parsed.error) throw new Error(parsed.error);
        return parsed.result;
    },
    /**
     * on registers a listener for events emitted from the Go backend.
     */
    on: function(name, cb) {
        const attachListener = () => {
            if (window.WailsBind && window.WailsBind.on) {
                window.WailsBind.on(name, cb);
            } else if (window.WailsEvents && window.WailsEvents.on) {
                window.WailsEvents.on(name, cb);
            } else {
                setTimeout(attachListener, 100);
            }
        };
        attachListener();
    }
};
