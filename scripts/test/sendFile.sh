#!/bin/bash

curl -XPOST http://localhost:9090/send-file -H 'Content-type: application/json' -d '{"filename":"ali.txt","data":"dsfsdfsdf","node_id":"hossain"}'