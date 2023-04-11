package stader

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	stader_backend "github.com/stader-labs/stader-node/shared/types/stader-backend"
	"github.com/stader-labs/stader-node/shared/utils/net"
	"net/http"
)

const merkleProofGetterApi = "https://v6s3vqe7va.execute-api.us-east-1.amazonaws.com/prod/merklesForElRewards/%d/%s"
const merkleProofAggregateGetterApi = "https://v6s3vqe7va.execute-api.us-east-1.amazonaws.com/prod/merklesForElRewards/proofs/%s"

// TODO - akhilesh - get an api which will return the all merkle proofs for all cycles for a given operator
func GetCycleMerkleProofsForOperator(cycle int64, operator common.Address) (*stader_backend.CycleMerkleProofs, error) {
	res, err := net.MakeGetRequest(fmt.Sprintf(merkleProofGetterApi, cycle, operator.Hex()), struct{}{})
	fmt.Printf("Making api call to %s\n", fmt.Sprintf(merkleProofGetterApi, cycle, operator.Hex()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while getting merkle proofs for operator %s", operator.Hex())
	}
	var cycleMerkleProofs stader_backend.CycleMerkleProofs
	err = json.NewDecoder(res.Body).Decode(&cycleMerkleProofs)
	if err != nil {
		return nil, err
	}
	return &cycleMerkleProofs, nil
}

func GetAllMerkleProofsForOperator(operator common.Address) ([]*stader_backend.CycleMerkleProofs, error) {
	res, err := net.MakeGetRequest(fmt.Sprintf(merkleProofAggregateGetterApi, operator.Hex()), struct{}{})
	fmt.Printf("Making api call to %s\n", fmt.Sprintf(merkleProofAggregateGetterApi, operator.Hex()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while getting ALL merkle proofs for operator %s", operator.Hex())
	}

	var allMerkleProofs []*stader_backend.CycleMerkleProofs
	err = json.NewDecoder(res.Body).Decode(&allMerkleProofs)
	if err != nil {
		return nil, err
	}
	return allMerkleProofs, nil
}
