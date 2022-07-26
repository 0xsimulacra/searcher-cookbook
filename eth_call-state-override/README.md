# eth_call state override

As you maybe know you can use the third param in eth_call to override the state while you call the endpoint.  
This opens some possibilities, like query the chain or write a smart contract you don't need to deploy.  

In the example you will find an example that overwrites an uni2 pair to query in bulk to get the current reserves.

### bulk call to get reserves
In the bulk-call/bulk-call.go you will find an example of how you can make a bulk call to eth_call and get the reserves from the pairs.  
The advantage of a bulk-call is, that you can run multiple smart contract calls with only one call to the node.  
Use this as an improvement to reduce calls in your bot ;)

```go 
go run bulk-call/bulk-call.go --rpc _wss_connection_string_
```

### explanation
```
//SPDX-License-Identifier: UNLICENSED
pragma solidity = 0.8.15;

contract Uni2ReserveInfo {

    function getReserves() public view returns (address pair, uint112 reserve0, uint112 reserve1) {
        assembly {
            let reserveHolder := sload(8)
            reserve0 := shr(144, shl(32, reserveHolder))
            reserve1 := shr(144, shl(144, reserveHolder))
        }
        pair = address(this);
    }
}
```
is our fake contract we created  

normally the uni2 pair contract returns  
_reserve0 uint112, _reserve1 uint112, _blockTimestampLast uint32    

to help summarize the data for later processing we want the following return values  
address pair, uint112 reserve0, uint112 reserve1

after you have your custom contract you need the deployed byte-code to use it to overwrite the state of the original pair contract  
```bash
solc Xyz.sol --bin-runtime --optimize
```
if you don't have solc, you can use remix, deploy the contract and you will find it in "Compilation Details" => RUNTIME BYTECODE

### Tips
Of course, you don't need to use assembly in your fake smart contract  
Here is an example without assembly
```
//SPDX-License-Identifier: UNLICENSED
pragma solidity = 0.6.12;

contract Uni3ReserveInfo {
    function getReserves(address pair) public view returns (address pairData, uint160 sqrtPriceX96, int24 tick, uint128 liquidity) {
        pairData = pair;
        (sqrtPriceX96, tick,,,,,) = IUniswapV3PoolData(pair).slot0();
        liquidity = IUniswapV3PoolData(pair).liquidity();
    }
}

interface IUniswapV3PoolData {
    function liquidity() external view returns (uint128);
    function slot0() external view returns (uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked);
}
```

Generally you can also just use the BatchCall without an overwrite to reduce calls. 