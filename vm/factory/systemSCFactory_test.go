package factory

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	arwenConfig "github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/vm"
	"github.com/ElrondNetwork/elrond-go/vm/mock"
	"github.com/ElrondNetwork/elrond-go/vm/systemSmartContracts/defaults"
	"github.com/stretchr/testify/assert"
)

func createMockNewSystemScFactoryArgs() ArgsNewSystemSCFactory {
	gasMap := arwenConfig.MakeGasMapForTests()
	gasMap = defaults.FillGasMapInternal(gasMap, 1)
	gasSchedule := mock.NewGasScheduleNotifierMock(gasMap)
	return ArgsNewSystemSCFactory{
		SystemEI:            &mock.SystemEIStub{},
		Economics:           &mock.EconomicsHandlerStub{},
		SigVerifier:         &mock.MessageSignVerifierMock{},
		GasSchedule:         gasSchedule,
		NodesConfigProvider: &mock.NodesConfigProviderStub{},
		Marshalizer:         &mock.MarshalizerMock{},
		Hasher:              &mock.HasherMock{},
		SystemSCConfig: &config.SystemSmartContractsConfig{
			ESDTSystemSCConfig: config.ESDTSystemSCConfig{
				BaseIssuingCost: "100000000",
				OwnerAddress:    "aaaaaa",
			},
			GovernanceSystemSCConfig: config.GovernanceSystemSCConfig{
				ProposalCost:     "500",
				NumNodes:         100,
				MinQuorum:        50,
				MinPassThreshold: 50,
				MinVetoThreshold: 50,
			},
			StakingSystemSCConfig: config.StakingSystemSCConfig{
				GenesisNodePrice:                     "1000",
				UnJailValue:                          "10",
				MinStepValue:                         "10",
				MinStakeValue:                        "1",
				UnBondPeriod:                         1,
				StakingV2Epoch:                       1,
				StakeEnableEpoch:                     0,
				NumRoundsWithoutBleed:                1,
				MaximumPercentageToBleed:             1,
				BleedPercentagePerRound:              1,
				MaxNumberOfNodesForStake:             100,
				ActivateBLSPubKeyMessageVerification: false,
				MinUnstakeTokensValue:                "1",
			},
			DelegationSystemSCConfig: config.DelegationSystemSCConfig{
				MinStakeAmount: "10",
				EnabledEpoch:   0,
				MinServiceFee:  0,
				MaxServiceFee:  10000,
			},
			DelegationManagerSystemSCConfig: config.DelegationManagerSystemSCConfig{
				BaseIssuingCost:    "10",
				MinCreationDeposit: "10",
				EnabledEpoch:       0,
			},
		},
		EpochNotifier:          &mock.EpochNotifierStub{},
		AddressPubKeyConverter: &mock.PubkeyConverterMock{},
	}
}

func TestNewSystemSCFactory_NilSystemEI(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.SystemEI = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilSystemEnvironmentInterface, err)

	value, _ := big.NewInt(0).SetString("2500000000000000000000", 10)
	fmt.Println(hex.EncodeToString(value.Bytes()))
}

func TestNewSystemSCFactory_NilEconomicsData(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.Economics = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilEconomicsData, err)
}

func TestNewSystemSCFactory_NilPubKeyConverter(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	arguments.AddressPubKeyConverter = nil
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, scFactory)
	assert.Equal(t, vm.ErrNilAddressPubKeyConverter, err)
}

func TestNewSystemSCFactory_Ok(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	scFactory, err := NewSystemSCFactory(arguments)

	assert.Nil(t, err)
	assert.NotNil(t, scFactory)
}

func TestSystemSCFactory_Create(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	scFactory, _ := NewSystemSCFactory(arguments)

	container, err := scFactory.Create()
	assert.Nil(t, err)
	assert.Equal(t, 6, container.Len())
}

func TestSystemSCFactory_CreateForGenesis(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	scFactory, _ := NewSystemSCFactory(arguments)

	container, err := scFactory.CreateForGenesis()
	assert.Nil(t, err)
	assert.Equal(t, 4, container.Len())
}

func TestSystemSCFactory_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	arguments := createMockNewSystemScFactoryArgs()
	scFactory, _ := NewSystemSCFactory(arguments)
	assert.False(t, scFactory.IsInterfaceNil())

	scFactory = nil
	assert.True(t, check.IfNil(scFactory))
}
