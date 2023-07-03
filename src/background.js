chrome.runtime.onInstalled.addListener(function () {
    console.log('Extension installed or updated.');
    // Perform any necessary setup or initialization tasks here
});

chrome.runtime.onMessage.addListener(function (message) {
    console.log('Received message:', message);
    // Handle the incoming message and perform any necessary actions here
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

chrome.tabs.onUpdated.addListener(function (tabId, changeInfo, tab) {
    console.log('Tab updated:', tab);
    // Perform any necessary actions based on the state of the tab here
});

//chrome.webRequest.onBeforeRequest.addListener(function (details) {
//    console.log('Request intercepted:', details);
//    // Modify or block the request as needed here
//    return { cancel: true };
//}, { urls: ['<all_urls>'] }, ['blocking']);
