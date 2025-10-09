// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title IOpenVmHalo2Verifier
 * @notice Interface for OpenVM Halo2 Verifier contract
 * @dev Based on openvm-solidity-sdk v1.4
 */
interface IOpenVmHalo2Verifier {
    /**
     * @notice Verify an OpenVM proof
     * @param publicValues The public values as bytes
     * @param proofData The proof data bytes
     * @param appExeCommit The application execution commitment
     * @param appVmCommit The application VM commitment
     * @dev Reverts if verification fails
     */
    function verify(
        bytes calldata publicValues,
        bytes calldata proofData,
        bytes32 appExeCommit,
        bytes32 appVmCommit
    ) external view;
}

