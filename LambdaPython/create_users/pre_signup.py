import json
import logging

import boto3
from botocore.exceptions import ClientError


logger = logging.getLogger()
logger.setLevel(logging.INFO)

cognito = boto3.client("cognito-idp")


def get_user_by_email(user_pool_id: str, email: str):
    """
    Find a LOCAL Cognito user with the given email.

    We ignore EXTERNAL_PROVIDER users because we only want
    the existing username/password Cognito account to be
    the destination of the link.
    """

    # Escape characters that have meaning in Cognito filter syntax
    safe_email = email.replace("\\", "\\\\").replace('"', '\\"')

    response = cognito.list_users(
        UserPoolId=user_pool_id,
        Filter=f'email = "{safe_email}"'
    )

    users = response.get("Users", [])

    # We only want a native Cognito user.
    local_users = [
        user
        for user in users
        if user.get("UserStatus") != "EXTERNAL_PROVIDER"
    ]

    if len(local_users) == 0:
        return None

    if len(local_users) > 1:
        raise Exception(
            f"Multiple local Cognito users found for email: {email}"
        )

    return local_users[0]


def get_attribute(user, attribute_name: str):
    """Get a Cognito user attribute."""
    for attribute in user.get("Attributes", []):
        if attribute.get("Name") == attribute_name:
            return attribute.get("Value")

    return None


def lambda_handler(event, context):
    logger.info("Received Cognito event: %s", json.dumps(event))

    # We only want to handle external provider sign-ins.
    if event.get("triggerSource") != "PreSignUp_ExternalProvider":
        return event

    try:
        user_pool_id = event["userPoolId"]
        provider_name = "Google"

        # Example:
        # event["userName"] = "Google_109220063452404746097"
        cognito_username = event["userName"]

        # Make sure this is actually the Google provider.
        if not cognito_username.lower().startswith("google_"):
            logger.info(
                "External provider is not Google. Skipping linking."
            )
            return event

        # Google user ID / Google sub.
        # Google_123456789
        #       ^^^^^^^^^
        google_sub = cognito_username.split("_", 1)[1]

        # Email supplied by the external provider and mapped by Cognito.
        user_attributes = event["request"].get("userAttributes", {})
        email = user_attributes.get("email")

        if not email:
            logger.warning(
                "No email found in external provider attributes."
            )
            return event

        logger.info(
            "Processing Google login for email=%s, google_sub=%s",
            email,
            google_sub
        )

        # ---------------------------------------------------------
        # 1. Find an existing native Cognito user by email
        # ---------------------------------------------------------

        existing_user = get_user_by_email(
            user_pool_id,
            email
        )

        # ---------------------------------------------------------
        # 2. No existing native account
        # ---------------------------------------------------------

        if existing_user is None:
            logger.info(
                "No existing local Cognito user found for %s. "
                "Allowing Cognito to create the Google user.",
                email
            )

            return event

        existing_username = existing_user["Username"]

        logger.info(
            "Existing local Cognito user found: %s",
            existing_username
        )

        # ---------------------------------------------------------
        # 3. Make sure the existing account's email is verified
        # ---------------------------------------------------------

        existing_email_verified = get_attribute(
            existing_user,
            "email_verified"
        )

        google_email_verified = user_attributes.get(
            "email_verified"
        )

        if existing_email_verified != "true":
            raise Exception(
                "An existing Cognito account was found with this email, "
                "but its email is not verified. "
                "Verify the existing account before linking Google."
            )

        if google_email_verified != "true":
            raise Exception(
                "The Google email is not verified. "
                "Google account cannot be linked."
            )

        # ---------------------------------------------------------
        # 4. Link Google identity → existing Cognito user
        # ---------------------------------------------------------

        logger.info(
            "Linking Google identity %s to Cognito user %s",
            google_sub,
            existing_username
        )

        cognito.admin_link_provider_for_user(
            UserPoolId=user_pool_id,

            # EXISTING local Cognito user
            DestinationUser={
                "ProviderName": "Cognito",
                "ProviderAttributeValue": existing_username
            },

            # NEW external Google identity
            SourceUser={
                "ProviderName": provider_name,
                "ProviderAttributeName": "Cognito_Subject",
                "ProviderAttributeValue": google_sub
            }
        )

        logger.info(
            "Successfully linked Google identity %s "
            "to Cognito user %s",
            google_sub,
            existing_username
        )

        # Let Cognito continue the sign-in.
        return event

    except ClientError as e:
        logger.error(
            "AWS Cognito error: %s",
            e,
            exc_info=True
        )

        # Fail the sign-in rather than allowing Cognito to
        # accidentally create another account.
        raise

    except Exception as e:
        logger.error(
            "Account linking failed: %s",
            e,
            exc_info=True
        )

        raise
