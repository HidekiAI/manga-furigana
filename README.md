# manga-furigana
OCR-based real-time dynamic furigana for manga plugin to Chromium/MSEdge


## Setup (Goole Cloud Vision API), build, and attach to Chromium/MSEdge
1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Select your project from the dropdown menu at the top of the page.
3. Click on the navigation menu icon in the top-left corner of the page and select "APIs & Services" > "Credentials".
4. Click on the "Create credentials" button and select "Service account key".
5. Select "New service account" and enter a name for the service account.
6. You'll note a "KEYS" tab at the top center, in which you'd click that tab and select "ADD KEY" > "Create new key" 
7. Select "JSON" as the key type and click on the "Create" button.
8. Dialog box with "Private key saved to your computer" will appear, in which you should most likely be prompted to download the JSON file, name the JSON file to `credentials.placeholder.json` and place it in the `manga-furigana` (src) directory.
   a. TODO: In future, switch over to "Workload Identity Federation" to avoid having to download the JSON file.
   b. Go to "APIs & Services" > "Library" and search for "Cloud Vision API", click on the result and click on the "Enable" button.
   c. Probably a good idea to go to "APIs & Services" > "Dashboard" and click on the "Enable APIs and Services" button and enable "Compute Engine API" as well.
9. Run the "build.sh" bash script to build the project.
10. Copy the built folder (/tmp/build/manga-furigana) to some persistent directory (e.g. /opt/manga-furigana).
11. Assuming Chrome/MSEdge is set to allow Developer extensions, go to the extensions page (chrome://extensions/ or edge://extensions/), enable Developer mode, and click on "Load unpacked" and select the directory from step 11.
12. The extension should now be loaded and ready to use.


* These steps are wrong for the purpose of browser extensions, the above methods are  for the purpose of getting authentications between service-to-service (S2S) and security-wise, it's a very bad idea.  The correct way to do this is to integrate it via client-to-service (C2S) method in which is done via OAuth2.
* Revise the steps above to use OAuth2 instead of S2S.

