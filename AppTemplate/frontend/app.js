/**
 * app.js - Wails Mobile Command Center Logic
 */
function updateOutput(text) {
    const output = document.getElementById('output');
    if (output.textContent === "Ready for commands...") {
        output.textContent = "";
    }
    output.textContent = `[${new Date().toLocaleTimeString()}] ${text}\n` + output.textContent;
}

window.addEventListener('DOMContentLoaded', () => {

    // --- Theme Toggle ---
    const themeToggle = document.getElementById('theme-toggle');
    const html = document.documentElement;

    themeToggle.addEventListener('click', () => {
        if (html.classList.contains('dark')) {
            html.classList.remove('dark');
            html.classList.add('light');
        } else {
            html.classList.remove('light');
            html.classList.add('dark');
        }
    });

    // --- Mouse Follow Polish ---
    document.addEventListener('mousemove', (e) => {
        const x = (e.clientX / window.innerWidth - 0.5) * 40;
        const y = (e.clientY / window.innerHeight - 0.5) * 40;

        const glow = document.querySelector('.absolute.bg-primary\\/10');
        if (glow) {
            glow.style.transform = `translate(${x}px, ${y}px)`;
        }
    });

    // --- Clear Logs ---
    document.getElementById('clear-logs').addEventListener('click', () => {
        document.getElementById('output').textContent = "Ready for commands...";
    });


    // --- Section 1: Notifications ---

    document.getElementById('request-notif').addEventListener('click', async () => {
        try {
            await Wails.CallGo('PermissionPlugin.Request', "android.permission.POST_NOTIFICATIONS");
            updateOutput("Requested notification permission.");
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('send-notif').addEventListener('click', async () => {
        try {
            const result = await Wails.CallGo('NotificationPlugin.Post', {
                id: 0,
                title: "Wails Mobile",
                body: "This is a unique verification notification.",
                importance: "HIGH"
            });
            const parsed = typeof result === 'string' ? JSON.parse(result) : result;
            updateOutput(`Notification posted (ID: ${parsed.id})`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });


    // --- Section 2: Biometrics ---

    document.getElementById('check-bio').addEventListener('click', async () => {
        try {
            const res = await Wails.CallGo('BiometricPlugin.CanAuthenticate');
            updateOutput(`Biometric Status: ${res.status} (Available: ${res.can_authenticate})`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('auth-bio').addEventListener('click', async () => {
        try {
            updateOutput("Starting biometric authentication...");
            await Wails.CallGo('BiometricPlugin.Authenticate', {
                title: "Verify Identity",
                subtitle: "Confirm biometric to proceed",
                description: "This test ensures the native prompt bridge is working.",
                negative_button_text: "Cancel"
            });
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });


    // --- Section 3: File Picker ---

    document.getElementById('pick-image').addEventListener('click', async () => {
        try {
            updateOutput("Opening image picker...");
            await Wails.CallGo('FilePickerPlugin.PickFile', {
                mime_type: "image/*",
                multiple: true
            });
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });


    // --- Section 4: Daemon ---

    document.getElementById('start-daemon')?.addEventListener('click', async () => {
        try {
            updateOutput("Starting background daemon...");
            await Wails.CallGo('DaemonPlugin.Start', {
                title: "Wails Mobile",
                message: "Core engine active in background",
                importance: "LOW"
            });
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('stop-daemon')?.addEventListener('click', async () => {
        try {
            updateOutput("Stopping background daemon...");
            await Wails.CallGo('DaemonPlugin.Stop');
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });


    // --- Event Listeners ---

    Wails.on('biometric:result', (result) => {
        updateOutput(`BIOMETRIC EVENT: ${result.status} (Success: ${result.success})`);
        if (result.error) updateOutput(`Detail: ${result.error}`);
    });

    Wails.on('filepicker:result', (result) => {
        if (result.error) {
            updateOutput(`FILEPICKER EVENT: ${result.error}`);
            return;
        }
        if (result.multiple) {
            updateOutput(`FILEPICKER EVENT: Selected ${result.uris.length} files`);
        } else {
            updateOutput(`FILEPICKER EVENT: Selected ${result.uri}`);
        }
    });

    Wails.on('permissions:changed', (data) => {
        updateOutput(`PERMISSION EVENT: ${data.permission} is ${data.granted ? 'GRANTED' : 'DENIED'}`);
    });

});
