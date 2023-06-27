document.addEventListener('DOMContentLoaded', function () {
    document.getElementById('scanButton').addEventListener('click', scanImages);
});

function scanImages() {
    chrome.identity.getAuthToken({ 'interactive': true }, function (token) {
        // Use the access token to authenticate requests to the Google Cloud Vision API
        var xhr = new XMLHttpRequest();
        xhr.open('POST', 'https://vision.googleapis.com/v1/images:annotate');
        xhr.setRequestHeader('Authorization', 'Bearer ' + token);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.onload = function () {
            // Handle the response
        };
        xhr.send(JSON.stringify(request));
    });

    chrome.identity.getAuthToken({ interactive: true }, function (token) {
        // Use the token to authenticate requests to the Google Cloud Vision API
        var xhr = new XMLHttpRequest();
        xhr.open('POST', 'https://vision.googleapis.com/v1/images:annotate');
        xhr.setRequestHeader('Authorization', 'Bearer ' + token);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.onload = function () {
            // Handle the response
        };
        xhr.send(JSON.stringify(request));
    });

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
                const xhr = new XMLHttpRequest();
                xhr.open('GET', image, true);
                xhr.responseType = 'blob';
                xhr.onload = function () {
                    const blob = xhr.response;
                    const reader = new FileReader();
                    reader.onload = function () {
                        const binaryData = reader.result;
                        chrome.runtime.sendNativeMessage('manga-furigana.codemonkeyninja.dev', { tabId, image: binaryData });
                    };
                    reader.readAsArrayBuffer(blob);
                };
                xhr.send();
            }
        });
    });
}

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

//chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
//    if (message.action === 'getNativeData') {
//        // Perform platform-specific operations or interact with native code here
//        // Send the result back to the content script
//
//        sendResponse({ data: 'Native data' });
//    }
//
//    // This script will handle communication between the extension and the native messaging host. 
//    // Here's an example of a background.js script that sends a message to the native host and receives a response:
//    chrome.runtime.sendNativeMessage('manga-furigana.codemonkeyninja.dev', request, function (response) {
//        sendResponse(response);
//    });
//    return true;
//});
//