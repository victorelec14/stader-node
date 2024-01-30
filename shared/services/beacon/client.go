/*
This work is licensed and released under GNU GPL v3 or any other later versions.
The full text of the license is below/ found at <http://www.gnu.org/licenses/>

(c) 2023 Rocket Pool Pty Ltd. Modified under GNU GPL v3. [1.4.7]

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package beacon

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/stader-labs/stader-node/shared/types/config"
	"github.com/stader-labs/stader-node/stader-lib/types"
)

// API request options
type ValidatorStatusOptions struct {
	Epoch *uint64
	Slot  *uint64
}

// API response types
type SyncStatus struct {
	Syncing  bool
	Progress float64
}
type Eth2Config struct {
	GenesisForkVersion           []byte
	GenesisValidatorsRoot        []byte
	GenesisEpoch                 uint64
	GenesisTime                  uint64
	SecondsPerSlot               uint64
	SlotsPerEpoch                uint64
	SecondsPerEpoch              uint64
	EpochsPerSyncCommitteePeriod uint64
}
type Eth2DepositContract struct {
	ChainID uint64
	Address common.Address
}
type BeaconHead struct {
	Epoch                  uint64
	FinalizedEpoch         uint64
	JustifiedEpoch         uint64
	PreviousJustifiedEpoch uint64
}
type ValidatorStatus struct {
	Pubkey                     types.ValidatorPubkey
	Index                      uint64
	WithdrawalCredentials      common.Hash
	Balance                    uint64
	Status                     ValidatorState
	EffectiveBalance           uint64
	Slashed                    bool
	ActivationEligibilityEpoch uint64
	ActivationEpoch            uint64
	ExitEpoch                  uint64
	WithdrawableEpoch          uint64
	Exists                     bool
}
type Eth1Data struct {
	DepositRoot  common.Hash
	DepositCount uint64
	BlockHash    common.Hash
}
type BeaconBlock struct {
	Slot                 uint64
	ProposerIndex        uint64
	HasExecutionPayload  bool
	Attestations         []AttestationInfo
	FeeRecipient         common.Address
	ExecutionBlockNumber uint64
}

type Committee struct {
	Index      uint64
	Slot       uint64
	Validators []uint64
}

type AttestationInfo struct {
	AggregationBits bitfield.Bitlist
	SlotIndex       uint64
	CommitteeIndex  uint64
}

type NodeVersion struct {
	Version string
}

// Beacon client type
type BeaconClientType int

const (
	// This client is a traditional "split process" design, where the beacon
	// client and validator process are separate and run in different
	// containers
	SplitProcess BeaconClientType = iota

	// This client is a "single process" where the beacon client and
	// validator run in the same process (or run as separate processes
	// within the same docker container)
	SingleProcess

	// Unknown / missing client type
	Unknown
)

type ValidatorState string

const (
	ValidatorState_PendingInitialized ValidatorState = "pending_initialized"
	ValidatorState_PendingQueued      ValidatorState = "pending_queued"
	ValidatorState_ActiveOngoing      ValidatorState = "active_ongoing"
	ValidatorState_ActiveExiting      ValidatorState = "active_exiting"
	ValidatorState_ActiveSlashed      ValidatorState = "active_slashed"
	ValidatorState_ExitedUnslashed    ValidatorState = "exited_unslashed"
	ValidatorState_ExitedSlashed      ValidatorState = "exited_slashed"
	ValidatorState_WithdrawalPossible ValidatorState = "withdrawal_possible"
	ValidatorState_WithdrawalDone     ValidatorState = "withdrawal_done"
)

// Beacon client interface
type Client interface {
	GetClientType() (BeaconClientType, error)
	GetSyncStatus() (SyncStatus, error)
	GetNodeVersion() (NodeVersion, error)
	GetEth2Config() (Eth2Config, error)
	GetEth2DepositContract() (Eth2DepositContract, error)
	GetAttestations(blockId string) ([]AttestationInfo, bool, error)
	GetBeaconBlock(blockId string) (BeaconBlock, bool, error)
	GetBeaconHead() (BeaconHead, error)
	GetValidatorStatusByIndex(index string, opts *ValidatorStatusOptions) (ValidatorStatus, error)
	GetValidatorStatus(pubkey types.ValidatorPubkey, opts *ValidatorStatusOptions) (ValidatorStatus, error)
	GetValidatorStatuses(pubkeys []types.ValidatorPubkey, opts *ValidatorStatusOptions) (map[types.ValidatorPubkey]ValidatorStatus, error)
	GetValidatorIndex(pubkey types.ValidatorPubkey) (uint64, error)
	GetValidatorSyncDuties(indices []uint64, epoch uint64) (map[uint64]bool, error)
	GetValidatorProposerDuties(indices []uint64, epoch uint64) (map[uint64]uint64, error)
	GetExitDomainData(domainType []byte, network config.Network) ([]byte, error)
	ExitValidator(validatorIndex, epoch uint64, signature types.ValidatorSignature) error
	Close() error
	GetEth1DataForEth2Block(blockId string) (Eth1Data, bool, error)
	GetCommitteesForEpoch(epoch *uint64) ([]Committee, error)
}
