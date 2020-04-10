#!/bin/bash

kubectl create secret generic http-auth --from-file=./config/auth.txt
kubectl create configmap agent-config --from-file=./config/stream.yaml
