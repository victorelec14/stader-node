package node

import (
	"fmt"
	"math/big"

	"github.com/stader-labs/stader-node/shared/services/gas"
	"github.com/stader-labs/stader-node/shared/services/stader"
	cliutils "github.com/stader-labs/stader-node/shared/utils/cli"
	"github.com/stader-labs/stader-node/shared/utils/math"
	"github.com/stader-labs/stader-node/stader-lib/utils/eth"
	"github.com/urfave/cli"
)

func ClaimRewards(c *cli.Context) error {
	staderClient, err := stader.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer staderClient.Close()

	// Check and assign the EC status
	err = cliutils.CheckClientStatus(staderClient)
	if err != nil {
		return err
	}

	// Print what network we're on
	err = cliutils.PrintNetwork(staderClient)

	// Check if we can Withdraw El Rewards
	canClaimRewardsResponse, err := staderClient.CanClaimRewards()
	if err != nil {
		return err
	}
	if canClaimRewardsResponse.NoRewards {
		fmt.Println("No rewards to claim.")
		return nil
	}

	// TODO:
	// 1. How much SD fee need to repay
	// 2. Based on the current Health Factor, you can claim upto <ETH amount> ETH.

	if canClaimRewardsResponse.SdStatusResponse.SdUtilizerLatestBalance.Cmp(big.NewInt(0)) > 0 {
		cliutils.PrintError(fmt.Sprintf("You need to first pay %f and close the utilization position to get back your funds. Execute the following command to repay your utilized SD stader-cli repay-sd --amount <SD amount> \n", eth.WeiToEth(canClaimRewardsResponse.SdStatusResponse.SdUtilizerLatestBalance)))
		return nil
	}

	err = gas.AssignMaxFeeAndLimit(canClaimRewardsResponse.GasInfo, staderClient, c.Bool("yes"))
	if err != nil {
		return err
	}

	// Prompt for confirmation
	if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf(
		"Are you sure you want to send rewards to your operator reward address?"))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Withdraw El Rewards
	res, err := staderClient.ClaimRewards()
	if err != nil {
		return err
	}
	fmt.Printf("Withdrawing %.6f ETH Rewards to Operator Reward Address: %s\n\n", math.RoundDown(eth.WeiToEth(res.OperatorRewardsBalance), 6), res.OperatorRewardAddress)
	cliutils.PrintTransactionHash(staderClient, res.TxHash)
	if _, err = staderClient.WaitForTransaction(res.TxHash); err != nil {
		return err
	}

	// Log & return
	fmt.Printf("Withdrawn %.6f ETH Rewards to Operator Reward Address: %s\n\n", math.RoundDown(eth.WeiToEth(res.OperatorRewardsBalance), 6), res.OperatorRewardAddress)
	return nil
}
