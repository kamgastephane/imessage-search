# **LEARNING ON THE GO**

Simple API built in [GOLANG](https://golang.org) to track my travelling request for work

# HOW TO RUN

Launch the mongo instance using docker-compose
then run the go app 
``` 
$ docker-compose up
$ go run main.go
```


# API available

#### Get the list of trips
````
GET http://localhost:8081/trips
````

#### Get a trip by ID
````
GET http://localhost:8081/trip/{tripId}
````

#### Create a trip
````
POST http://localhost:8081/trip/
````
###### Example of a JSON payload
[
        {
            "city": "Douala",
            "country": "Camerun",
            "leavingBetween": {
                "from": "2019-06-20T10:00:00Z",
                "to": "2019-06-25T10:00:00Z"
            },
            "arrivingBetween": {
                "from": "2019-06-20T10:00:00Z",
                "to": "2019-06-25T10:00:00Z"
            },
            "notes": "Meeting with PO"
        }
    ]
#### Update a trip
````
PUT http://localhost:8081/trip/{tripId}
````
###### Example of a JSON payload
{
    "quotationId": "99f4f084-2456-4f9e-881a-785dd0c330d2",
    "trips": [
        {
            "city": "douala",
            "country": "Camerun",
            "leavingBetween": {
                "from": "2019-06-20T10:00:00Z",
                "to": "2019-06-25T10:00:00Z"
            },
            "arrivingBetween": {
                "from": "2019-06-20T10:00:00Z",
                "to": "2019-06-25T10:00:00Z"
            },
            "notes": "Meeting with PO"
        }
    ],
    "creationDate": "2019-05-11T18:01:47.49Z",
    "status": "CLOSED"
}

#### Update the status of a trip
````
PUT http://localhost:8081/trip/{tripId}/status/{state}
````
you can update the status of a trip
The possible available state are:
- NEW
- OPENED
- CLOSED

# TODO
Setup a Front-end using https://www.typeform.com and deploy on an EC2 instance