# Backend Notification APIs
Notification system API, which is capable of receiving a message and depending on the category of the message and the users subscribed to them, said users will be notified to the medium that they themselves chose.

#### Running make help will display many useful functions for running the API.

First, run Docker Compose to start the database. After that, run the API, and you can test the server with a request like this:

```bash
curl -X POST http://localhost:8080/notify \
-H "Content-Type: application/json" \
-d '{"category": "Sports", "message": "New Sports Event!"}'
```
