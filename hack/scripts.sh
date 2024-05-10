#!/bin/bash

kubectl apply -f app.yaml
#create a route and replace in tekton config


oc expose service event-listener

kubectl apply -f TektonConfig.yaml

#create tasks

kubectl apply -f tasks.yaml

#create pipeline

kubectl apply -f goodbye-pipeline.yaml


#Execute PipelineRun

kubectl apply -f goodbye-pipelinerun.yaml -n new



