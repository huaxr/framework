#!/usr/bin/python

import requests

url="http://127.0.0.1:8888/test/11"

while 1:
    r = requests.get(url)
    state=json.loads(r.text).get("code")
    if state == 1:
        break
    print("query", state)