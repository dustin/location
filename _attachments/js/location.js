function map_init(app) {
    var latlng = new google.maps.LatLng(37, -122);
    var myOptions = {
        zoom: 8,
        center: latlng,
        mapTypeId: google.maps.MapTypeId.ROADMAP
    };
    var map = new google.maps.Map(document.getElementById("map"),
                                  myOptions);
}
