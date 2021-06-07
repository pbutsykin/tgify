# tgify
TGify - it's simple utility for redirect console command output to your telegram chat

# Example of usage:
```
$ ping 8.8.8.8 | tgify

$ tail -f /tmp/log | tgify

$ tgify --args dmesgdmesg -l emerg,alert,crit,err,war
```

##### To stop sending output just send to chat: "Stop" or "S" (case insensitive)


# Why is it needed?
Sometimes it can be convenient to watch the progress of a long-term task from telegram chat.

Also, the utility can simplify the development of programs that need to send human-readable notifications.
For example, if we want to develop a program that will search for big discounts in online stores. It will be
enough to develop the program that prints notifications to stdout.

# How to confgure?
For configuration you need to create config file **~/.tgify/config.yam** :
```
token:  777777777:AAh81ftRp7dsFGh8Q1UaEsksCH8 #your telegram bot token
chatIds: [
  1234567890
]  #chat id list for sending messages
```
