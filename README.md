# Highload API written in Golang

## Docs
Documentation available on `/swagger/index.html` route

## Frameworks
- Fiber
- Gorm
- Testify
- swag (provide swagger)

## Description
Collect data from external sensors by http POST request and return arithmetic mean for some time interval.
Sensors are in rooms, rooms are in zones. Only owner and admin/superuser has access to them.

## Usage
- Register by POST `/auth/register` route
- Create zones by POST `/zone`
- Create rooms by POST `/room`
- Create sensors by POST `/sensor`
- Sensors send data into `/external/sensors_data` POST route in real-time
