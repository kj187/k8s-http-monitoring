#!/usr/bin/env bash

if [ -z "${REPOSITORY}" ] ; then error_exit "REPOSITORY not set"; fi
if [ -z "${REF}" ] ; then error_exit "REF not set"; fi

IMAGE_ID=$(echo "${REPOSITORY}" | tr '[A-Z]' '[a-z]')
echo "IMAGE_ID: ${IMAGE_ID}"

VERSION=$(echo "${REF}" | sed -e 's,.*/\(.*\),\1,')  # Strip git ref prefix from version
[[ "${REF}" == "refs/tags/"* ]] && VERSION=$(echo ${VERSION} | sed -e 's/^v//')  # Strip "v" prefix from tag name
[ "${VERSION}" == "master" ] && VERSION=latest  # Use Docker `latest` tag convention
echo "VERSION: ${VERSION}"

docker tag $IMAGE_ID $IMAGE_ID:$VERSION
docker push $IMAGE_ID