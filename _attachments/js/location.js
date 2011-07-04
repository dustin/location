var currentMarkers = { };

function doStuff(app, map) {
    var bounds = map.getBounds();
    var ne = bounds.getNorthEast(), sw = bounds.getSouthWest();
    var bbox = [sw.lng(), sw.lat(), ne.lng(), ne.lat()];
    console.log("Current bounds", bbox);
    if (!bbox) {
        return;
    }
    app.spatial('points', {bbox: bbox.join(','),
                           success: function(res) {
                               var newItems = 0;
                               res.rows.forEach(function(r) {
                                   if (!currentMarkers[r.id]) {
                                       ++newItems;
                                       currentMarkers[r.id] = new google.maps.Marker({
                                           position: new google.maps.LatLng(r.geometry.coordinates[1],
                                                                            r.geometry.coordinates[0]),
                                           map: map,
                                           title: "At " + r.value
                                       });
                                   }
                               });
                               console.log("Planted " + newItems + " markers.");
                           }});
}

function map_init(app) {
    var timer = undefined;
    var latlng = new google.maps.LatLng(37, -122);
    var myOptions = {
        zoom: 8,
        center: latlng,
        mapTypeId: google.maps.MapTypeId.ROADMAP
    };
    var map = new google.maps.Map(document.getElementById("map"),
                                  myOptions);
    google.maps.event.addListener(map, 'bounds_changed', function() {
        if (timer) {
            clearTimeout(timer);
        }
        timer = setTimeout(function() {
            doStuff(app, map);
        }, 100);
    });
}
