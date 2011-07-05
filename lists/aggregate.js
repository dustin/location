function(head, req) {
    log(req.query);
    var minWidth = req.query['minWidth'] || 0.01,
        minHeight = req.query['minHeight'] || 0.01;

    var row;
    // Since there's no key ordering, I have to collect them all and sort them first.
    var readings = [];
    while ((row = getRow())) {
        readings.push([Math.floor(parseInt(row.id) / 1000), row.value,
                       row.geometry.coordinates[0],
                       row.geometry.coordinates[1]]);
    }
    readings.sort(function(a, b) {
        return b[0] - a[0];
    });

    // Now I can perform useful aggregation

    function tooFar(lng1, lat1, lng2, lat2) {
        return Math.abs(lng1 - lng2) > minWidth
            || Math.abs(lat1 - lat2) > minHeight;
    }

    var output = [];
    var prev = 0, prevlbl = '', prevlng = 0, prevlat = 0;
    readings.forEach(function(r) {
        var ts = r[0], label = r[1], lng = r[2], lat = r[3];

        if (ts - prev > 7200 || tooFar(lng, lat, prevlng, prevlat)) {
            if (output.length > 0) {
                output[output.length - 1][1] = label + " - " + prevlbl;
            }
            output.push([ts, label, lng, lat]);
        }

        prev = ts;
        prevlng = lng;
        prevlat = lat;
        prevlbl = label;
    });

    send(JSON.stringify({rows: output}));
}