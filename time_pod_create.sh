#!/usr/bin/env bash
programname=$0

usage() { printf "Usage: $0 [-y <pod yaml name>] [-u whether to roll Ratify before each create]" 1>&2; exit 1; }

uninstall=false
while getopts ":i:y:r:n:u" o; do
  case "${o}" in
    y) yaml_name=${OPTARG} ;;
    u) uninstall=true ;;
    *) usage ;;
  esac
done
shift "$((OPTIND-1))"
if [ -z "${yaml_name}" ] || [ -z "${uninstall}" ]; then
  usage
fi

for i in {1..100}; do
  if [ "${uninstall}" = true ];
  then
    helm uninstall ratify -n ratify-service
    helm install ratify ../ratify/charts/ratify --atomic --namespace ratify-service --create-namespace --set image.repository=anlandu/ratify --set image.tag=nosleep
  fi
  export KUBECTLPERFINDEX="1image$i"
  yq -i '.metadata.name = strenv(KUBECTLPERFINDEX)' ./pods/${yaml_name}.yaml
  start=$(($(date +%s%N)/1000000))
  kubectl apply -f ./pods/${yaml_name}.yaml
  echo $(($(date +%s%N)/1000000 - $start))
  kubectl delete -f ./pods/${yaml_name}.yaml
done