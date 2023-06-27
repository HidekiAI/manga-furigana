# manga-furigana: ocr-based real-time dynamic furigana for manga plugin to chromium/msedge

This is a Chrome extension that uses the Google Cloud Vision API to perform OCR on manga images and display furigana (Japanese pronunciation guides) in real-time.

## Prerequisites

Before you can use this extension, you will need to set up a Google Cloud Vision API project and obtain a service account key file in JSON format. Here are the general steps to do this:

1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project or select an existing project.
3. Enable the Google Cloud Vision API for your project.
4. Create a service account for your project and download the service account key file in JSON format.
5. Store the service account key file in a secure location, such as a private repository or a secure file storage service, and restrict access to the file to only those users who need it.

TODO: Determine if we need to export the account key (JSON) file (steps 4 and 5 above) since we'll be using the Google Identity API to obtain an OAuth2 access token for the Google Cloud Vision API.  If we do need to export the account key (JSON) file, then probably, we'll need to create a proxy server that will do server-to-server and on that server, the configuration key JSON file will only reside (privately).  And then, have client-to-server communication (browser-to-server) to get the access token.  This way, the configuration key JSON file will not be exposed to the client (browser).  The proxy server will be the one to communicate with the Google Cloud Vision API.  The client (browser) will only communicate with the proxy server.  The proxy server will be the one to communicate with the Google Cloud Vision API.  The client (browser) will only communicate with the proxy server.  The proxy server will be the one to communicate with the Google Cloud Vision API.  The client (browser) will only communicate with the proxy server.  The proxy server will be the one to communicate with the Google Cloud Vision API.  The client (browser) will only communicate with the proxy server.  The proxy server will be the one to communicate with the Google Cloud Vision API.  The client (browser) will only communicate with the proxy server.  The proxy server will be the one to communicate with the Google Cloud Vision API.  The client (browser) will only communicate with the proxy server.  The proxy server will be the one to communicate with the Google Cloud Vision API.  The client (browser) will only communicate with the proxy server.  The proxy server will be the one to communicate with the Google Cloud Vision API.  The client (browser) will only communicate with the proxy server.

## Authentications (behind the scene)

When using the Google Identity API to obtain an OAuth2 access token for the Google Cloud Vision API, the user will be prompted to sign in to their Google account if they are not already signed in. The user will then be prompted to grant permission to the extension to access the Google Cloud Vision API using the OAuth2 scopes that were specified in the `manifest.json` file.

The user's Google account does not need to be the same account that was used to create the Google Cloud Vision API project. However, the user must have permission to access the Google Cloud Vision API using the OAuth2 scopes that were specified in the `manifest.json` file.

When the user grants permission to the extension to access the Google Cloud Vision API, an OAuth2 access token will be obtained for the user's Google account. This access token can then be used to authenticate requests to the Google Cloud Vision API on behalf of the user.

So, in summary, the user does not need to indicate that the Google account they created the project on is the account they need to log in with. They simply need to sign in to their Google account and grant permission to the extension to access the Google Cloud Vision API using the OAuth2 scopes that were specified in the `manifest.json` file.  This way, if you want your friends and family to use your Google Cloud Vision Service, you can simply share the extension with them and they can sign in to their own Google account and grant permission to the extension to access the Google Cloud Vision API using the OAuth2 scopes that were specified in the `manifest.json` file, and will be charged only to single Google Cloud Vision API project (billings).  You can also set the ceiling amount in the billings so that if your friends and family abuse the service, you will not be charged more than the ceiling amount.



## Usage

To use this extension, simply install it in your Chrome or Microsoft Edge browser and navigate to a page with manga images. The extension will automatically detect the manga images and display furigana in real-time.




## License

This extension is licensed under the MIT License. See the `LICENSE` file for details.