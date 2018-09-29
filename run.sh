#bin/bash

nohup ./GatewayServer/GatewayServer & 0 > gatewayserver.error 1 > gatewayserver.log
nohup ./LoginServer/LoginServer & 0 > loginserver.error 1 > loginserver.log
