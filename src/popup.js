document.addEventListener("DOMContentLoaded", function () {
    // convert each image on the page to text using OCR
    chrome.tabs.query({ active: true, currentWindow: true }, function (tabs) {
        // create a canvas element to draw the image to
        var canvas = document.createElement("canvas");
        var context = canvas.getContext("2d");

        chrome.tabs.sendMessage(tabs[0].id, { greeting: "greeting:popup.js" }, function (response) {
            if (response && response.images) {
                // loop through each image on the page
                for (var i = 0; i < response.images.length; i++) {
                    var image = response.images[i];

                    // draw the image to the canvas
                    canvas.width = image.width;
                    canvas.height = image.height;
                    context.drawImage(image, 0, 0);

                    // perform OCR on the image using Tesseract.js
                    Tesseract.recognize(canvas.toDataURL()).then(function (result) {
                        // translate the OCR result to hiragana and katakana
                        var japaneseText = result.text;
                        var phoneticText = japanese.toHiragana(japaneseText) + japanese.toKatakana(japaneseText);

                        // create a text box below the image with the phonetic text
                        var textBox = document.createElement("textarea");
                        textBox.value = phoneticText;
                        image.parentNode.insertBefore(textBox, image.nextSibling);
                    });
                }
            }
        });
    });
});