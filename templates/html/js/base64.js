{{ define "js/base64" }}

'use strict';

function base64urlToBuffer(baseurl64String) {
    // Base64url to Base64
    const padding = "==".slice(0, (4 - (baseurl64String.length % 4)) % 4);
    const base64String = baseurl64String.replace(/-/g, "+").replace(/_/g, "/") + padding;

    // Base64 to binary string
    const str = atob(base64String);

    // Binary string to buffer
    const buffer = new ArrayBuffer(str.length);
    const byteView = new Uint8Array(buffer);
    for (let i = 0; i < str.length; i++) {
        byteView[i] = str.charCodeAt(i);
    }

    return buffer;
}

function bufferToBase64url(buffer) {
    // Buffer to binary string
    const byteView = new Uint8Array(buffer);
    let str = "";
    for (const charCode of byteView) {
        str += String.fromCharCode(charCode);
    }

    // Binary string to base64
    const base64String = btoa(str);

    // Base64 to base64url
    return base64String.replace(/\+/g, "-").replace(/\//g, "_").replace(/=/g, "");
}

{{ end }}
