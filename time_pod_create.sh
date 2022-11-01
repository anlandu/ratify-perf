#!/usr/bin/env bash
for i in {1..100}; do
  helm uninstall ratify -n ratify-service
  helm install ratify ../ratify/charts/ratify --atomic --namespace ratify-service --create-namespace --set image.repository=anlandu/ratify --set image.tag=nosleep
  export KUBECTLPERFINDEX="1image$i"
  yq -i '.metadata.name = strenv(KUBECTLPERFINDEX)' ./pods/1image.yaml
  start=$(($(date +%s%N)/1000000))
  kubectl apply -f ./pods/1image.yaml
  echo $(($(date +%s%N)/1000000 - $start))
  kubectl delete -f ./pods/1image.yaml
done