<img src="https://wtop.com/wp-content/uploads/2019/07/AP_100073447735.jpg" width="600" height="400" />

# E.T. Phone Home

![build][build]
![Go Report Card][grc]
[![License][LIC-BADGE]][LIC]

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


[grc]: https://goreportcard.com/badge/github.com/mainflux/callhome
[build]: https://github.com/mainflux/callhome/actions/workflows/ci.yml/badge.svg
[LIC]: LICENCE
[LIC-BADGE]: https://img.shields.io/badge/License-Apache_2.0-blue.svg
