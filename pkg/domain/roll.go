package domain

import (
	"math/rand"

	"github.com/Knetic/govaluate"
	"github.com/pkg/errors"
)

type Roll struct {
	ID     string `db:"id,omitempty" json:"id,omitempty"`
	GameID string `db:"game_id,omitempty" json:"game_id,omitempty"`

	DieExpressions []int `db:"die_expressions,omitempty" json:"die_expressions,omitempty"`
	DieResults     []int `db:"die_results,omitempty" json:"die_results,omitempty"`

	ModifierExpression string  `db:"modifier_expression,omitempty" json:"modifier_expression,omitempty"`
	ModifierResult     float64 `db:"modifier_result,omitempty" json:"modifier_result,omitempty"`
}

func (r *Roll) CalculateRoll() error {

	modifierExpr, err := govaluate.NewEvaluableExpression(r.ModifierExpression)
	if err != nil {
		return errors.Wrap(err, "cannot parse modifier expression")
	}

	modifierRes, err := modifierExpr.Eval(nil)
	if err != nil {
		return errors.Wrap(err, "cannot evaluate modifier expression")
	}

	if modifierRes, ok := modifierRes.(float64); ok {
		r.ModifierResult = modifierRes
	} else {
		return errors.New("modifier result dose not resemble float64")
	}

	if r.DieExpressions == nil {
		r.DieExpressions = []int{}
	}

	r.DieResults = []int{}
	for _, dieExpr := range r.DieExpressions {
		r.DieResults = append(r.DieResults, rand.Intn(dieExpr)+1)
	}

	return nil
}
