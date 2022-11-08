#!/usr/bin/env bash

# NOTE: This script assumes you are pushing to an ACR which has admin
# credentials enabled and has REGISTRY_USERNAME and REGISTRY_PASSWORD
# populated

set -e
programname=$0

usage() { printf "Usage: $0 [-p <ACR repo name>] [-r <ACR registry name>] [-n <number of signatures to create per subject>] [-s <number of subjects to create]" 1>&2; exit 1; }
creds() { printf "Please set NOTATION_USERNAME and NOTATION_PASSWORD" 1>&2; exit 1; }

while getopts "p:r:n:s:" o; do
  case "${o}" in
    p) repo=${OPTARG} ;;
    r) registry=${OPTARG} ;;
    n) num_sigs=${OPTARG} ;;
    s) num_subjects=${OPTARG} ;;
    *) usage ;;
  esac
done

# check all parameters are specified
if [ -z "${registry}" ] || [ -z "${repo}" ] || [ -z "${num_sigs}" ] || [ -z "${num_subjects}" ]; then
  usage
fi

# check registry credentials are populated as environment variables
if [ -z "${REGISTRY_USERNAME}" ] || [ -z "${REGISTRY_PASSWORD}" ]; then
  creds
fi

# login into ACR with admin credentials
docker login ${registry}.azurecr.io -u ${REGISTRY_USERNAME} -p ${REGISTRY_PASSWORD}

# for each subject
for ((i=1;i<=${num_subjects};i++)); do
  # build a unique scratch dockerfile and build image
  echo $'FROM scratch\nCMD ["echo", "image '${i}'"]' > Dockerfile
  docker build . -t ${registry}.azurecr.io/${repo}:${i}
  # push new image to specified registry, repo, and tag
  docker push ${registry}.azurecr.io/${repo}:${i}
  # delete Dockerfile and image created to reduce clutter
  rm Dockerfile
  docker image rm ${registry}.azurecr.io/${repo}:${i}
  # add specified number of signatures
  for ((j=1;j<=${num_sigs};j++)); do
      notation sign -u ${REGISTRY_USERNAME} -p ${REGISTRY_PASSWORD} ${registry}.azurecr.io/${repo}:${i}
  done
done
set +e
