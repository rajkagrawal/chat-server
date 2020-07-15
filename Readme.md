#Application workings :

When a user(not connected in room) sends a message all the connected individual client(not joined the room) will receive the messages.<br/>
When a user(connected in room) sends a message all the client connected(joined) in a room will receive the messages.<br/>
For http api, since the user doesnt receive streaming messages, so when user uses Fetch api endpoint, it is returned with all the 
messages sent/received. The client should be responsible to filter out the messages marked for room/individula chat message. 




#curl request examples

Request body use same template for all the type of request which are user_id,room_id,message


when a person wants to post a message : <br/>
curl "http://localhost:8080/post" -d '{"user_id":"username","room_id":"", "message":"message to send"}'

when user joins the room room_id is must <br/>
curl "http://localhost:8080/post" -d '{"user_id":"username","room_id":"roomname", "message":"this message will go in the room"}'

when a user wants to get list of rooms created : 
curl "http://localhost:8080/command" -d '{"user_id":"username","room_id":"roomname", "message":"\\rooms"}'

when a user wants to query the message : 
curl "http://localhost:8080/fetch" -d '{"user_id":"username","room_id":"roomname", "message":""}'

#run the chat server
Go to the root dir i.e ~/chat-server <br/>
Make necessary changes in config file i.e config.json <br/>
~/chat-server$  go run cmd/main.go -configfile=./config.json
 
#telnet commands examples 
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

