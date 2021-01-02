[![codecov](https://codecov.io/gh/edusalguero/roteiro/branch/main/graph/badge.svg?token=DCDP4YT6T3)](https://codecov.io/gh/edusalguero/roteiro)
![CI](https://github.com/edusalguero/roteiro/workflows/run%20tests/badge.svg)
----
# Roteiro
Roteiro is  **r**ide-p**o**oling mobili**t**y v**e**hIcle **r**oute **o**ptimization

---

## Roteiro API
#### Version: v1

#### POST /problem

###### Summary:

Solve a problem with the given description

###### Description:

The endpoint can solve the vehicle routing problem stated in the request body synchronously.

###### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Success |
| 400 | Error |

#### POST /problem-long

###### Summary:

Queue a long-running problem with the given description

###### Description:

It will trigger a long-running background task. 
The GET solution endpoint could be used to ping to fetch the results. 
This is necessary for larger requests that take longer than 20 seconds.
When you POST to the above URL, it will return immediately with a 202 HTTP code and the problem_id

###### Responses

| Code | Description |
| ---- | ----------- |
| 202 | Problem queued |
| 400 | Error |

#### GET /solution/{problem_id}

###### Summary:

Get the solution for the given problem_id

###### Description:

It allows to recover the solution of a previously solved problem, or a previously queued problem. 
When the problem has not yet been resolved, the status code is 409 (Processing). 
Upon completion, the status code is 200. The response that you would normally get directly from the synchronous endpoint is now in the output.

###### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| problem_id | path | ID of related problem | Yes | string (uuid) |

###### Responses

| Code | Description |
| ---- | ----------- |
| 200 | Success |
| 400 | Error |
| 404 | Not found |
| 409 | Processing. The problem is not solved yet |
