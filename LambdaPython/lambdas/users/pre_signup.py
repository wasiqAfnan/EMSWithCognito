import json
import logging
import boto3
from botocore.exceptions import ClientError

logger = logging.getLogger()
logger.setLevel(logging.INFO)

cognito = boto3.client("cognito-idp")


def get_user_by_email(user_pool_id: str, email: str):
    """Retrieve a native Cognito user (username/password) matching the provided email."""
    safe_email = email.replace("\\", "\\\\").replace('"', '\\"')

    response = cognito.list_users(
        UserPoolId=user_pool_id,
        Filter=f'email = "{safe_email}"'
    )

    # Filter out federated users; we only want native accounts
    local_users = [
        user for user in response.get("Users", [])
        if user.get("UserStatus") != "EXTERNAL_PROVIDER"
    ]

    if not local_users:
        return None

    if len(local_users) > 1:
        raise Exception(f"Multiple native Cognito users found for email: {email}")

    return local_users[0]


def get_attribute(user, attribute_name: str):
    """Extract a specific attribute value from a Cognito user object."""
    for attr in user.get("Attributes", []):
        if attr.get("Name") == attribute_name:
            return attr.get("Value")
    return None


def lambda_handler(event, context):
    # 1. Ensure trigger is from an external provider (OAuth)
    if event.get("triggerSource") != "PreSignUp_ExternalProvider":
        return event

    try:
        user_pool_id = event["userPoolId"]
        cognito_username = event["userName"]
        
        # 2. Ensure the external provider is Google
        if not cognito_username.lower().startswith("google_"):
            logger.info("External provider is not Google. Skipping account linking.")
            return event

        google_sub = cognito_username.split("_", 1)[1]
        
        user_attributes = event["request"].get("userAttributes", {})
        email = user_attributes.get("email")

        if not email:
            logger.warning("No email attribute provided by Google.")
            return event

        # 3. Lookup existing native user account by email
        existing_user = get_user_by_email(user_pool_id, email)

        if not existing_user:
            logger.info("No existing native account found. Allowing normal Google signup.")
            return event

        existing_username = existing_user["Username"]

        # 4. Enforce email verification on the existing account before linking
        existing_email_verified = get_attribute(existing_user, "email_verified")
        if existing_email_verified != "true":
            raise Exception("Existing native account email is not verified. Cannot link Google identity.")

        # 5. Link the Google identity to the existing native Cognito user        
        cognito.admin_link_provider_for_user(
            UserPoolId=user_pool_id,
            DestinationUser={
                "ProviderName": "Cognito",
                "ProviderAttributeValue": existing_username
            },
            SourceUser={
                "ProviderName": "Google",
                "ProviderAttributeName": "Cognito_Subject",
                "ProviderAttributeValue": google_sub
            }
        )

        logger.info("Account linking successful.")
        return event

    except ClientError as e:
        logger.error("AWS Cognito API error", exc_info=True)
        raise
    except Exception as e:
        logger.error("Account linking failed", exc_info=True)
        raise
