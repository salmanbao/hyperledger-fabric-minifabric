/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
)

func main() {

	chaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating IoT data verification chaincode: %v", err)
		return
	}
	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting IoT data verification chaincode: %v", err)
	}

}
