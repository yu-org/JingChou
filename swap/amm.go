package swap

import (
	"errors"
	"math/big"
	"sync"
)

// --------------------------
// 常量与辅助
// --------------------------
var (
	ONE             = big.NewInt(1)
	DefaultFeeNum   = big.NewInt(997) // numerator for fee (997/1000 -> 0.3%)
	DefaultFeeDen   = big.NewInt(1000)
	Zero            = big.NewInt(0)
	ErrInsufficient = errors.New("insufficient liquidity")
)

// ceilDiv: ceil(a/b) for big.Int (returns floor((a + b -1)/b))
func ceilDiv(a, b *big.Int) *big.Int {
	q := new(big.Int).Div(a, b)
	r := new(big.Int).Mod(a, b)
	if r.Cmp(Zero) > 0 {
		q.Add(q, ONE)
	}
	return q
}

type Pool struct {
	mu       sync.Mutex
	Reserve0 *big.Int // token0 reserve
	Reserve1 *big.Int // token1 reserve
	FeeNum   *big.Int // fee numerator (e.g., 997)
	FeeDen   *big.Int // fee denominator (e.g., 1000)
}

// NewPool 创建一个新池（初始储备可以为 0）
func NewPool(res0, res1 *big.Int, feeNum, feeDen *big.Int) *Pool {
	// copy inputs to avoid aliasing
	r0 := new(big.Int).Set(res0)
	r1 := new(big.Int).Set(res1)
	fn := new(big.Int).Set(feeNum)
	fd := new(big.Int).Set(feeDen)
	return &Pool{
		Reserve0: r0,
		Reserve1: r1,
		FeeNum:   fn,
		FeeDen:   fd,
	}
}

// AddLiquidity 简单把 token0/token1 加到储备里（不计算 LP token）
func (p *Pool) AddLiquidity(amount0, amount1 *big.Int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Reserve0.Add(p.Reserve0, new(big.Int).Set(amount0))
	p.Reserve1.Add(p.Reserve1, new(big.Int).Set(amount1))
}

// GetAmountOut (Uniswap V2 公式):
// amountOut = (amountIn * feeNum * reserveOut) / (reserveIn*feeDen + amountIn*feeNum)
func (p *Pool) GetAmountOut(amountIn, reserveIn, reserveOut *big.Int) *big.Int {
	if amountIn.Cmp(Zero) <= 0 {
		return big.NewInt(0)
	}
	if reserveIn.Cmp(Zero) <= 0 || reserveOut.Cmp(Zero) <= 0 {
		return big.NewInt(0)
	}
	amountInWithFee := new(big.Int).Mul(amountIn, p.FeeNum)                           // amountIn * feeNum
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)                        // amountInWithFee * reserveOut
	denom := new(big.Int).Add(new(big.Int).Mul(reserveIn, p.FeeDen), amountInWithFee) // reserveIn*feeDen + amountInWithFee
	return new(big.Int).Div(numerator, denom)
}

// GetAmountIn (Uniswap V2 逆推，向上取整):
// amountIn = ceil( reserveIn * amountOut * feeDen / ((reserveOut - amountOut) * feeNum) )
func (p *Pool) GetAmountIn(amountOut, reserveIn, reserveOut *big.Int) (*big.Int, error) {
	if amountOut.Cmp(Zero) <= 0 {
		return big.NewInt(0), nil
	}
	if reserveIn.Cmp(Zero) <= 0 || reserveOut.Cmp(Zero) <= 0 {
		return nil, ErrInsufficient
	}
	if amountOut.Cmp(reserveOut) >= 0 {
		return nil, ErrInsufficient
	}
	num := new(big.Int).Mul(reserveIn, amountOut)  // reserveIn * amountOut
	num.Mul(num, p.FeeDen)                         // * feeDen
	den := new(big.Int).Sub(reserveOut, amountOut) // reserveOut - amountOut
	den.Mul(den, p.FeeNum)                         // * feeNum
	if den.Cmp(Zero) == 0 {
		return nil, errors.New("division by zero in GetAmountIn")
	}
	// ceil division:
	amountIn := ceilDiv(num, den)
	return amountIn, nil
}

// Swap0For1 用户给 amountIn token0，池子返回 amountOut token1，并更新储备
func (p *Pool) Swap0For1(amountIn *big.Int) (*big.Int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if amountIn.Cmp(Zero) <= 0 {
		return big.NewInt(0), nil
	}
	// compute amountOut
	amountOut := p.GetAmountOut(amountIn, p.Reserve0, p.Reserve1)
	if amountOut.Cmp(Zero) == 0 {
		return big.NewInt(0), nil
	}
	if amountOut.Cmp(p.Reserve1) > 0 {
		return nil, ErrInsufficient
	}
	// update reserves: x += amountIn, y -= amountOut
	p.Reserve0.Add(p.Reserve0, new(big.Int).Set(amountIn))
	p.Reserve1.Sub(p.Reserve1, new(big.Int).Set(amountOut))
	return amountOut, nil
}

// Swap1For0 用户给 amountIn token1，池子返回 amountOut token0，并更新储备
func (p *Pool) Swap1For0(amountIn *big.Int) (*big.Int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if amountIn.Cmp(Zero) <= 0 {
		return big.NewInt(0), nil
	}
	amountOut := p.GetAmountOut(amountIn, p.Reserve1, p.Reserve0)
	if amountOut.Cmp(Zero) == 0 {
		return big.NewInt(0), nil
	}
	if amountOut.Cmp(p.Reserve0) > 0 {
		return nil, ErrInsufficient
	}
	p.Reserve1.Add(p.Reserve1, new(big.Int).Set(amountIn))
	p.Reserve0.Sub(p.Reserve0, new(big.Int).Set(amountOut))
	return amountOut, nil
}

// Inspect 返回当前储备的拷贝
func (p *Pool) Inspect() (r0, r1 *big.Int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return new(big.Int).Set(p.Reserve0), new(big.Int).Set(p.Reserve1)
}
