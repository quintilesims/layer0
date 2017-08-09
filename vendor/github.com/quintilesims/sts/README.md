# System Test Service

## API

#### `GET /health`
Returns the current health of the service.

Response:
```
{
  "TimeCreated": "string",
  "Mode": "string"
}
```

#### `POST /health`
Update the current mode of the service. 
Valid options for `mode` are:
* `"normal"` - make the service operate normally
* `"slow"` - make the service wait 20 seconds before responding to future `GET /health` requests
* `"die"` - make the service exit its main process 5 seconds after receiving the request

Request:
```
{
  "Mode": "string"
}
```

Response:
```
{
  "TimeCreated": "string",
  "Mode": "string"
}
```


#### `GET /command`
Returns the commands ran by the service.

Response:
```
[
  {
    "Name": "string",
    "Args": ["string"],
    "Output": "string"
  },
  ...
]
```

#### `GET /command/:name`
Returns the commands with the specified name.

Response:
```
{
  "Name": "string",
  "Args": ["string"],
  "Output": "string"
}
```

#### `POST /command`
Executes the command(s) given by `Args`. 
The output can be retrieved `GET /commands/:name`

Request:
```
{
  "Name": "string",
  "Args": ["string"]
}
```

Response:
```
{
  "Name": "string",
  "Args": ["string"],
  "Output": "string"
}
```

## Update Image
Run `make release`
