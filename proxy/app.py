import calendar
import logging
import os
import time
from flask import Flask, request, Response
from dotenv import load_dotenv
import requests

logging.basicConfig(filename="record.log", level=logging.INFO)
logging.getLogger("werkzeug").setLevel("WARNING")  # to filter out werkzeug logs

# configurations loaded from environment variables
load_dotenv()
PROXY_PORT = os.environ.get("PROXY_PORT")
WDB_RETRO_URL = os.environ.get("WDB_RETRO_URL")
WDB_RETRO_CLUSTER = os.environ.get("WDB_RETRO_CLUSTER")
WDB_RETRO_TOKEN = os.environ.get("WDB_RETRO_TOKEN")

app = Flask(__name__)

@app.route("/connect", methods=["POST", "GET", "DELETE", "PATCH"])
def connect():
    forward_url = "{0}/connect?cluster={1}&token={2}".format(WDB_RETRO_URL, WDB_RETRO_CLUSTER, WDB_RETRO_TOKEN)
    res = requests.request(
        method=request.method,
        url=forward_url,
        data=request.get_data(),
        allow_redirects=False,
    )

    log_event(method=request.method, url=request.url, forwarded_to=WDB_RETRO_URL)

    return Response(res.content, res.status_code)

@app.route('/get-started', methods = ['GET'])
def get_started():
    forward_url = "{0}/get-started".format(WDB_RETRO_URL)
    res = requests.request(
        method=request.method,
        url=forward_url,
        data=request.get_data(),
        allow_redirects=False,
    )

    log_event(method=request.method, url=request.url, forwarded_to=WDB_RETRO_URL)
    return Response(res.content, res.status_code)


def log_event(method, url, forwarded_to):
    current_GMT = time.gmtime()
    time_stamp = calendar.timegm(current_GMT)
    app.logger.info(
        {
            "method": method,
            "request": url,
            "forwarded_to": forwarded_to,
            "timestamp": time_stamp,
        }
    )

if __name__ == "__main__":
    app.run(port=PROXY_PORT, host='0.0.0.0')