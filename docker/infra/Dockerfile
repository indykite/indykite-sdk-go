# checkov:skip=CKV_DOCKER_2:ensure that HEALTHCHECK instructions have been added
FROM golang:1.27-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc


SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

ENV APPUSER="appuser"
ENV APPGROUP="appgroup"
ENV APPUSER_HOME="/home/$APPUSER"
ENV CLOUDSDK_INSTALL_DIR="${APPUSER_HOME}"
ENV PATH="$PATH:${APPUSER_HOME}/google-cloud-sdk/bin"

# hadolint ignore=DL3018
RUN apk add --update --no-cache \
        bash \
        curl \
        jq \
        git \
        make \
        python3 \
        openssh-client \
    # Add new user and not using root to run the tests for security reasons
    && addgroup -S "$APPGROUP" --gid 10001 \
    && adduser -S "$APPUSER" --uid 10001 \
        -G "$APPGROUP" \
        --disabled-password \
        --gecos "" \
        --home "$APPUSER_HOME" \
    && mkdir -p "${APPUSER_HOME}/github" \
    && chown -R "$APPUSER":"$APPGROUP" "$APPUSER_HOME" \
    # Install gscloud for posting the results in the bucket
    && curl -sSL https://sdk.cloud.google.com | bash \
    && apk info -v


COPY startscript.sh "${APPUSER_HOME}/startscript.sh"
RUN chmod +x "${APPUSER_HOME}/startscript.sh"

# Switch to user (numeric UID:GID of $APPUSER:$APPGROUP)
USER 10001:10001
WORKDIR "${APPUSER_HOME}/github"

# Checkout github repository
ENV CI=true

# trivy:ignore:AVD-DS-0026 - TODO: Add HEALTHCHECK instruction in your Dockerfile
# HEALTHCHECK
CMD ["../startscript.sh"]
