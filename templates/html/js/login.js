{{ define "js/login" }}

'use strict';

async function requestPublicKeyCredentialRequestOptions() {
    const rawResponse = await fetch(`/login/challenge`).catch(() => null);
    if (rawResponse === null || !rawResponse.ok) {
        return null;
    }

    const response = await rawResponse.json();

    let publicKeyCredentialRequestOptions = {
        rpId: response.rpId,
        challenge: base64urlToBuffer(response.challenge),
        timeout: response.timeout,
        userVerification: response.userVerification,
    };

    return publicKeyCredentialRequestOptions;
}

async function postPublicKeyCredential(credential) {
    const encodedCredential = {
        authenticatorAttachment: credential.authenticatorAttachment,
        id: credential.id,
        rawId: bufferToBase64url(credential.rawId),
        response: {
            authenticatorData: bufferToBase64url(credential.response.authenticatorData),
            clientDataJSON: bufferToBase64url(credential.response.clientDataJSON),
            signature: bufferToBase64url(credential.response.signature),
            userHandle: bufferToBase64url(credential.response.userHandle),
        },
        type: credential.type,
    };

    const response = await fetch(`/login/challenge`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(encodedCredential),
    }).catch(() => null);

    return response !== null && response.ok;
}

async function login() {
    const publicKeyCredentialRequestOptions = await requestPublicKeyCredentialRequestOptions();
    if (publicKeyCredentialRequestOptions === null) {
        displayErrorMessage('Login challenge request failed.', '');
        return;
    }

    let publicKeyCredential;
    try {
        publicKeyCredential = await navigator.credentials.get({
            publicKey: publicKeyCredentialRequestOptions,
        });
    } catch (e) {
        displayErrorMessage('A client-side error occurred during challenge signing.', e);
        return;
    }

    const loginSuccessful = await postPublicKeyCredential(publicKeyCredential);
    if (!loginSuccessful) {
        displayErrorMessage('Permission denied! Your credential is not authorized.', '');
        return;
    }

    const url = new URL(window.location.href);
    const params = new URLSearchParams(url.search);
    const loc = params.get('loc');

    if (loc !== null) {
        window.location.replace(atob(loc));
    } else {
        window.location.replace('/');
    }
}

function main() {
    login();
}

main();

{{ end }}