import boto3
import json
import os

# Iris flower categories.
irisCategory = {
    0: 'setosa',
    1: 'versicolor',
    2: 'virginica'
}


def handler(event, context):
    """
    Lambda handler. Processes MongoDB Change Events, invokes SageMaker enpoint
    with input read from event and writes results back to event bus.
    """

    print(json.dumps(event))

    try:
        # Read environment variables.
        SAGEMAKER_ENDPOINT = os.environ['model_endpoint']
        REGION_NAME = os.environ['region_name']
        EVENTBUS_NAME = os.environ['eventbus_name']

        # Enable the SageMaker runtime.
        runtime = boto3.client(
            'runtime.sagemaker', region_name=REGION_NAME
        )

        # Read MongoDB change event.
        doc = event['detail']['fullDocument']
        payload = json.dumps({'Input': doc['data']})

        # Predict from model.
        response = runtime.invoke_endpoint(
            EndpointName=SAGEMAKER_ENDPOINT,
            ContentType='application/json',
            Body=payload
        )
        output = json.loads(response['Body'].read().decode())['Output']

        # Write result back to eventbus.
        prediction = [irisCategory[catID] for catID in output]
        response = push_to_eventbus(
            EVENTBUS_NAME, REGION_NAME, prediction, doc['_id']
        )

        print(json.dumps(response))

        return response
    except Exception as ex:
        print("Exception: " + str(ex))
        raise ex


# Push events to eventbus.
def push_to_eventbus(EVENTBUS_NAME, REGION_NAME, prediction, inputID):
    # Set up client for AWS.
    client = boto3.client(
        'events',
        region_name=REGION_NAME
    )

    # Create JSON for pushing to eventbus.
    detailJsonString = {
        "prediction": prediction,
        "inp_id": inputID
    }

    # Put events to eventbus.
    response = client.put_events(
        Entries=[
            {
                'Source': 'user-event',
                'DetailType': 'user-preferences',
                'Detail': json.dumps(detailJsonString),
                'EventBusName': EVENTBUS_NAME
            }
        ]
    )

    return response
