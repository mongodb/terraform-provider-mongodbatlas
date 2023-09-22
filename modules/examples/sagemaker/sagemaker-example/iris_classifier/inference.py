import joblib
import os
import json


def model_fn(model_dir):
    """
    Deserialize fitted model.
    """

    model = joblib.load(os.path.join(model_dir, "model.joblib"))
    return model


def input_fn(request_body, request_content_type):
    """
    input_fn

    request_body: The body of the request sent to the model.
    request_content_type: (string) specifies the format/variable type of the request
    """

    if request_content_type == 'application/json':
        request_body = json.loads(request_body)
        inpVar = request_body['Input']
        return inpVar
    else:
        raise ValueError("This model only supports application/json input")


def predict_fn(input_data, model):
    """
    predict_fn
        input_data: returned array from input_fn above
        model (sklearn model) returned model loaded from model_fn above
    """

    return model.predict(input_data)


def output_fn(prediction, content_type):
    """
    output_fn
        prediction: the returned value from predict_fn above
        content_type: the content type the endpoint expects to be returned. Ex: JSON, string

    """

    respJSON = {'Output': prediction.tolist()}
    return respJSON
