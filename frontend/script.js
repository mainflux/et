//create map object and set default positions and zoom level
var map = L.map('map').setView([20, 0], 2);
L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'}).addTo(map);
  
var Icon = L.icon({
    iconUrl: 'My project.png',
    //shadowUrl: 'http://leafletjs.com/docs/images/leaf-shadow.png',

    iconSize:     [20, 20], // size of the icon
    //shadowSize:   [50, 64], // size of the shadow
    iconAnchor:   [10, 20], // point of the icon which will correspond to marker's location
    //shadowAnchor: [4, 62],  // the same for the shadow
    popupAnchor:  [0, -10] // point from which the popup should open relative to the iconAnchor
});


async function logJSONData() {
    const res = await fetch("./public/sample.json")
    const data = await res.json()
    console.log(data.telemetry);
    data.telemetry.forEach(tel => {
        L.marker([tel.latitude, tel.longitude], {icon: Icon}).bindPopup(
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