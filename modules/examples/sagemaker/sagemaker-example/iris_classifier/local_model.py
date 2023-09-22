"""
A simple ML model that can classify different species of the Iris flower.

The Iris dataset contains 3 classes of flowers, Setosa, Versicolor, and Virginica.
Each class contains 4 features, 'Sepal length', 'Sepal width', 'Petal length', and 'Petal width'.
The aim of the Iris flower classification is to predict flowers based on their specific features.
"""

from sklearn.datasets import load_iris
from sklearn import model_selection, linear_model, metrics
import joblib

# Load Iris dataset.
data = load_iris()
x_data, y_data = load_iris(return_X_y=True)
features = data['feature_names']
labels = {
    label: name for label, name in zip(
        range(0, 3), data['target_names']
    )
}
print("\nFeatures: " + str(features))
print("Labels: " + str(labels))

# Total 150 samples.
# Split 60:40 for training and testing samples respectively.
x_train, x_test, y_train, y_test = model_selection.train_test_split(
    x_data,
    y_data,
    test_size=0.4,
    random_state=42
)

# Building and fitting the model.
classifier_model = linear_model.LogisticRegression(
    penalty='l2',
    C=100,
    random_state=10,
    multi_class='auto',
    max_iter=1000
)
classifier_model.fit(x_train, y_train)

# Prediction.
predictions = classifier_model.predict(x_test)
accuracy = metrics.accuracy_score(y_test, predictions)
print("\nAccuracy of Logistic Regression model:", accuracy)

# Create model.joblib file with current object value of classifier_model.
with open('model.joblib', 'wb') as f:
    joblib.dump(classifier_model, f)

# Construct predictor object from stored model for testing.
with open('model.joblib', 'rb') as f:
    predictor = joblib.load(f)

# Testing.
print("\nTesting following input: ")
sample_input = [[6.1, 2.8, 5.6, 2.2], [6.1, 2.8, 4.7, 1.2]]  # [2 1]
print(sample_input)

output = predictor.predict(sample_input)
print("\nOutput:")
print(output)
print([labels[cat_id] for cat_id in output], "\n")
