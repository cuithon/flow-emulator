/*
 * Flow Emulator
 *
 * Copyright Flow Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package emulator_test

import (
	"fmt"
	"testing"

	flowgo "github.com/onflow/flow-go/model/flow"

	"github.com/onflow/flow-emulator/emulator"
	"github.com/stretchr/testify/require"
)

func TestCommonContractsDeployment(t *testing.T) {

	t.Parallel()

	//only test monotonic and emulator ( mainnet / testnet is used for remote debugging )
	chains := []flowgo.Chain{
		flowgo.Emulator.Chain(),
		flowgo.MonotonicEmulator.Chain(),
	}

	for _, chain := range chains {
		contracts := emulator.NewCommonContracts(chain)

		b, err := emulator.New(
			emulator.Contracts(contracts),
			emulator.WithChainID(chain.ChainID()),
		)
		require.NoError(t, err)

		for _, contract := range contracts {

			require.Equal(t, contract.Address.Hex(), chain.ServiceAddress().Hex())

			scriptCode := fmt.Sprintf(`
			access(all) fun main() {
				getAccount(0x%s).contracts.get(name: "%s") ?? panic("contract is not deployed")
	    	}`, contract.Address, contract.Name)

			scriptResult, err := b.ExecuteScript([]byte(scriptCode), [][]byte{})
			require.NoError(t, err)
			require.NoError(t, scriptResult.Error)

		}
	}
}
