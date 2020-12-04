// Copyright 2020 The Reed Developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package pow

import (
	"fmt"
	"github.com/reed/types"
	"testing"
)

func TestCalcNextDifficulty(t *testing.T) {
	oldDiff := DifficultyLimit()

	p := &types.Block{
	}
	p.Timestamp = 1606458831
	p.BigNumber = oldDiff

	diff := CalcDifficulty(1605249231, p)

	if diff.Cmp(&oldDiff) == -1 {
		fmt.Println("YES")
	} else {
		t.Error("difficulty error")
	}

}

func TestDifficultyLimit(t *testing.T) {
	l := DifficultyLimit()
	fmt.Println(l)
}

func TestDifficultyAdjustmentInterval(t *testing.T) {
	if DifficultyAdjustmentInterval() != 2016 {
		t.Error("DifficultyAdjustmentInterval != 2016")
	}
}
