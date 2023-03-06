package node

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stader-labs/stader-node/stader-lib/stader"
	"math/big"
)

func EstimateOnboardNodeOperator(pnr *stader.PermissionlessNodeRegistryContractManager, mevSocialize bool, operatorName string, operatorRewarderAddress common.Address, opts *bind.TransactOpts) (stader.GasInfo, error) {
	return pnr.PermissionlessNodeRegistryContract.GetTransactionGasInfo(opts, "onboardNodeOperator", mevSocialize, operatorName, operatorRewarderAddress)
}

func OnboardNodeOperator(pnr *stader.PermissionlessNodeRegistryContractManager, mevSocialize bool, operatorName string, operatorRewarderAddress common.Address, opts *bind.TransactOpts) (*types.Transaction, error) {
	tx, err := pnr.PermissionlessNodeRegistry.OnboardNodeOperator(opts, mevSocialize, operatorName, operatorRewarderAddress)
	if err != nil {
		return nil, fmt.Errorf("Could not onboard node operator: %w", err)
	}

	return tx, nil
}

func GetOperatorId(pnr *stader.PermissionlessNodeRegistryContractManager, nodeAddress common.Address, opts *bind.CallOpts) (*big.Int, error) {
	operatorId, err := pnr.PermissionlessNodeRegistry.OperatorIDByAddress(opts, nodeAddress)
	if err != nil {
		return nil, err
	}

	return operatorId, nil
}

func GetOperatorInfo(pnr *stader.PermissionlessNodeRegistryContractManager, operatorId *big.Int, opts *bind.CallOpts) (struct {
	Active                  bool
	OptedForSocializingPool bool
	OperatorName            string
	OperatorRewardAddress   common.Address
	OperatorAddress         common.Address
}, error) {
	operatorInfo, err := pnr.PermissionlessNodeRegistry.OperatorStructById(nil, operatorId)
	if err != nil {
		return struct {
			Active                  bool
			OptedForSocializingPool bool
			OperatorName            string
			OperatorRewardAddress   common.Address
			OperatorAddress         common.Address
		}{}, err
	}

	return operatorInfo, nil
}