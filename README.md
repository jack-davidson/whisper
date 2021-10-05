# Whisper

## Database Documentation
Whisper uses postgresql for the database.

### Schema
**TODO**: document database schema

## API documentation

### `/newuser` Given a username and password, create a new user.
    Headers:
        Name     - username for new user
        Password - password for new user

    Example: curl -H 'Name: MyUserName' -H 'Password: Password123' http://localhost:3000/newuser

### `/message`: Given a name, password, and recipient, send a message.
    Headers:
        Name     - username of sender
        Password - password of sender
        For      - username of recipient of message

    Example: curl -H 'Name: MyUserName' -H 'Password: Password123' -H 'For: AnotherUserName' http://localhost:3000/message

### `/messages`: Given a username, get all of its messages (inbox).
    Headers:
        Name     - username of user
        Password - password of user

    Example: curl -H 'Name: MyUserName' -H 'Password: Password123' http://localhost:3000/messages

### `/user`: Given a username, lookup the user's id.
    Headers:
        Name - username of user
    
    Example: curl -H 'Name: MyUserName' http://localhost:3000/user

### `/deleteuser`: Given a username and a password, delete the user.
    Headers:
        Name     - username of user
        Password - password of user

    Example: curl -H 'Name: MyUserName' -H 'Password: Password123' http://localhost:3000/deleteuser
