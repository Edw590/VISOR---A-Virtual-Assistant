# V.I.S.O.R. - A Virtual Assistant
V.I.S.O.R., my in-development assistant, written in Go

## Notice
This project is a part of a bigger project, consisting of the following:
-   [V.I.S.O.R. - Android Version Assistant](https://github.com/DADi590/VISOR---A-better-Android-assistant)
-   [V.I.S.O.R. - A Virtual Assistant](https://github.com/Edw590/VISOR---A-Virtual-Assistant)

## Introduction
These are the desktop client and server versions of my in-development assistant, V.I.S.O.R. (the Android version is just above on the link).

The client one needs a bit more development to let it keep running on my computer like I let the Android one on my phone, but it's getting there.

The server one runs 24/7 on my Raspberry Pi (but it can run on Windows too. It's supported on both). This is supposed to be VISOR's "operations center", where all 24/7 things run. As an example, the RSS Feed Notifier and Email Sender modules. The notifier checks for news on the feeds and the sender sends the emails that the notifier queues on it. All always running.

## Questions
Feel free to create an Issue or a Discussion with any questions you have about this. I'm trying to make it as generic as possible for anyone to be able to use it, but I may forget to document things or something, so if you need anything, feel free to ask about it.

## Modules
| Number | Name | Client/Server | Description |
|-|-|-|-|
| N/A | Utils | Both | This is not a module, it's just a separate global package (all non-library modules are `main` packages; this one is the `Utils` package), but I'm writing about it because it can (and does) have utilities to communicate between modules. For example, it has utilities to queue emails to the Email Sender, so one just needs to call the function there to queue an email.
| 0 | **V.I.S.O.R.** | Both | These are the main VISOR programs (includes the server app and the client app in the same number), and are in this repository. Through here all communication are made with him.
| 1 | **Modules Manager** | Both | Manages all of VISOR's modules. It's responsible for keeping them running all the time and restarting them in case they stop for any reason.
| 2 | **S.M.A.R.T. Checker** | Server | Runs S.M.A.R.T. tests on the given disks and checks the S.M.A.R.T. information after the tests are done.
| 3 | **Speech** | Client (Windows only) | This is VISOR's speech module. It speaks or queues a notification about anything that needs to be spoken.
| 4 | **RSS Feed Notifier** | Server | Checks RSS feeds and queues an email about any news. Currently it's tested on YouTube channels *and playlists* (something YouTube didn't do nor does), and on StackExchange feeds. May work in others, but I didn't test (haven't needed so far).
| 5 | **Email Sender** | Server | Sends the emails that are queued for it to send. It works with cURL (the `curl` command), so it must be installed on the system and be on the PATH. It also works by sending an EML file containing the raw email information.
| 6 | **Online Information Checker** | Server | Checks the Internet for information like weather and news and updates a file with the information it got. This file can then be read by apps to get the information back, already ready for usage.
| 7 | **GPT Communicator** | Server | Sends and receives text to and from a local LLM (Large Language Model), like Llama 3 through the llama.cpp project.
| 8 | **Website Backend** | Server | It's the backend of VISOR's website. It is responsible for handling the requests from the frontend.
| 9 | **Reminders Reminder** | Client | Checks reminders and warns when one is triggered. The reminders are fetched from the server.
| 10 | **System Checker** | Client | Sends information about the system to the server, like Wi-Fi networks and Bluetooth devices in range. Or the state of the Wi-Fi adapter and the Bluetooth adapter. Or the screen brightness. Or others. The server can use this information to determine where the device is and if it's being used or not.
| 11 | **Speech Recognition** | Client | Currently only checks if the phrase "Hey VISOR" is spoken and shows the UI, but later should be used to detect normal speech to interact with VISOR.
| 12 | **User Locator** | Server | Locates the user based on everything the server knows about the user (the user must configure some things first) and based on all devices communicating with the server. For example, if the phone is communicating and the user is always with the phone, then the user is near the phone whether it's being used or not. With the computer, it must be being used because the user may leave the computer and go have lunch but not the phone (that must still be configured - the "AlwaysWith" device).

## Libraries
| Name | Description |
|-|-|
| **[Advanced Commands Detection](https://github.com/DADi590/Advanced-Commands-Detection)** | Detects commands in a sentence of words (a link because it's in another repository). It's the module that understands the user communication (voice or text - as long as it uses words). It can detect no so simple sentences of multiple commands, and understands the meaning of "don't", "it", and "and". Example of a complex sentence it can successfully understand (without the punctuation - it must not be present): `"turn it on. turn on the wifi, and and the airplane mode, get it it on. no, don't turn it on. turn off airplane mode and also the wifi, please."` (ignores/warns about the meaningless "it", turns on the Wi-Fi, and turns off the airplane mode and the Wi-Fi).
| **GPT** | Sends and gets text to and from the GPT Communicator module through the Website Backend module.
| **Online Information Getter** | Gets the information the Online Internet Checker module put on the website files.
| **Registry** | This is like the Windows Registry. Has keys and values with data inside. Useful to communicate inside the app between different modules for example. Or to show the state of various things to the user, since there's a screen with all the values and their data listed.
| **SpeechQueue** | This is an implementation of a speech queue to be used for the Speech modules of VISOR.
| **ULComm** | Makes the communication between the client and the server's User Locator module ("User Locator Communicator"). It's used to send device info from the clients to the server.

## Developer notes
This began as a Python project (in 2020), but even using an IDE got the project confusing. So I'm translating it to Go, because using Go solves the issues I was having (forces me to organize the code well enough to not have cyclic imports, since it won't compile if they exist, and forces the types on variables - ah, and I can mess with pointers, and that's nice).

### - To use the project
- Download this main project and the Advanced Commands Detection module (`git clone --recursive [repo link here]`). Then go to the `ClientCode` folder and run the command `go build -tags=client .`. For the server, go to the `ServerCode` folder and run the command `go build -tags=server .`. Finally move both files to the `bin` folder.
- Next go on that `bin` folder and edit the JSON file with your values and rename the file to PersonalConsts_EOG.json. VISOR needs an email of his own btw. Also needs a website (I use Nginx for it). I'll try to remove that requirement soon. But for full functionality (like communication between the app and the server) the website must exist.
- Go on each module folder and copy the JSON file to `data/UserData/MOD_[module number here]` (create the folders if they don't exist) and configure it (in case the module needs one).
- Start the client or the server executables and that's it.

#### Supported OSes
The entire project is supposed to be able to be run on Unix-like and Windows OSes (multi-platform project). If by chance any module is not supported on any operating system, it will refuse to run on the unsupported OS(es) - even though it can probably still be compiled for them (just not ran). In case there is a module like this, it will be warned on the Modules list above. This probably just means I haven't had the time or interest to program it for that OS and not because it really can't be run there.

To change it to run on Windows or Linux, just compile to the OS you want, put the binaries in the bin folder and configure the path on the PersonalConsts_EOG.json file. Nothing else needs to be done to change things from running on either OS.

## About
### - License
This project is licensed under Apache 2.0 License - http://www.apache.org/licenses/LICENSE-2.0.
