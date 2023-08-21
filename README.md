# Hotel reservation backend

# Project environment variables

```
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=mysupersecretNOBODYKNOWS
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017
```

# Project outline

- users -> book room from an hotel
- admins -> going to check reservation/bookings
- Authentication and authorizations -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database managements -> seeding, migration

#

## Database => mongodb

# Docs

https://mongodb.com/docs/drivers/go/current/quick-start

# Go client

go get go.mongodb.org/mongo-driver/mongo

##

## Web framework => fiber

# Docs

https://gofiber.io

# Install gofiber

go get github.com/gofiber/fiber/v2

#

## Docker

docker run --name mongodb -d mongo:latest -p 27017:27017
