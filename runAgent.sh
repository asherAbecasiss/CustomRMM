#!/bin/bash





function run(){

    cd $(pwd)/agentFlow && ./agentFlow
}





while true; do run & sleep 3600; done