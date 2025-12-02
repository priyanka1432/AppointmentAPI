Here is the API documentation formatted as a clean Markdown file. You can copy the code block below and save it as API\_DOCS.md.

\# AppointmentAPI Documentation



\*\*Base URL\*\*: `http://localhost:8080/api`



\## Overview

This API allows coaches to set their availability schedules and users to view slots, book appointments, and manage their bookings.



---



\## Endpoints



\### 1. Set Coach Availability

Define the working hours for a coach on a specific day of the week.



\- \*\*URL\*\*: `/coaches/availability`

\- \*\*Method\*\*: `POST`

\- \*\*Content-Type\*\*: `application/json`



\*\*Request Body:\*\*



| Field | Type | Description |

| :--- | :--- | :--- |

| `coach\_id` | Integer | The unique ID of the coach. |

| `day` | String | Day of the week (e.g., "Monday", "Tuesday"). Case-insensitive. |

| `start\_time` | String | Start time in `HH:MM` (24-hour) format. |

| `end\_time` | String | End time in `HH:MM` (24-hour) format. |



\*\*Example Request:\*\*

```json

{

&nbsp; "coach\_id": 1,

&nbsp; "day": "Monday",

&nbsp; "start\_time": "09:00",

&nbsp; "end\_time": "12:00"

}



Responses:

&nbsp;\* 201 Created: Availability successfully saved.

&nbsp;\* 400 Bad Request: Validation error (e.g., start time is after end time).

&nbsp;\* 500 Internal Server Error: Database error.

2\. Get Available Slots

Fetch all available 30-minute booking slots for a specific coach on a specific date. The system automatically calculates slots based on the coach's schedule and removes currently booked slots.

&nbsp;\* URL: /users/slots

&nbsp;\* Method: GET

Query Parameters:

| Parameter | Type | Required | Description |

|---|---|---|---|

| coach\_id | Integer | Yes | The ID of the coach. |

| date | String | Yes | The date to check in YYYY-MM-DD format. |

Example Request:

GET /users/slots?coach\_id=1\&date=2025-12-15

Success Response (200 OK):

{

&nbsp; "slots": \[

&nbsp;   "2025-12-15T09:00:00Z",

&nbsp;   "2025-12-15T09:30:00Z",

&nbsp;   "2025-12-15T10:00:00Z"

&nbsp; ]

}



Note: Times are returned in RFC3339 (UTC) format.

3\. Book an Appointment

Book a specific 30-minute slot with a coach.

&nbsp;\* URL: /users/bookings

&nbsp;\* Method: POST

&nbsp;\* Content-Type: application/json

Request Body:

| Field | Type | Description |

|---|---|---|

| user\_id | Integer | The ID of the user making the booking. |

| coach\_id | Integer | The ID of the coach. |

| datetime | String | The exact timestamp of the slot in RFC3339 format. |

Example Request:

{

&nbsp; "user\_id": 1,

&nbsp; "coach\_id": 1,

&nbsp; "datetime": "2025-12-15T09:00:00Z"

}



Responses:

&nbsp;\* 201 Created: Booking successful.

&nbsp;\* 400 Bad Request: Invalid date format or slot not within coach's availability.

&nbsp;\* 409 Conflict: Slot is already booked (prevents double-booking).

4\. View User Bookings

Retrieve a list of all upcoming appointments for a specific user.

&nbsp;\* URL: /users/bookings

&nbsp;\* Method: GET

Query Parameters:

| Parameter | Type | Required | Description |

|---|---|---|---|

| user\_id | Integer | Yes | The ID of the user. |

Example Request:

GET /users/bookings?user\_id=1

Success Response (200 OK):

\[

&nbsp; {

&nbsp;   "id": 5,

&nbsp;   "user\_id": 1,

&nbsp;   "coach\_id": 1,

&nbsp;   "start\_time": "2025-12-15T09:00:00Z",

&nbsp;   "created\_at": "2025-12-01T10:00:00Z"

&nbsp; }

]



5\. Cancel Booking

Cancel an existing appointment. This endpoint validates that the booking belongs to the specified user before deletion.

&nbsp;\* URL: /users/bookings/:userid/:id

&nbsp;\* Method: DELETE

Path Parameters:

| Parameter | Type | Description |

|---|---|---|

| userid | Integer | The ID of the user who owns the booking. |

| id | Integer | The ID of the booking to cancel. |

Example Request:

DELETE /users/bookings/1/5

Responses:

&nbsp;\* 200 OK: Booking successfully cancelled.

&nbsp;\* 500 Internal Server Error: Failed to delete from database.

<!-- end list -->





