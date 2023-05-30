import flask,json
from flask import request

server = flask.Flask(__name__)

@server.route('/',methods=['post'])
def main():
    # 处理paramsJson
    params = json.loads(request.data)
    x = params['x']

    x += 5

    # 生成returnJson
    rets = {'x' : x}
    return json.dumps(rets)

if __name__== '__main__':
    server.run(debug=True,port = 8888,host='0.0.0.0')