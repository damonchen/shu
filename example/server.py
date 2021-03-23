#!/usr/bin/env python
#coding=utf-8

import json
from flask import Flask, request


app = Flask(__name__)


@app.route('/api/v1/login', methods=["POST"])
def login():
    data = request.get_json()
    print(data)

    return json.dumps({'status': 'aaa'})


if __name__ == '__main__':
    app.run()
