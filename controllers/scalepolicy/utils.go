package scalepolicy

import (
	"fmt"
	scalev1 "github.com/ShivamJha2436/kubehalo/api/v1"
)

func FormatScalePolicyID(sp *scalev1.ScalePolicy) string {
	return fmt.Sprintf("%s/%s", sp.Namespace, sp.Name)
}
