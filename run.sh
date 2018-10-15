#bin/bash

nohup ./GatewayServer/GatewayServer & 0 > ./log/gatewayserver.error 1 > ./log/gatewayserver.log
nohup ./LoginServer/LoginServer & 0 > ./log/loginserver.error 1 > ./log/loginserver.log
