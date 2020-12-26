## F_WhatsappGo

 A small project rewritten in go which sends links to telegram from whatsapp chats and email.
 
 #### Requirements:
  1) Whatsapp account. 
  2) Email account.  //To-Do
 
 #### Steps for running:
  1) Rename `.config.json` to `config.json`
  2) Fill the `config.json` file
  3) Grab the latest binary from the releases or just `go build` it in root of this project. 
  4) Run the binary`./F_WhatsappGO` and scan the qr code which will appear on the terminal.   
 #### To-DO:
  1) Select individual chats. Right now it only works globally.
  2) Add email support
  3) <del>Needs docker to run because the wrapper library refuses to work properly without it.</del> Selenium can go to hell for all I care.

 #### Filtering text:
  1) You can set `filter_mode` to `blacklist` or `whitelist` to filter stuff accordingly to your need.
  2) Filling the links in`whitelist_links` will force the program to **not** ignore the messages containing the links whereas `blacklist_links` will force the program to ignore the same.

#### Legal:
This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by WhatsApp or any of its affiliates or subsidiaries. This is an independent and unofficial software. Use at your own risk.

##### SIde-Note about the rewrite :-
* This program uses [this](https://github.com/Rhymen/go-whatsapp) library (all credits to the repo owners and maintainers) which is based on a reverse engineered whatsapp's webui library.
* This enables anybody to use it without any dependency outside of this program unlike python which needs dockers and its own headache. Rewriting was fun and very educational at the same time.
