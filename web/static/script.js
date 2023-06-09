//create map object and set default positions and zoom level
var map = L.map('map').setView([20, 0], 2);
L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'}).addTo(map);

function logJSONData() {
    var mapData = {{.mapData}};
    console.log(data.telemetry);
    data.telemetry.forEach(tel => {
        L.marker([tel.latitude, tel.longitude]).bindPopup(
            `<h3>Deployment details</h3>
            <p>IP Address:\t${tel.ip_address}</p>
            <p>version:\t${tel.mainflux_version}</p>
            <p>last seen:\t${tel.last_seen}</p>
            <p>country:\t${tel.country}</p>
            <p>city:\t${tel.city}</p>`
        ).addTo(map);
    });
}
logJSONData();