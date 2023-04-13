<img src="https://wtop.com/wp-content/uploads/2019/07/AP_100073447735.jpg" width="600" height="400" />

# E.T. Phone Home

Server to recive and store information regarding mainfflux deployments. This information includes:

- IP Address
- Mainflux Version
- Last Seen
- Mianflux Service

## Usage
To Run:

```bash
go run ./cmd/homing-server/main.go
\```


### Requirements
- [IP to Location database](https://lite.ip2location.com/)
- [Google service account](https://developers.google.com/identity/protocols/oauth2/service-account) - Don't forget to share the google sheets document with service account email.
- [Mainflux grpc authentication service](https://github.com/mainflux/mainflux/tree/master/auth)

example .env file:
```
MF_USERS_LOG_LEVEL="info"
MF_JAEGER_URL="localhost:6831"
MF_GCP_CRED="*.json"
MF_SPREADSHEET_ID=''
MF_SHEET_ID=0
MF_IP_DB="*.BIN"
MF_AUTH_GRPC_PORT=8181
MF_AUTH_GRPC_URL=localhost:8181
MF_AUTH_GRPC_TIMEOUT=1s
```

Spreadsheet ID and sheet ID can be found on the google sheets url: 
`https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>`