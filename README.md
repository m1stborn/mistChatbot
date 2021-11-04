# mistbot

A line chatbot that send notification to you through LINE notify when a streamer go online.

## Line Usage

1. Add this account on Line : https://lin.ee/xeaZqKh
2. Connect to Line notify by clicking the link mistbot provide to you after you add mistbot to your friend.
3. Now you can start and subscribe some channel !

## Basic bot command

All the chatbot command starting with "/" :
1. `/sub [twitch ID]` : subscribe a channel
   * example: /sub never_loses
2. `/del [twitch ID]` : delete a channel
   * example: /del qq7925168
3. `/list` : list all the channel that you subscribe

## Demo
| Subscribe channel  | Delete channel  | Notify message   |
|---|---|---|
|![image](./assets/sub_command.jpg) |![image](./assets/del_command.jpg)|![image](assets/notify_message.jpg)   |

## TODO
- [ ] Automatically renew Twitch API token.
- [ ] Integrate Youtube live stream.
