function(doc) {
    emit({"type": "Point", "coordinates": [doc.longitude, doc.latitude]}, doc.ts);
}
