The application business :

When a user(not connected in room) sends a message all the connected individual client(not joined the room) will receive the messages.
When a user(connected in room) sends a message all the client connected(joined) in a room will receive the messages.




#curl request examples

Request body use same template for all the type of request which are user_id,room_id,message


when a person wants to post a message
curl "http://localhost:8080/post" -d '{"user_id":"username","room_id":"", "message":"message to send"}'

when user joins the room room_id is must
curl "http://localhost:8080/post" -d '{"user_id":"username","room_id":"roomname", "message":"this message will go in the room"}'

when a user wants to get list of rooms created
curl "http://localhost:8080/command" -d '{"user_id":"username","room_id":"roomname", "message":"\\rooms"}'

when a user wants to query the message 
curl "http://localhost:8080/fetch" -d '{"user_id":"username","room_id":"roomname", "message":""}'


#telnet commands examples 
:~/Desktop/poc$ telnet 127.0.0.1 5050
Trying 127.0.0.1...
Connected to 127.0.0.1.
Escape character is '^]'.
Please enter your username :: dummyusername

\create roomname       #when this command is run a room is create
time : 2020-07-13T16:46:20, senderID : system, msg : room created rommname
\join roomname         #when this command is run user joins the room so that messages in that room can be sent/received


Type \help on console to get the list of available options 
