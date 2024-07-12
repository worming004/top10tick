# This repository show how to implement a top 10 of tick name in stock marker

## Known flow (eventually to do list)

- No usage of schema registry
- For simplicity of demo, all ticks name are generated randomly. Each time you restart services, new tick names are generated
- Top 10 is all time biggest delta. Not by day like usually real apps do

## Install initial kafka cluster

Follow https://github.com/worming004/kafka-minikube-getting-started to install a kafka cluster

## Huge number of pods

Minikube can be started with increased number of pods. Like `min start --cpus=12 --memory=80g --extra-config=kubelet.max-pods=1000`

