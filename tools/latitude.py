#!/usr/bin/env python
# -*- coding: utf-8 -*-
#
# Started as a google sample latitude example.  Now grabs all latitude
# data and loads it into couchdb.
#
# Copyright 2010 Google Inc. All Rights Reserved.
# Copyright 2011 Dustin Sallings All Lefts Reversed

"""Simple command-line example for Latitude.

Command-line application that sets the users
current location.
"""

from apiclient.discovery import build

import time
import httplib2
import pickle
import signal

import couchdb

from apiclient.discovery import build
from apiclient.oauth import FlowThreeLegged
from apiclient.ext.authtools import run
from apiclient.ext.file import Storage

ISO8601 = "%Y-%m-%dT%H:%M:%S"

DB = couchdb.Server(os.getenv('COUCH_SERVER'))[os.getenv('COUCH_DB')]

# Uncomment to get detailed logging
# httplib2.debuglevel = 4

def ts(r):
  return time.ctime(float(int(r['timestampMs']) / 1000))

def store(records):
  for r in records:
    r['_id'] = r['timestampMs']
    r['ts'] = time.strftime(ISO8601, time.gmtime(int(r['timestampMs']) / 1000))

  DB.update(records)

def main():
  signal.alarm(47)
  storage = Storage('latitude.dat')
  credentials = storage.get()
  if credentials is None or credentials.invalid == True:
    auth_discovery = build("latitude", "v1").auth_discovery()
    flow = FlowThreeLegged(auth_discovery,
                           # You MUST have a consumer key and secret tied to a
                           # registered domain to use the latitude API.
                           #
                           # https://www.google.com/accounts/ManageDomains
                           consumer_key='west.spy.net',
                           consumer_secret='kF52E2QkcMYeWI5JBxUdqHkE',
                           user_agent='google-api-client-python-latitude/1.0',
                           domain='west.spy.net',
                           scope='https://www.googleapis.com/auth/latitude',
                           xoauth_displayname='Google API Latitude Example',
                           location='all',
                           granularity='best'
                           )

    credentials = run(flow, storage)

  http = httplib2.Http()
  http = credentials.authorize(http)

  service = build("latitude", "v1", http=http)

  more = True
  max_time = ""
  args = {'max_results': '1000', 'granularity': 'best'}

  # See if we're picking up from somewhere...
  try:
    lastKey = list(DB.view('_all_docs', limit=1, descending=True, startkey='9'))[0].id
    args['min_time'] = lastKey
    print "Resuming since", lastKey
  except IndexError:
    print "Not resuming."

  while more:
    r = service.location().list(**args).execute()['items']
    store(r)
    more = len(r) == 1000

    print "Got %d records from %s - %s" % (len(r), ts(r[-1]), ts(r[0]))
    args['max_time'] = r[-1]['timestampMs']

if __name__ == '__main__':
  main()
