package imgUtil

import "github.com/spf13/viper"

func IsValidScalingOption(o int) bool {
	for _, it := range ScalingOptions() {
		if it == o {
			return true
		}
	}
	return false
}

func ScalingOptions() []int {
	return viper.GetIntSlice("scale.options")
}
