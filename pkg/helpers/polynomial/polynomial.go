package polynomial

import (
	"crypto/rand"

	"filippo.io/edwards25519"
	"github.com/taurusgroup/frost-ed25519/pkg/helpers/common"
)

type Polynomial struct {
	coefficients []edwards25519.Scalar
}

// NewPolynomial generates a Polynomial f(X) = secret + a1*X + ... + at*X^t,
// with coefficients in Z_q, and degree t.
func NewPolynomial(degree uint32, constant *edwards25519.Scalar) *Polynomial {
	var polynomial Polynomial
	polynomial.coefficients = make([]edwards25519.Scalar, degree+1)

	// Set the constant term to the secret
	polynomial.coefficients[0].Set(constant)

	var randomBytes [64]byte
	for i := uint32(1); i <= degree; i++ {
		_, _ = rand.Read(randomBytes[:64])
		polynomial.coefficients[i].SetUniformBytes(randomBytes[:64])
	}

	return &polynomial
}

// Evaluate evaluates a polynomial in a given variable index
// We use Horner's method: https://en.wikipedia.org/wiki/Horner%27s_method
func (p *Polynomial) Evaluate(index uint32) *edwards25519.Scalar {
	var result, x edwards25519.Scalar
	//result := edwards25519.NewScalar()
	common.SetScalarUInt32(&x, index)
	//x := common.NewScalarUInt32(index)
	// revers order
	for i := len(p.coefficients) - 1; i >= 0; i-- {
		// b_n-1 = b_n * x + a_n-1
		result.MultiplyAdd(&result, &x, &p.coefficients[i])
	}
	return &result
}

// EvaluateMultiple performs evaluate on a slice of indices and returns the result as a map
func (p *Polynomial) EvaluateMultiple(indices []uint32) map[uint32]*edwards25519.Scalar {
	shares := make(map[uint32]*edwards25519.Scalar, len(indices))
	for _, index := range indices {
		shares[index] = p.Evaluate(index)
	}
	return shares
}

// Degree is the highest power of the Polynomial
func (p *Polynomial) Degree() uint32 {
	return uint32(len(p.coefficients)) - 1
}

// Size is the number of coefficients of the polynomial
// It is equal to Degree+1
func (p *Polynomial) Size() int {
	return len(p.coefficients)
}
