#!/usr/bin/env bash

set -e
programname=$0

usage() { printf "Usage: $0 [-i <local image name>] [-p <ACR repo name>] [-r <ACR registry name>] [-n <number of signatures to create>] [-t <whether each signature should have a distinct tag]" 1>&2; exit 1; }
creds() { printf "Please set NOTATION_USERNAME and NOTATION_PASSWORD" 1>&2; exit 1; }

separate_tags=false
while getopts ":i:p:r:n:t" o; do
  case "${o}" in
    i) image=${OPTARG} ;;
    p) repo=${OPTARG} ;;
    r) registry=${OPTARG} ;;
    n) num_sigs=${OPTARG} ;;
    t) separate_tags=true ;;
    *) usage ;;
  esac
done
shift "$((OPTIND-1))"
if [ -z "${image}" ] || [ -z "${registry}" ] || [ -z "${repo}" ] || [ -z "${num_sigs}" ]; then
  usage
fi

if [ -z "${NOTATION_USERNAME}" ] || [ -z "${NOTATION_PASSWORD}" ]; then
  creds
fi

docker tag ${image}:latest ${registry}.azurecr.io/${repo}:latest
docker push ${registry}.azurecr.io/${repo}:latest
if [ "${separate_tags}" = true ];
then
  notation sign ${registry}.azurecr.io/${repo}:latest
fi

for ((i=1;i<=${num_sigs};i++)); do
  if [ "${separate_tags}" = true ];
  then
    docker tag ${image}:latest ${registry}.azurecr.io/${repo}:${i}
    docker push ${registry}.azurecr.io/${repo}:${i}
  else
    notation sign ${registry}.azurecr.io/${repo}:latest
  fi
done
set +e
