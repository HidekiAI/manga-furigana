// check if the operating system is Windows
function isWindows() {
    return navigator.userAgent.indexOf("Win") != -1;
}

// check if the operating system is Linux
function isLinux() {
    return navigator.userAgent.indexOf("Linux") != -1;
}

// load the appropriate Kakasi library based on the operating system
var script = document.createElement("script");
if (isWindows()) {
    script.src = chrome.runtime.getURL("kakasi-windows/kakasi.js");
} else if (isLinux()) {
    script.src = chrome.runtime.getURL("kakasi-linux/kakasi.js");
}
document.head.appendChild(script);


// content.js
chrome.runtime.sendMessage({ action: 'getNativeData' }, function (response) {
    console.log(response.data);
});


// This script will be injected into web pages and interact with the background script
// Here's an example content.js script that sends a message to the background script:
chrome.runtime.sendMessage({ message: 'Hello from content.js' }, function(response) {
    console.log('Response from background:', response);
  });
