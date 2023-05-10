<img src="https://wtop.com/wp-content/uploads/2019/07/AP_100073447735.jpg" width="600" height="400" />

# E.T. Phone Home

[![build][ci-badge]][ci-url]
[![go report card][grc-badge]][grc-url]
[![coverage][cov-badge]][cov-url]
[![license][license]](LICENSE)

Server to receive and store information regarding mainflux deployments. This information includes:

- IP Address
- Mainflux Version
- Last Seen
- Mainflux Service

## Usage
To Run:

```bash
make docker-image-server
make docker-image-ui
make run
```


### Requirements
- [IP to Location database](https://lite.ip2location.com/)
