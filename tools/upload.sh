#!/bin/bash
# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

set -eux -o pipefail
ARTIFACTORY_URL="https://hcss.jfrog.io/artifactory"
ARTIFACTORY_REPO=lh-cdc-charts
for TARBALL in $(ls ~/project/tarballs); do
   ARTIFACT_MD5_CHECKSUM=$(md5sum ~/project/tarballs/${TARBALL} | awk '{print $1}')
   ARTIFACT_SHA1_CHECKSUM=$(sha1sum ~/project/tarballs/${TARBALL} | awk '{ print $1 }')
   ARTIFACT_SHA256_CHECKSUM=$(sha256sum ~/project/tarballs/${TARBALL} | awk '{ print $1 }')
   
   # Push helm package
   curl -u ${JFROG_USERNAME}:${JFROG_PASSWORD} \
      --header "X-Checksum-MD5:${ARTIFACT_MD5_CHECKSUM}" \
      --header "X-Checksum-Sha1:${ARTIFACT_SHA1_CHECKSUM}" \
      --header "X-Checksum-Sha256:${ARTIFACT_SHA256_CHECKSUM}" \
      --show-error --fail \
      --connect-timeout 5 --max-time 120 --retry 5 \
      --retry-delay 0 --retry-max-time 120 \
      -T tarballs/${TARBALL} \
      "${ARTIFACTORY_URL}/${ARTIFACTORY_REPO}/${TARBALL}"
   set -x
done

