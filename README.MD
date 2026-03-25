# Notifications manager system -THC

This is a basic notification management system for authenticated users that 
allows each user to manage and send notifications through different channels.

### Features

- Create a new user
- Get user data an token authorization / Login
- Update user data
- Create a new notification linked to a user
- Get a list of notifications
- Update a specific notification
- Delete a specific notification

### Author

Julián Imá Kherí Rentería

- Github: https://github.com/imakheri
- LinkedIn: https://www.linkedin.com/in/imakheri/
- Email: imakheri@gmail.com

### Badges

[![CircleCI](https://dl.circleci.com/status-badge/img/circleci/AoMkXQmCK1qb1Ns4CwHguW/DARMbqK4h4EKBtmTqw8kSv/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/circleci/AoMkXQmCK1qb1Ns4CwHguW/DARMbqK4h4EKBtmTqw8kSv/tree/main) [![Coverage Status](https://coveralls.io/repos/github/Imakheri/notification-system-thch/badge.svg?branch=main)](https://coveralls.io/github/Imakheri/notification-system-thch?branch=main)

### Deployed app running on Heroku

- [API Swagger](https://imakheri-notification-system.herokuapp.com/swagger/index.html)

### Table of contents

1. [Features](#features)
2. [Author](#author)
3. [Badges](#badges)
4. [Deployed app running on Heroku](#deployed-app-running-on-heroku)
5. [Technologies](#technologies)
6. [Routes](#routes)
7. [Pre-requisites](#pre-requisites)
8. [How to run app](#how-to-run-app)
9. [Authentication & Endpoints](#authentication--endpoints)
9. [How to run tests](#how-to-run-tests)
10. [Areas to improve](#areas-to-improve)
11. [Decisions made](#decisions-made)
12. [Development Tools](#development-tools)
    1. [Database migrations script](#database-migrations-script)
    2. [Environment variables](#environment-variables)
    3. [Mocks](#mocks)
13. [Deployment](#deployment)

### Technologies

- Programming language: Golang v1.25.0
- API Framework: Gin-gonic v1.12.0
- Database: MySQL v8.0
- Object-relational mapping framework: GORM v1.31.1
- Mocking framework: Gomock v0.6.0
- Containers: Docker, Docker-compose

### Routes

- Local: [API Swagger](https://localhost:8080/swagger/index.html)

### Pre-requisites

- Docker and Docker-compose installed
- No services running on localhost port 8080 and 3306

### How to run app

1. Just execute the next script to initialize portable (dockerized) version

````
chmod 711 ./up_dev.sh
./up_dev.sh
````
2. Go to swagger and test the app or consume API through curl or postman

3. Press Control + C to stop the app.

### Authentication & Endpoints

The application provides **7 endpoints** categorized by their access requirements:

Public Access: The `Create User` and `Login` endpoints require an **API Key** managed by the system administrator. This must be included in the request header as:
````
api_key: <your_admin_key>
````
User Authenticated Access: The remaining endpoints are protected and require a JSON Web Token (JWT), which is generated upon a successful login.

Note: All protected requests must include the token in the Authorization header using the Bearer scheme:
````
Authorization: Bearer <your_jwt_token>
````


### How to run tests

1. Execute the next script to run all application tests
````
chmod 711 ./up_test.sh
./up_test.sh
````

2. Press Control + C to stop the app

### Areas to improve

- The implementation IoC (Inversion of Control) pattern could make more readable main API file
- ORM uses auto migrations. Non-automated migrations would be the best option

### Decisions made

- Hexagonal Architecture: to be able to handle further changes in a proper way.
- Strategy Pattern: the problem statement specifies that the notification must be able to be sent through different channels.
  Due to the nature of each channel, there are specific conditions that must be met, therefore, they have different logic.
  The ability to define specific algorithms and keep them separate makes them interchangeable, allowing to determine which one will be used during execution.  
  This is why the strategy design pattern was implemented.
- GORM: it is the most popular ORM in Golang, for this reason it is easy to find fixes and people that know how to use it.
- Docker: makes the application portable. This allows to launch the application or run the tests with a single command, simplifying the process for anyone who wants to test the application in a virtualized environment.
- Testing: unit tests were primarily implemented in the domain layer (entities, handlers, and use cases).
  Mocks were used to define responses from adjacent layers. Error cases were considered happy paths.
  For the repository and database layers, end-to-end integration tests were implemented, following only happy paths, as error cases were already covered in the unit tests.
  Additionally, a large test was added to evaluate the functions `create user`, `create notification`, and `retrieve all notifications for the user who created the notification`.
- Gomocks: in order to save time and be maintainable in the future, this dependency was used to mock up the layers.
  Within the industry it is a standard, since it is very rare for mocks to be built by hand, in addition to providing us with specific methods aimed entirely at testing.
  A file was created where the parts of the code that you want to mock are compiled so that they can be generated automatically using a script in the terminal.
- JSON Web Token: because it is an open standard (RFC 7519) for securely transmitting information between a client and a server. It is a modern industry standard for authorization because it is compact, self-contained, and ensures data integrity through digital signatures.
  By using JWT, the server remains stateless, eliminating the need to store session data in memory and significantly improving project scalability.
- Artificial Intelligence: not a single line of this project was generated by AI

### Development tools

#### Database migrations script

If the database is just created and you want to execute a Model AutoMigrate
````
go run cmd/db/main.go -migrate
````

If the migration is already done, the "channels" table needs to be seeded, therefore please execute:
````
go run cmd/db/main.go -seed
````

If you want to do those two steps at once, please execute:
````
go run cmd/db/main.go -all
````

#### Environment variables

They should be defined. To find an example of the values you can use .env.example

#### DISCLAIMER: It is important to clarify that sensitive information was hardcoded into the Docker-compose files to achieve portability and ease of deploying the project in a virtualized environment. However, different global variables will be defined for the production version.

#### Mocks

If you want to create the mocks to run the tests locally or contribute to the development of the application, please execute:

````
make mocks
````

### Deployment

This project uses CircleCI for remote builds and automated testing. Upon every push, the platform verifies the build integrity and ensures all tests pass.
Once validated, code coverage data is automatically transmitted to Coveralls for detailed analysis and reporting.

