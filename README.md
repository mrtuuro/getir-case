## Getting Started

This project contains 3 endpoints. 1 endpoint let's us retrieve data from 
MongoDB. Other 2 endpoints uses in-memory database. One is to store and the other one is to retrieve data.

## Requirements

* Go 1.22+
* MongoDB 1.14

## Running Project

1. Install the project's source code.
2. Create `.env` file.
3. Insert your environment variables. For this project we only need MongoDB URI `MONGODB_URI` and a port `PORT` to run the project.
4. Mongo URI should be in this format: `"mongodb+srv://<username>:<password>@<clusterName>.mongodb.net/?retryWrites=true&w=majority"`
5. While in project directory run the command `make run`

## API ENDPOINTS

---

### ```POST /api/records```

This endpoint is used to retrieve data from MongoDB by given request body.

### Request Body
```
{
    "startDate": "2016-01-26",
    "endDate": "2018-02-02",
    "minCount": 2700,
    "maxCount": 3000
}
```

* `startDate` field must be `string` type but in “YYYY-MM-DD” format.
* `endDate` field must be `string` type but in "“YYYY-MM-DD”" format.
* `minCount` field must be `int` type.
* `maxCount` field must be `int` type.

### Response
```
{
    "code": 0,
    "msg": "Success",
    "records": [
        {
            "key": "TAKwGc6Jr4i8Z487",
            "createdAt": "2017-01-28T01:22:14.398Z",
            "totalCount": 2800
        },
        {
            "key": "NAeQ8eX7e5TEg7oH",
            "createdAt": "2017-01-27T08:19:14.135Z",
            "totalCount": 2900
        }
    ]
}
```
---
### ```POST /api/in-memory```

This endpoint is used to store data in an in-memory database.
### Request Body
```
{
    "key": "key1",
    "value": "value1"
}
```

* `key` field must be `string` type.
* `value` field must be `string` type.
### Response
```
{
    "key": "key1",
    "value": "value1"
}
```

---

### ```GET /api/memory```

This endpoint is used to retrieve data from an in-memory database.
### Query Params
```
?key=key1
```

### Response
```
{
    "key": "key1",
    "value": "value1"
}
```
---

