# CrypTML

This is a simple (and naive) little project to allow wrapping encrypted HTML into
a little decryption wrapper, that will do the decryption client side
using browser native decryption capabilities.

## Motivation

I sometimes produce markdown documents that I want to share in a rendered
form but that contain sensitive information. Most filesharing systems (dropbox,
nextcloud, seafile, google drive, etc.) will only allow to download the
attached files but not show the HTML as-is. To keep the content "safe" from
unintended visitors and not having to trust the server I share it with, I can
now encrypt it but still produce a link that can be viewed immediately.

The security of course relies only on the fact that the key is not practically
guessable. As soon as the (full) URL leaks, the security is gone.

## Example

* Create the HTML file you want to share.
* Call `cryptml source/myfile.html target/myfile.html`
* Note the printed key.
* Upload `target/myfile.html` to the HTTP server of your choice.
* Go to your browser and visit: `https://<yourserver>/myfile.html#<key>`
  (replacing `<key>` with the one you noted earlier).