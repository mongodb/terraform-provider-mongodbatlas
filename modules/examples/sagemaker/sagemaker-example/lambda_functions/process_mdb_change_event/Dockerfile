FROM amazon/aws-lambda-python:3.8

ADD app.py ./

RUN pip install boto3

CMD ["app.handler"]
