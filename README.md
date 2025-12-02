AppointmentAPI

A small Go (Gin + GORM) REST API to manage coach availabilities and user bookings (appointments).

## Overview

This project exposes endpoints to:

 * Add coach availability windows (day + start/end time)

 * List available 30-minute slots for a coach on a given date

 * Book a 30-minute slot for a user

 * List a user's upcoming bookings

 * Cancel a booking



The API uses Gin for HTTP routing and middleware, and GORM with a MySQL driver for persistence.

## Design choices & assumptions

Design choices


* Gin is used for lightweight, performant HTTP handling and middleware.

* GORM provides easy ORM features and automatic migrations (used in db.InitDb()).

* MySQLis the default database (see db.InitDb() which uses a MySQL DSN). The project expects appointmentdb by default.

* Business logic is separated into repo, service and controller layers.

* Availability windows are stored per coach as a day name + `HH:MM` start/end strings. Available 30-minute slots are computed at runtime.

* Bookings use a DB-level unique index (coach + start\_time) to prevent double-booking; the service maps DB duplicate-key errors to a conflict response.

Assumptions



* The code shown is the source of truth. No additional hidden configuration files are required.

* Default HTTP port: :8080 (set in main.go).

* Default MySQL DSN in db.InitDb() is: `root:root@tcp(127.0.0.1:3306)/appointmentdb?charset=utf8mb4\&parseTime=True&loc=Local` — update this for your environment.

* The application auto-migrates the schemas for `User`, `Coach`, `Availability`, and `Booking` on startup.



## Requirements



 * Go 1.20+ (or the version in `go.mod`)

 * MySQL (server) or Docker (to run MySQL container)

 * `git` (if cloning/pushing)



---

## Setup & Run Locally (Simple Steps)

1. Clone the repository

git clone https://github.com/priyanka1432/AppointmentAPI.git

cd AppointmentAPI


2. Download Go dependencies

go mod tidy


3. Start MySQL

Make sure MySQL is running on your system and a database named appointmentdb exists.
Default credentials in the code:

* username: root

* password: root

* host: 127.0.0.1:3306

You can create the DB using:

CREATE DATABASE appointmentdb;

4. Run the project

go run main.go

If everything is correct, you will see a message like:

database connected and migrated

and the server will start at: [http://localhost:8080](http://localhost:8080)

## API endpoints

Base path: `/api`

1. Add coach availability

POST `/api/coaches/availability`
Request body (JSON)

{

 "coach_id": 1,

 "day": "Monday",

 "start_time": "09:00",

"end_time": "12:00"

}

* day must be one of: `Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday` (case-insensitive)

* `start_time` and `end_time` must be `HH:MM` (24h)

* Window must be at least 30 minutes and `start_time < end_time`

Responses:

* `201 Created` with the created availability JSON on success

* `400 Bad Request` for validation errors

* `500 Internal Server Error` for DB errors

2. Get available slots for a coach on a date



GET `/api/users/slots?coach_id=<coach_id>\&date=YYYY-MM-DD`

* `date` must be in `YYYY-MM-DD` format
* Response format:

{ "slots": \["2025-12-10T09:00:00Z", "2025-12-10T09:30:00Z", ...] }

Slots are returned in RFC3339 (UTC) format. The server computes 30-minute slots within availability windows.

3\. Book a slot

POST `/api/users/bookings`

Request body (JSON):

{

  "user_id": 1,

  "coach_id": 1,

  "datetime": "2025-12-10T09:00:00Z"

}





 * `datetime` must be RFC3339 and align to :00 or :30 minutes (e.g., `09:00` or `09:30` UTC).

 * If the slot is already booked, the API returns `409 Conflict` with `slot already booked`.



Responses:



 * `201 Created` with the booking object on success

 * `400 Bad Request` for invalid payload or if slot is not within availability

 * `409 Conflict` if DB unique constraint is violated (slot already booked)



4. Get a user's upcoming bookings



GET `/api/users/bookings?user_id=<user_id>`



Response: JSON array of booking objects (ordered ascending by start time). Only bookings with `start\_time >= now` are returned.



5\. Cancel a booking



 * *DELETE * * `/api/users/bookings/:userid/:id`



 * Both `userid` and `id` (booking id) must match; the service validates that the booking belongs to the user before deleting.

 * On success: `200 OK` with `{ "status": "cancelled", "id": <id> }`



---

 Database schema (GORM models)



 * `User` — `{ id, name }`

 * `Coach` — `{ id, name }`

 * `Availability` — `{ id, coach_id, day, start\_time, end\_time }`

 * `Booking` — `{ id, user_id, coach_id, start\_time, created\_at }` with a unique index on `(coach_id, start\_time)` to prevent double bookings


After creating the appointmentdb database, you can add sample users and coaches using the following SQL queries in MySQL Workbench:

INSERT INTO users (name) VALUES 
('Priyanka'),
('Rahul'),
('Sneha'),
('Amit');


INSERT INTO coaches (name) VALUES
('Coach A'),
('Coach B'),
('Coach C');

---

## Examples (Postman)



Below are simple examples you can directly use in Postman.



 1. Add Coach Availability



Method: POST

URL:`http://localhost:8080/api/coaches/availability`

Body (JSON):

{

  "coach_id": 1,

  "day": "Monday",

  "start_time": "09:00",

  "end_time": "12:00"

}



2. Get Available Slots

Method: GET

URL:http://localhost:8080/api/users/slots?coach_id=1&date=2025-12-15

3. Book a Slot

Method: POST

URL:`http://localhost:8080/api/users/bookings`

Body (JSON):





{

  "user_id": 1,

  "coach_id": 1,

  "datetime": "2025-12-15T09:00:00Z"

}

4.Get User Bookings



Method: GET

URL:http://localhost:8080/api/users/bookings?user_id=1



5. Cancel a Booking



Method: DELETE

URL: http://localhost:8080/api/users/bookings/1/5


## Error handling
 * The project centralizes error handling in `internal/middleware.ErrorHandler()` which maps `internal/errors.AppError` to structured JSON responses with HTTP status codes.

 * Logger middleware logs each request method and path after processing.

FOLDER STRUCTURE
```
APPOINTMENTAPI
├── internal                           # All core application logic
│   ├── controller                     # Route handlers / API controllers
│   │   └── controller.go
│   ├── db                             # Database connection & setup
│   │   └── db.go
│   ├── errors                         # Custom error types & responses
│   │   └── errors.go
│   ├── middleware                     # Middlewares (logger, error handler)
│   │   ├── errorhandler.go
│   │   └── logger.go
│   ├── models                         # Data models / structs
│   │   └── models.go
│   ├── repo                           # Repository interfaces & implementations
│   │   ├── repo.go
│   │   └── repo_impl.go
│   ├── service                        # Business logic layer
│   │   └── service.go
│   └── utils                          # Utility/helper functions
│       └── utils.go
│
├── APIDOCS.md                         # API endpoint documentation
├── go.mod                             # Go module definition
├── go.sum                             # Dependency checksum file
├── main.go                            # Entry point of the application
└── README.md                          # Project documentation
```


