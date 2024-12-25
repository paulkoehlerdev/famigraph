{{ define "js/register" }}

'use strict';

async function requestPublicKeyCredentialCreationOptions() {
    const rawResponse = await fetch(`/register/challenge`).catch(() => null);
    if (rawResponse === null || !rawResponse.ok) {
        return null;
    }

    const response = await rawResponse.json();

    return {
        rp: response.rp,
        user: {
            id: base64urlToBuffer(response.user.id),
            name: response.user.name,
            displayName: response.user.displayName,
        },
        challenge: base64urlToBuffer(response.challenge),
        pubKeyCredParams: response.pubKeyCredParams,
        timeout: response.timeout,
        authenticatorSelection: response.authenticatorSelection,
        attestation: response.attestation,
    };
}

async function postPublicKeyCredential(credential) {
    const encodedCredential = {
        authenticatorAttachment: credential.authenticatorAttachment,
        id: credential.id,
        rawId: bufferToBase64url(credential.rawId),
        response: {
            attestationObject: bufferToBase64url(credential.response.attestationObject),
            clientDataJSON: bufferToBase64url(credential.response.clientDataJSON),
        },
        type: credential.type,
    };

    const response = await fetch(`/register/challenge`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(encodedCredential),
    }).catch(() => null);

    if (response.ok) {
        return null;
    }

    return await response.text();
}

async function register() {
    hideErrorMessage();

    const publicKeyCredentialCreationOptions = await requestPublicKeyCredentialCreationOptions();
    if (publicKeyCredentialCreationOptions === null) {
        displayErrorMessage('Registration challenge request failed.', '');
        return;
    }

    let publicKeyCredential;
    try {
        publicKeyCredential = await navigator.credentials.create({
            publicKey: publicKeyCredentialCreationOptions,
        });
    } catch (e) {
        displayErrorMessage('A client-side error occurred during challenge signing.', e);
        return;
    }

    const errorMessage = await postPublicKeyCredential(publicKeyCredential);
    if (errorMessage !== null) {
        displayErrorMessage('The server rejected your signed challenge.', errorMessage);
        return;
    }

    window.location.replace("/");
}

function main() {
    register();
}

main();

{{ end }}