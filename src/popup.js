document.addEventListener('DOMContentLoaded', function () {
    //document.getElementById("scanButton").addEventListener('click', scanImages);
    var scanButton = document.getElementById("scanButton");
    console.log(scanButton);
    scanButton.addEventListener('click', scanImages);
});

function scanImages() {
    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
        chrome.runtime.sendNativeMessage('manga-furigana.codemonkeyninja.dev', { tabId: tabs[0].id });
    });

    // [MY_CLIENT_ID] is a Client ID from one of your credentials in the Google Cloud console.
    // get clientID to tbe used for authentication against accounts.google.com oauth2
    clientId = "1094731870219-e5nljpbhmf5e0hu88erijmint450esft.apps.googleusercontent.com";     // client_id=[MY_CLIENT_ID].apps.googleusercontent.com

    // Note: Edge does not support chrome.identity.getAuthToken - As an alternate, you can use launchWebAuthFlow to fetch an OAuth2 token to authenticate users.
    // note that documentation only use scope of 'scope=https://www.googleapis.com/auth/cloud-platform'
    const scopes = [
        'https://www.googleapis.com/auth/cloud-platform',   // scope=https://www.googleapis.com/auth/cloud-platform
        //'https://www.googleapis.com/oauth2/v1/certs',
        //'https://accounts.google.com/o/oauth2/auth',
        //'https://oauth2.googleapis.com/token',
    ];  // cloud-platform or cloud-vision?

    // [MY_REDIRECT_URI] is the corresponding Authorized redirect URIs from the same credential in the Google Cloud console. 
    // If no redirect URI is specified, you must specify a trusted URI, for example https://www.google.com. 
    // The redirect URI defines where the HTTP response is sent. For production, you must specify your application's 
    // auth endpoint, which handles responses from the OAuth 2.0 server. For more information, see Using OAuth 2.0 
    // for Web Server Applications.
    redirect_url = chrome.identity.getRedirectURL();     // redirect_uri=[MY_REDIRECT_URI]
    if ((redirect_url == null) || (redirect_url == "undefiled") || (redirect_url == "")) {
        redirect_url = "https://localhost:8080/";
    }

    // log clientID and redirect_url
    console.log(clientId);
    console.log(redirect_url);


    // sample requesting URL from 'https://cloud.google.com/appengine/docs/admin-api/accessing-the-api' documentation:
    //      https://accounts.google.com/o/oauth2/v2/auth?
    //          response_type=token&
    //          client_id=[MY_CLIENT_ID].apps.googleusercontent.com&
    //          scope=https://www.googleapis.com/auth/cloud-platform&
    //          redirect_uri=[MY_REDIRECT_URI]
    urlNoScopes = "https://accounts.google.com/o/oauth2/v2/auth?response_type=token&client_id=" + clientId + "&redirect_uri=" + encodeURIComponent(redirect_url);
    urlAppendedWithScopes = urlNoScopes + "&scope=" + scopes.join(' ');
    // `https://accounts.google.com/o/oauth2/auth?client_id=${clientId}&redirect_uri=${encodeURIComponent(redirect_url)}&response_type=token&scope=${scopes.join(' ')}`,
    console.log(urlAppendedWithScopes);
    chrome.identity.launchWebAuthFlow(
        {
            url: urlAppendedWithScopes,
            'interactive': true
        },
        function (redirect_url) {
            if (chrome.runtime.lastError || !redirect_url) {
                console.log(chrome.runtime.lastError);
                alert('Failed to authenticate with Google Cloud Platform. Please try again.');
                return;
            }

            // extract the access token from redirect_url
            const accessToken = redirect_url.match(/access_token=([^&]+)/)[1];

            // Use the access token to authenticate requests to the Google Cloud Vision API
            var xhr = new XMLHttpRequest();
            xhr.open('POST', 'https://vision.googleapis.com/v1/images:annotate');
            xhr.setRequestHeader('Authorization', 'Bearer ' + token);
            xhr.setRequestHeader('Content-Type', 'application/json; charset=utf-8');
            xhr.setRequestHeader('x-content-type-options', 'nosniff'); // Add x-content-type-options header
            xhr.onload = function () {
                // Handle the response
                const response = JSON.parse(xhr.responseText);
                console.log(response);

                // only continue if response is valid
                if (!response || !response.responses || !response.responses[0] || !response.responses[0].textAnnotations) {
                    return;

                    // collect all the images from this current tab
                    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
                        // get ID of the current tab that is in focus
                        const currentTab = tabs.findIndex(tab => tab.active);
                        const tabId = tabs[currentTab].id;

                        // Get all the images on the current tab
                        chrome.tabs.executeScript(tabId, { code: 'Array.from(document.images).map(img => img.src)' }, function (result) {
                            // Get the images from the result
                            const images = result[0] || [];

                            // Send the images to the native Go background logic
                            for (const image of images) {
                                const xhrPerImage = new XMLHttpRequest();
                                xhrPerImage.open('GET', image, true);
                                xhrPerImage.responseType = 'blob';
                                xhrPerImage.onload = function () {
                                    const blob = xhrPerImage.response;
                                    const reader = new FileReader();
                                    reader.onload = function () {
                                        const binaryData = reader.result;
                                        chrome.runtime.sendNativeMessage('manga-furigana.codemonkeyninja.dev', { tabId, image: binaryData });
                                    };
                                    reader.readAsArrayBuffer(blob);
                                };
                                xhrPerImage.send();
                            }
                        });
                    });
                    chrome.runtime.sendMessage({ images });
                };
                xhr.send();
            };
        });
};

chrome.runtime.onMessage.addListener(function (message) {
    if (message.images) {
        const resultArea = document.getElementById('resultArea');
        resultArea.innerHTML = '';

        for (const image of message.images) {
            const textbox = document.createElement('textarea');
            textbox.value = image.tokenizedText;
            resultArea.appendChild(textbox);
        }
    }
});
