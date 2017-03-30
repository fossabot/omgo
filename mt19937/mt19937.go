/*
MIT License

Copyright (c) 2017 Master.G

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package mt19937

import "math"

const (
	MT_N = 624
	M    = 397
	// Constant vector A
	MATRIX_A = 0x9908B0DF
	// Most significant w-r bits
	UPPER_MASK = 0x80000000
	// least significant w-r bits
	LOWER_MASK = 0x7FFFFFFF

	FULL_MASK = 0xFFFFFFFF
	DEF_SEED  = 0x012BD6AA
)

type (
	Context struct {
		// State vector
		mt [MT_N]uint32
		// mti == N + 1 -> mt[N] not initialized
		mti uint32
	}
)

var (
	defaultMT [MT_N]uint32
	mag01     [2]uint32
)

func init() {
	var i uint32
	for i = 1; i < MT_N; i++ {
		// See Knuth TAOCP Vol2. 3rd Ed. P.106 for multiplier.
		defaultMT[i] = 1812433253*(defaultMT[i-1]^(defaultMT[i-1]>>30)) + i
		defaultMT[i] &= FULL_MASK
	}
	mag01[0] = 0
	mag01[1] = MATRIX_A
}

// NewContext create a new MT19937 RNG context with a given uint32 seed
func NewContext(seed uint32) *Context {
	ctx := new(Context)
	ctx.init(seed)
	return ctx
}

// init MT19937 context with a given uint32 seed
func (ctx *Context) init(seed uint32) {
	copy(ctx.mt, defaultMT)
	ctx.mt[0] = seed & FULL_MASK
}

// NewContextWithArray create a new MT19937 RNG context with a given uint32 array seed
func NewContextWithArray(seed []uint32) *Context {
	ctx := NewContext(DEF_SEED)
	var i, j, k uint32
	if MT_N > len(seed) {
		k = MT_N
	} else {
		k = uint32(len(seed))
	}

	for ; k > 0; i-- {
		ctx.mt[i] = (ctx.mt[i] ^ ((ctx.mt[i-1] ^ (ctx.mt[i-1] >> 30)) * 1664525)) + seed[j] + j
		ctx.mt[i] &= FULL_MASK
		i++
		j++
		if j >= MT_N {
			ctx.mt[0] = ctx.mt[MT_N-1]
			i = 1
		}
		if j >= uint32(len(seed)) {
			j = 0
		}
	}

	for k = MT_N - 1; k > 0; k-- {
		ctx.mt[i] = (ctx.mt[i] ^ ((ctx.mt[i-1] ^ (ctx.mt[i-1] >> 30)) * 1566083941)) - i
		ctx.mt[i] &= FULL_MASK
		i++

		if i >= MT_N {
			ctx.mt[0] = ctx.mt[MT_N-1]
			i = 1
		}
	}

	ctx.mt[0] = UPPER_MASK

	return ctx
}

// NextUInt32 generates the next pseudorandom uint32 number
func (ctx *Context) NextUInt32() uint32 {
	var y uint32
	var kk int
	// mag01[x] = x * MATRIX_A for x = 0, 1
	if ctx.mti >= MT_N {
		if ctx.mti == MT_N+1 {
			ctx.init(5489)
		}

		for kk = 0; kk < MT_N-M; kk++ {
			y = (ctx.mt[kk] & UPPER_MASK) | (ctx.mt[kk+1] & LOWER_MASK)
			ctx.mt[kk] = ctx.mt[kk+M] ^ (y >> 1) ^ mag01[y&0x1]
		}

		for ; kk < MT_N-1; kk++ {
			y = (ctx.mt[kk] & UPPER_MASK) | (ctx.mt[kk+1] & LOWER_MASK)
			ctx.mt[kk] = ctx.mt[kk+(M-MT_N)] ^ (y >> 1) ^ mag01[y&0x1]
		}

		y = (ctx.mt[MT_N-1] & UPPER_MASK) | (ctx.mt[0] & LOWER_MASK)
		ctx.mt[MT_N-1] = ctx.mt[M-1] ^ (y >> 1) ^ mag01[y&0x1]
		ctx.mti = 0
	}
	y = ctx.mt[ctx.mti]
	ctx.mti++
	// Tempering
	y ^= y >> 11
	y ^= (y << 7) & 0x9D2C5680
	y ^= (y << 15) & 0xEFC60000
	y ^= y >> 18

	return y
}

// NextInt32 generates the next pseudorandom int32 number
func (ctx *Context) NextInt32() int32 {
	return int32(ctx.NextUInt32())
}

// NextInt generates the next pseudorandom 32bit int number
func (ctx *Context) NextInt() int {
	return int(ctx.NextUInt32())
}

// NextInt generates the next pseudorandom float32 number, between [0.0 ~ 1.0)
func (ctx *Context) NextFloat32() float32 {
	return float32(ctx.NextUInt32()) / float32(math.MaxUint32+1)
}

// NextInt generates the next pseudorandom float64 number, between [0.0 ~ 1.0)
func (ctx *Context) NextFloat64() float64 {
	return float64(ctx.NextUInt32()) / float64(math.MaxUint32+1)
}
