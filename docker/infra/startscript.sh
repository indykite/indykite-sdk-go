#!/bin/bash
# Adding a retry so let the istio start properly before the tests start

# pipefail - BASH only, not supported in POSIX Shell
set -o errexit -o nounset # -o pipefail

: "${APPUSER_HOME:?User home path is required}"
: "${BRANCH:="master"}"
: "${GITHUB_USER:?GitHub User is required}"
: "${GITHUB_TOKEN:?GitHub Token is required}"
: "${GITHUB:?GitHub FQDN is required}"
: "${SLACK_WEBHOOK_URL:=""}"
: "${BUCKET_NAME:?Bucket name is required}"
BUCKET_PATH="gs://${BUCKET_NAME}"

i=0
while [[ ${i} -le 5 ]]; do
    git clone --single-branch --branch "${BRANCH}" "https://${GITHUB_USER}:${GITHUB_TOKEN}@${GITHUB}"
    retVal=$?
    if [[ ${retVal} -ne 0 ]]; then
        sleep 1
        ((i++))
    else
        break
    fi
done
if [[ ${retVal} -ne 0 ]]; then
    echo "Failed to clone https://****:****@${GITHUB}"
    exit "${retVal}"
fi

cd ./*/ || exit
export RUN_ENV=${RUN_ENV:=${BRANCH}}
export RELEASE_VERSION="${RELEASE_VERSION:=unknown}"
export BUCKET_NAME="${BUCKET_NAME:=sdk_results_deploy}"
export SECRET_NAME=${SECRET_NAME:=goSdkTests}
export SDK_AUDIT_TABLE_NAME="${SDK_AUDIT_TABLE_NAME:=audit_log}"

# setup test environment variables
CUSTOMER_ID=$(gcloud secrets versions access latest --secret="${SECRET_NAME}" | jq --raw-output .goSdkTests.customerID)
LOCATION_ID=$(gcloud secrets versions access latest --secret="${SECRET_NAME}" | jq --raw-output .goSdkTests.appAgentCredentials.appSpaceId)
INDYKITE_APPLICATION_CREDENTIALS=$(gcloud secrets versions access latest --secret="${SECRET_NAME}" | jq --raw-output .goSdkTests.appAgentCredentials)
INDYKITE_SERVICE_ACCOUNT_CREDENTIALS=$(gcloud secrets versions access latest --secret="${SECRET_NAME}" | jq --raw-output .goSdkTests.serviceAccountCredentials)
export CUSTOMER_ID LOCATION_ID INDYKITE_APPLICATION_CREDENTIALS INDYKITE_SERVICE_ACCOUNT_CREDENTIALS

# setup reporting variables
run_date=$(date +%Y%m%d-%H%M)
result_file_name="${RELEASE_VERSION}_results_sdk_${RUN_ENV}_${run_date}_report.html"
storage="https://storage.cloud.google.com/${BUCKET_NAME}/${result_file_name}"

make install-tools report integration 2>output.txt
retVal=$?

# we are moving away of this kind of slack messaging, so it is optional for now
if [[ -n "${SLACK_WEBHOOK_URL}" ]]; then
    echo "Send results to Slack channel"

    if [[ ${retVal} -ne 0 ]]; then
        message="Test errors: ${retVal}"
        attachment_message=":alert: Tests failed :alert:"
        repair_message="To see the errors, open the logs or go to the github and launch the integration tests manually in actions"
        colour="#FF0000"
    else
        message="All tests passed"
        attachment_message=":heavy_check_mark: All Passed :heavy_check_mark:"
        repair_message="All good"
        colour="#008000"
    fi
    github_sha=$(git rev-parse --short HEAD)

    blocks=$(jq -n \
        --arg sha "${github_sha}" \
        --arg storage "${storage}" \
        --arg message "${message}" \
        --arg repair "${repair_message}" \
        --arg colour "${colour}" \
        --arg attach "${attachment_message}" \
        '{
    "blocks": [
      {"type": "divider"},
      {"type": "section", text: {"type": "mrkdwn", "text": "Test results - *Go SDK tests* - \($sha) - <\($storage)?authuser=0|Logs>"}},
      {"type": "section", fields: [
        {"type": "mrkdwn", text: $message},
        {"type": "mrkdwn", text: $repair}
      ]},
      {"type": "divider"}
    ],
    attachments: [
      {title: $message, color: $colour, fields: [$attach]}
    ]
  }')

    curl --header "Content-Type: application/json" --data "${blocks}" -X POST "${SLACK_WEBHOOK_URL}"
fi

echo "Send results to GCP bucket"

mv test-report.html "${result_file_name}"

echo "Copying the results to google cloud"
echo "Copying ${result_file_name} to ${BUCKET_PATH}"
gsutil -q cp "${result_file_name}" "${BUCKET_PATH}"

echo "Logs: ${storage}?authuser=0"
exit "${retVal}"
