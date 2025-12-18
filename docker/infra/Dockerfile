# checkov:skip=CKV_DOCKER_2:ensure that HEALTHCHECK instructions have been added

FROM golang:1.25-alpine@sha256:72567335df90b4ed71c01bf91fb5f8cc09fc4d5f6f21e183a085bafc7ae1bec8

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

# Switch to user
USER "$APPUSER"
WORKDIR "${APPUSER_HOME}/github"

# Checkout github repository
ENV CI=true

# trivy:ignore:AVD-DS-0026 - TODO: Add HEALTHCHECK instruction in your Dockerfile
# HEALTHCHECK
CMD ["../startscript.sh"]
