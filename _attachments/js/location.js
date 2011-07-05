var currentMarkers = { };

function doStuff(app, map) {
    var bounds = map.getBounds();
    var ne = bounds.getNorthEast(), sw = bounds.getSouthWest();
    var bbox = [sw.lng(), sw.lat(), ne.lng(), ne.lat()];

    console.log("Current bounds: ", bbox);
    if (!bbox) {
        return;
    }

    var pixWidth = $('#map').width();
    var pixHeight = $('#map').height();
    var minMove = 5; // Must move at least five pixels
    console.log("Pixel bounds:",pixWidth, pixHeight);
    var mh = Math.abs(bbox[1] - bbox[3]) / (pixHeight / minMove),
        mw = Math.abs(bbox[0] - bbox[2]) / (pixWidth / minMove);

    app.spatial_list('aggregate', 'points',
                     {bbox: bbox.join(','), minHeight: mh, minWidth: mw,
                      success: function(res) {
                          var newItems = 0;
                          res.rows.forEach(function(r) {
                              if (!currentMarkers[r[0]]) {
                                  ++newItems;
                                  currentMarkers[r[0]] = new google.maps.Marker({
                                      position: new google.maps.LatLng(r[3],
                                                                       r[2]),
                                      map: map,
                                      title: "At " + r[1]
                                  });
                              }
                          });
                          console.log("Planted " + newItems + " markers.");
                      }});
}

function clearMarkers(map) {
    if (currentMarkers) {
        for (i in currentMarkers) {
            currentMarkers[i].setMap(null);
        }
        currentMarkers = { };
    }
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

    function redraw() {
        if (timer) {
            clearTimeout(timer);
        }
        timer = setTimeout(function() {
            doStuff(app, map);
        }, 100);
    }

    google.maps.event.addListener(map, 'bounds_changed', redraw);
    google.maps.event.addListener(map, 'zoom_changed', function() {
        clearMarkers(map);
        redraw();
    });
}
