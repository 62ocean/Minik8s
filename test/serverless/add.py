import flask, json
from flask import request

server = flask.Flask(__name__)


@server.route('/', methods=['post'])
def main():
    # 处理paramsJson
    params = json.loads(request.data)
    for param in params:
        if param['Name'] == 'x':
            x = param['Value']
        if param['Name'] == 'y':
            y = param['Value']

    # 函数逻辑
    x = x + y

    # 生成returnJson
    rets = [{'Name': 'x', 'Value': x}]
    return json.dumps(rets)


if __name__ == '__main__':
    server.run(debug=True, port=8888, host='0.0.0.0')
