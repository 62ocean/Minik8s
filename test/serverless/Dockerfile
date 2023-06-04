FROM python:3.11
WORKDIR ./addFive
ADD . .
RUN pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple
EXPOSE 8888
CMD ["python", "./addFive.py"]
