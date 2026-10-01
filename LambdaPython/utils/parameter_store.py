import boto3

ssm = boto3.client("ssm", region_name="ap-south-1")


def get_parameter(name):
    response = ssm.get_parameter(
        Name=name,
        WithDecryption=True
    )

    return response["Parameter"]["Value"]