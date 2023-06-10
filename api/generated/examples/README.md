# How to use foo client examples 

## Prerequisite

must have a custom apiserver implements hello.zeng.dev APIGrous runs behind kube-apiserver

## Play

basic client and applyconfiguration

        go run client/client.go

watch with informer then list foos from cache with lister

        go run watch/watch.go