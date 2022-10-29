#!/usr/bin/env bash
for i in {1..100}; do
  export KUBECTLPERFINDEX="75images$i"
  yq -i '.metadata.name = strenv(KUBECTLPERFINDEX)' ./pods/75images.yaml
  start=$(($(date +%s%N)/1000000))
  kubectl apply -f ./pods/75images.yaml
  echo $(($(date +%s%N)/1000000 - $start))
  kubectl delete -f ./pods/75images.yaml
done