#!/bin/bash

kubectl apply -f deployment.yaml
kubectl apply -fd service.yaml
#create a route and replace in tekton config

kubectl apply -f TektonConfig.yaml

#create tasks

kubectl apply -f hello.yaml
kubectl apply -f goodbye.yaml

#create pipeline

kubectl apply -f goodbye-pipeline.yaml


#Execute PipelineRun

kubectl apply -f goodbye-pipelinerun.yaml



