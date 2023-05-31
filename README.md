<img src="https://wtop.com/wp-content/uploads/2019/07/AP_100073447735.jpg" width="600" height="400" />

# E.T. Phone Home

Server to receive and store information regarding mainflux deployments. This information includes:

- IP Address
- Mainflux Version
- Last Seen
- Mainflux Service

## Usage
To Run:

```bash
go run ./cmd/homing-server/main.go
```


### Requirements
- [IP to Location database](https://lite.ip2location.com/)
- [Google service account](https://developers.google.com/identity/protocols/oauth2/service-account) - Don't forget to share the google sheets document with service account email. This provides a .json file.

Spreadsheet ID and sheet ID can be found on the google sheets url: 
`https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>`