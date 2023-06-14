chrome.runtime.onMessage.addListener(function (request, sender, sendResponse) {
    if (message.action === 'getNativeData') {
        // Perform platform-specific operations or interact with native code here
        // Send the result back to the content script

        sendResponse({ data: 'Native data' });
    }

    // This script will handle communication between the extension and the native messaging host. 
    // Here's an example of a background.js script that sends a message to the native host and receives a response:
    chrome.runtime.sendNativeMessage('dev.codemonkeyninja.manga-furigana', request, function (response) {
        sendResponse(response);
    });
    return true;
});
