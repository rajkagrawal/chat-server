[![Build Status](https://travis-ci.org/rajaanova/chat-server.svg?branch=master)](https://travis-ci.org/rajaanova/chat-server)

# Application workings :

When a user(not connected in room) sends a message all the connected individual client(not joined the room) will receive the messages.<br/>
When a user(connected in room) sends a message all the client connected(joined) in a room will receive the messages.<br/>
For http api, since the user doesnt receive streaming messages, so when user uses Fetch api endpoint, it is returned with all the 
messages sent/received. The client should be responsible to filter out the messages marked for room/individula chat message. 



# curl request examples

Request body use same template for all the type of request which are user_id,room_id,message


when a person wants to post a message : <br/>
curl "http://localhost:8080/post" -d '{"user_id":"username","room_id":"", "message":"message to send"}'

when user joins the room room_id is must <br/>
curl "http://localhost:8080/post" -d '{"user_id":"username","room_id":"roomname", "message":"this message will go in the room"}'

when a user wants to get list of rooms created : 
curl "http://localhost:8080/command" -d '{"user_id":"username","room_id":"roomname", "message":"\\rooms"}'

when a user wants to query the message : 
curl "http://localhost:8080/fetch" -d '{"user_id":"username","room_id":"roomname", "message":""}'

# run the chat server without docker 

1. Go to the root dir i.e ~/chat-server <br/>
2. Make necessary changes in config file i.e config.json <br/>

_~/chat-server$_ go run cmd/main.go -configfile=./config.json

# run the chat server with docker


_~/chat-server$_ docker build . -t chatserver<br/>
_~/chat-server$_ docker run -p 5050:5050 -p 8080:8080 --name=rajchat --rm --mount type=bind,source="$(pwd)"/logs,target=/logs  chatserver
 
# telnet commands examples 
$ telnet 127.0.0.1 5050
<pre>
Trying 127.0.0.1...
Connected to 127.0.0.1.
Escape character is '^]'.
Please enter your username ::raj
some commands to use...
\rooms:                 list all the rooms/channels
\create <room_name>:    create a new room
\join <room_name>:      join a room
\exitroom:              removes users from the room
\help:                  prints all commands
\exit:                  go offline
\unsubscribe <user_id>:block the messages from user
start chatting now ...
time : 2020-07-15T08:52:06, senderID : ravi, msg : hey raj
ok ravi
time : 2020-07-15T08:52:17, senderID : raj, msg : ok ravi
time : 2020-07-15T08:52:30, senderID : htpuser, msg : hey raj
</pre>

# Approach taken while writing the application

1. Since this application is multiclient a loop is initiated to listen to the connection request.
2. For every connection run 2 go routines one to receive client input and another to send the messages on chat/terminal window.
3. While sending the messages, need to take care of messages which are like command(added switch case and regex to handle specific command) and another is normal messages that needs to be sent to connected client<br/>
4. For http api since the connection is not long live, to maintain the chats to be delivere a cache is maintained. Whenever a http client fetches this record client has to filter on their part the appropriate message to be displayed.<br/>

# Further developments
1. Rest api can be exposed to demarcate the message to be returned.
2. Few more functionality such as active users, \subscribe(unsubscribed user can be subscribed back) command can be added.
3. Cleanup tasks such as few goroutines are still lurking behind which might add to leakage, some policy regarding cache msg expriation can be taken care of.
