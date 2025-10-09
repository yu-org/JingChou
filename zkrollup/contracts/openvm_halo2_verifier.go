// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// OpenVmHalo2VerifierMetaData contains all meta data concerning the OpenVmHalo2Verifier contract.
var OpenVmHalo2VerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"publicValues\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"appExeCommit\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"appVmCommit\",\"type\":\"bytes32\"}],\"name\":\"verify\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// OpenVmHalo2VerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use OpenVmHalo2VerifierMetaData.ABI instead.
var OpenVmHalo2VerifierABI = OpenVmHalo2VerifierMetaData.ABI

// OpenVmHalo2Verifier is an auto generated Go binding around an Ethereum contract.
type OpenVmHalo2Verifier struct {
	OpenVmHalo2VerifierCaller     // Read-only binding to the contract
	OpenVmHalo2VerifierTransactor // Write-only binding to the contract
	OpenVmHalo2VerifierFilterer   // Log filterer for contract events
}

// OpenVmHalo2VerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type OpenVmHalo2VerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenVmHalo2VerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OpenVmHalo2VerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenVmHalo2VerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OpenVmHalo2VerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenVmHalo2VerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OpenVmHalo2VerifierSession struct {
	Contract     *OpenVmHalo2Verifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// OpenVmHalo2VerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OpenVmHalo2VerifierCallerSession struct {
	Contract *OpenVmHalo2VerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// OpenVmHalo2VerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OpenVmHalo2VerifierTransactorSession struct {
	Contract     *OpenVmHalo2VerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// OpenVmHalo2VerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type OpenVmHalo2VerifierRaw struct {
	Contract *OpenVmHalo2Verifier // Generic contract binding to access the raw methods on
}

// OpenVmHalo2VerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OpenVmHalo2VerifierCallerRaw struct {
	Contract *OpenVmHalo2VerifierCaller // Generic read-only contract binding to access the raw methods on
}

// OpenVmHalo2VerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OpenVmHalo2VerifierTransactorRaw struct {
	Contract *OpenVmHalo2VerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOpenVmHalo2Verifier creates a new instance of OpenVmHalo2Verifier, bound to a specific deployed contract.
func NewOpenVmHalo2Verifier(address common.Address, backend bind.ContractBackend) (*OpenVmHalo2Verifier, error) {
	contract, err := bindOpenVmHalo2Verifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OpenVmHalo2Verifier{OpenVmHalo2VerifierCaller: OpenVmHalo2VerifierCaller{contract: contract}, OpenVmHalo2VerifierTransactor: OpenVmHalo2VerifierTransactor{contract: contract}, OpenVmHalo2VerifierFilterer: OpenVmHalo2VerifierFilterer{contract: contract}}, nil
}

// NewOpenVmHalo2VerifierCaller creates a new read-only instance of OpenVmHalo2Verifier, bound to a specific deployed contract.
func NewOpenVmHalo2VerifierCaller(address common.Address, caller bind.ContractCaller) (*OpenVmHalo2VerifierCaller, error) {
	contract, err := bindOpenVmHalo2Verifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OpenVmHalo2VerifierCaller{contract: contract}, nil
}

// NewOpenVmHalo2VerifierTransactor creates a new write-only instance of OpenVmHalo2Verifier, bound to a specific deployed contract.
func NewOpenVmHalo2VerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*OpenVmHalo2VerifierTransactor, error) {
	contract, err := bindOpenVmHalo2Verifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OpenVmHalo2VerifierTransactor{contract: contract}, nil
}

// NewOpenVmHalo2VerifierFilterer creates a new log filterer instance of OpenVmHalo2Verifier, bound to a specific deployed contract.
func NewOpenVmHalo2VerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*OpenVmHalo2VerifierFilterer, error) {
	contract, err := bindOpenVmHalo2Verifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OpenVmHalo2VerifierFilterer{contract: contract}, nil
}

// bindOpenVmHalo2Verifier binds a generic wrapper to an already deployed contract.
func bindOpenVmHalo2Verifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OpenVmHalo2VerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OpenVmHalo2Verifier.Contract.OpenVmHalo2VerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpenVmHalo2Verifier.Contract.OpenVmHalo2VerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OpenVmHalo2Verifier.Contract.OpenVmHalo2VerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OpenVmHalo2Verifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpenVmHalo2Verifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OpenVmHalo2Verifier.Contract.contract.Transact(opts, method, params...)
}

// Verify is a free data retrieval call binding the contract method 0x24270d54.
//
// Solidity: function verify(bytes publicValues, bytes proofData, bytes32 appExeCommit, bytes32 appVmCommit) view returns()
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierCaller) Verify(opts *bind.CallOpts, publicValues []byte, proofData []byte, appExeCommit [32]byte, appVmCommit [32]byte) error {
	var out []interface{}
	err := _OpenVmHalo2Verifier.contract.Call(opts, &out, "verify", publicValues, proofData, appExeCommit, appVmCommit)

	if err != nil {
		return err
	}

	return err

}

// Verify is a free data retrieval call binding the contract method 0x24270d54.
//
// Solidity: function verify(bytes publicValues, bytes proofData, bytes32 appExeCommit, bytes32 appVmCommit) view returns()
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierSession) Verify(publicValues []byte, proofData []byte, appExeCommit [32]byte, appVmCommit [32]byte) error {
	return _OpenVmHalo2Verifier.Contract.Verify(&_OpenVmHalo2Verifier.CallOpts, publicValues, proofData, appExeCommit, appVmCommit)
}

// Verify is a free data retrieval call binding the contract method 0x24270d54.
//
// Solidity: function verify(bytes publicValues, bytes proofData, bytes32 appExeCommit, bytes32 appVmCommit) view returns()
func (_OpenVmHalo2Verifier *OpenVmHalo2VerifierCallerSession) Verify(publicValues []byte, proofData []byte, appExeCommit [32]byte, appVmCommit [32]byte) error {
	return _OpenVmHalo2Verifier.Contract.Verify(&_OpenVmHalo2Verifier.CallOpts, publicValues, proofData, appExeCommit, appVmCommit)
}
