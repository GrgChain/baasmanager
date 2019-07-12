## Cli 测试
本操作是进入K8s pod的cli容器里执行

### 创建channel
peer channel create -o orderer0.luotestfabricorderer:7050 -c mychannel -f ./channel-artifacts/mychannel.tx \
--tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/luotestfabricorderer/orderers/orderer0.luotestfabricorderer/msp/tlscacerts/tlsca.luotestfabricorderer-cert.pem
### 加入channel
peer channel join -b mychannel.block
### 更新组织锚节点
peer channel update -o orderer0.luotestfabricorderer:7050 -c mychannel  -f ./channel-artifacts/Org1MSPanchors.tx \
--tls --cafile  /etc/hyperledger/crypto-config/ordererOrganizations/luotestfabricorderer/orderers/orderer0.luotestfabricorderer/msp/tlscacerts/tlsca.luotestfabricorderer-cert.pem
### 下载链码
go get github.com/hyperledger/fabric-samples
### 安装链码
peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric-samples/chaincode/chaincode_example02/go/
### 实例化链码
peer chaincode instantiate -o orderer0.luotestfabricorderer:7050 \
--tls --cafile  /etc/hyperledger/crypto-config/ordererOrganizations/luotestfabricorderer/orderers/orderer0.luotestfabricorderer/msp/tlscacerts/tlsca.luotestfabricorderer-cert.pem \
-C mychannel -n mycc -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "OR ('Org1MSP.member')"
### 查询
peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
### 调用
peer chaincode invoke -o orderer0.luotestfabricorderer:7050 \
--tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/luotestfabricorderer/orderers/orderer0.luotestfabricorderer/msp/tlscacerts/tlsca.luotestfabricorderer-cert.pem \
-C mychannel -n mycc \
--peerAddresses peer0.luotestfabricorg1:7051 \
--tlsRootCertFiles /etc/hyperledger/crypto-config/peerOrganizations/luotestfabricorg1/peers/peer0.luotestfabricorg1/tls/ca.crt \
-c '{"Args":["invoke","a","b","10"]}'