# Momenton Developer Challenge

 - `make build`
 - `make run`
 - `curl -X POST  http://localhost:3333/login --header 'Content-Type: application/json' --data-raw '{"User": "admin","Password": "password"}'`
 - `curl http://localhost:3333/user --header 'Authorization: Bearer <ACCESS_TOKEN>'` (this is a private path, e.g. test with no/incorrect token)
 - `curl -X POST  http://localhost:3333/refresh --header 'Content-Type: application/json' --data-raw '{"refresh_token": "<REFRESH_TOKEN>"}'`

## Notes

Password is set in the `.env` file. "password" for lack of a better idea. 

Accepts any user name, stores the user name and returns it when you GET the `/user` endpoint.