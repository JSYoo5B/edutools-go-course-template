# Style: code inside a test

For a testing tool (or any API that is normally used from tests) the exercise is the body of a visible test.
Setup and cleanup are given, the learner writes the calls, and the given assertions at the end grade the result.
A small hidden test makes sure those assertions were not edited away.

In `test/hello_test.go`, send `GET /hello?name=Go` to the test server and read the response body into `body`.

Use this style when learners will use the subject inside their own tests. Keep the given assertions short and
specific so their failure messages explain what is missing.
