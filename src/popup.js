document.addEventListener('DOMContentLoaded', function () {
    document.getElementById('scanButton').addEventListener('click', scanImages);
});

function scanImages() {
    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
        chrome.runtime.sendNativeMessage('manga-furigana.codemonkeyninja.dev', { tabId: tabs[0].id });
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
