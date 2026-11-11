#!/usr/bin/env bash
set -uo pipefail

PROJECT=dwk-gke-513717
CLUSTER=dwk-cluster
ZONE=europe-north1-b
REGION=europe-north1

read -r -p "Delete cluster $CLUSTER and all images in $PROJECT? Type nuke: " answer
[ "$answer" = "nuke" ] || { echo "Aborted"; exit 1; }

if gcloud container clusters describe "$CLUSTER" --zone "$ZONE" --project "$PROJECT" >/dev/null 2>&1; then
  KUBECONFIG=$(mktemp)
  export KUBECONFIG
  trap 'rm -f "$KUBECONFIG"' EXIT
  gcloud container clusters get-credentials "$CLUSTER" --zone "$ZONE" --project "$PROJECT" || exit 1

  kubectl delete gateway,httproute,ingress --all --all-namespaces --timeout=300s

  for ns in $(kubectl get namespaces -o name | grep -Ev '^namespace/(default|kube-.*|gke-.*|gmp-.*)$'); do
    kubectl delete "$ns" --timeout=300s
  done
  kubectl delete deployment,statefulset,pvc --all -n default --timeout=300s

  for _ in $(seq 30); do
    [ -z "$(kubectl get pv -o name)" ] && break
    sleep 10
  done

  gcloud container clusters delete "$CLUSTER" --zone "$ZONE" --project "$PROJECT" --quiet || exit 1
else
  echo "Cluster $CLUSTER not found, skipping"
fi

for name in $(gcloud artifacts repositories list --project "$PROJECT" --location "$REGION" --filter='format=DOCKER' --format='value(name)'); do
  repo=${name##*/}
  gcloud artifacts repositories delete "$repo" --project "$PROJECT" --location "$REGION" --quiet || exit 1
  gcloud artifacts repositories create "$repo" --project "$PROJECT" --location "$REGION" --repository-format=docker || exit 1
done

echo "Leftovers that still cost money, if any:"
gcloud compute disks list --project "$PROJECT" --filter='name~^pvc-'
gcloud compute forwarding-rules list --project "$PROJECT"
