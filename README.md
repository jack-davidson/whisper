# Whisper

## API documentation

### `/deleteuser`: Given a username and a password, delete the user.
    Headers:
        Name     - username of user
        Password - password of user

### `/message`: Given a name, password, and recipient, send a message.
    Headers:
        Name     - username of sender
        Password - password of sender
        For      - username of recipient of message

### `/messages`: Given a username, get all of its messages (inbox).
    Headers:
        Name     - username of user
        Password - password of user

### `/newuser` Given a username and password, create a new user.
    Headers:
        Name     - username for new user
        Password - password for new user

### `/user`: Given a username, lookup the user's id.
    Headers:
        Name - username of user
