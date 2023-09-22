FROM amazon/aws-lambda-python:3.8

ADD app.py ./

RUN pip install 'pymongo[srv]'==3.11

CMD ["app.handler"]