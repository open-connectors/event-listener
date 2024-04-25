#!/bin/bash

kubectl apply -f app.yaml -n new
#create a route and replace in tekton config

kubectl apply -f TektonConfig.yaml

#create tasks

kubectl apply -f tasks.yaml -n new

#create pipeline

kubectl apply -f goodbye-pipeline.yaml -n new


#Execute PipelineRun

kubectl apply -f goodbye-pipelinerun.yaml -n new



