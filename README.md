# Tasks Management API Assessment

Go Developer assessment for creative advanced technologies.

## Explanation

This is the solution to the task management api assessment.

To avoid unnecessary complexity and overengineering i decided to:

* Use just logs rather than implementing tracing with open telemetry.
* Maintain a simple folder structure.
* Write less tests for duplicated use cases (mainly because it's an assessment)


## Testing

I did not test all parts of the codebase to avoid spending too much time on the project, since this is just and assessment project.

To show how i test my codebase i have a sample full coverage test for each use case.

* [Service Layer](components/users/service_test.go)
* [Http Handlers](http-handlers/users_test.go)


## Database

This API uses mongodb as the primary database.


## How to execute / use

* Using docker **(recommended)** run `docker-compose up` and connect to `localhost:5555`
* Using local mongodb and go installation, run `go run main.go` then adjust environment variables in `.env` to fit your current setup.


## Endpoints

##### Create new user

POST:  `/users/`

Sample Payload:

```json
{
   "firstName": "Wisdom",
   "lastName": "Matthew",
   "email": "talk2wisdommatt@gmail.com"
}
```


---

##### Get User

GET: `/users/{userId}`

---

##### Get Users

GET: `/users/?lastId=&limit=20`

`lastId` and `limit` url parameters are used for pagination.

---

##### Delete User

DELETE: `/users/{userId}`

---

##### Create Task

POST: `/tasks/`

Sample Payload:

```json
{
    "title": "Run 20 minutes",
    "startTime": "2022-02-18T11:01:00.000+00:00",
    "endTime": "2022-02-18T12:00:00.000+00:00",
    "userId": "6212c3112e46aabc11bbee1c",
    "reminderPeriod": "2022-02-18T12:00:00.000+00:00"
}
```

---

##### Get Task

GET: `/tasks/{taskId}`

---

##### Get Tasks

GET: `/users/{userId}/tasks?lastId=&pagination=20`

---

##### Delete Task

DELETE: `/tasks/{taskId}`
