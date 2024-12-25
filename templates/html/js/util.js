{{ define "js/util" }}

'use strict';

function displayErrorMessage(message = '', errorDetails = '') {
    document.querySelector('div.error > p.error__message').innerText = message;
    document.querySelector('div.error > p.error__details').innerText = errorDetails;
    document.querySelector('div.error').classList.remove('hidden');

    localStorage.removeItem('username');
}

function hideErrorMessage() {
    if (!document.querySelector('div.error').classList.contains('hidden')) {
        document.querySelector('div.error').classList.add('hidden');
    }
}

{{ end }}
