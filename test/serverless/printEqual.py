import flask,json
from flask import request

server = flask.Flask(__name__)

@server.route('/',methods=['post'])
def main():
    return "equal!"

if __name__== '__main__':
    server.run(debug=True,port = 8888,host='0.0.0.0')