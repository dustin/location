#!/usr/bin/env python

import re
import sys
import json

import couchdb

import tripit

DB = couchdb.Server('http://127.0.0.1:5984/')['tripit']

idmap = {
    'WeatherObject': lambda x: x['trip_id'] + "_" + x['date'],
    'Profile': 'screen_name'
}

FLOAT_KEY_RE = re.compile(r'temp_c|wind_speed|precipitation|latitude|longitude')

BOOL_KEY_RE = re.compile('^is_')

def getId(doc):
    getter = idmap.get(doc['type'], 'id')
    if callable(getter):
        i = getter(doc)
    else:
        i = doc[getter]
    return doc['type'] + '_' + i

def cleanup(h):
    for k in h:
        if isinstance(h[k], dict):
            cleanup(h[k])
        elif isinstance(h[k], list):
            for e in h[k]:
                if isinstance(e, dict):
                    cleanup(e)
        elif FLOAT_KEY_RE.search(k):
            h[k] = float(h[k])
        elif BOOL_KEY_RE.search(k):
            h[k] = h[k] == 'true'

def main(argv):
    if len(argv) < 5:
        print "Usage: example.py api_url consumer_key consumer_secret authorized_token authorized_token_secret"
        return 1

    api_url = argv[0]
    consumer_key = argv[1]
    consumer_secret = argv[2]
    authorized_token = argv[3]
    authorized_token_secret = argv[4]

    oauth_credential = tripit.OAuthConsumerCredential(consumer_key,
                                                      consumer_secret,
                                                      authorized_token,
                                                      authorized_token_secret)
    t = tripit.TripIt(oauth_credential, api_url = api_url)
    ob = t.list_trip([('past', 'true'), ('include_objects', 'true')])
    docs = []
    for otype, values in ob.iteritems():
        if isinstance(values, list):
            print "Processing", otype
            for d in values:
                doc = {'type': otype}
                doc.update(d)
                doc['_id'] = getId(doc)
                cleanup(doc)
                docs.append(doc)

    DB.update(docs)

if __name__ == "__main__":
    main(sys.argv[1:])
