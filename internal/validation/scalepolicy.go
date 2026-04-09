package validation

import (
	"fmt"
	"slices"
	"time"

	kubehalov1 "github.com/ShivamJha2436/kubehalo/api/kubehalo/v1"
)

var validScheduleDays = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

// ValidateScalePolicy performs structural and logical validation shared by runtime components.
func ValidateScalePolicy(policy *kubehalov1.ScalePolicy) error {
	switch {
	case policy.Spec.TargetRef.Kind == "":
		return fmt.Errorf("spec.targetRef.kind must not be empty")
	case policy.Spec.TargetRef.Name == "":
		return fmt.Errorf("spec.targetRef.name must not be empty")
	case policy.Spec.TargetRef.Namespace == "":
		return fmt.Errorf("spec.targetRef.namespace must not be empty")
	case policy.Spec.Metric.Query == "":
		return fmt.Errorf("spec.metric.query must not be empty")
	case policy.Spec.Metric.Threshold < 0:
		return fmt.Errorf("spec.metric.threshold must be non-negative")
	case policy.Spec.MinReplicas <= 0:
		return fmt.Errorf("spec.minReplicas must be greater than zero")
	case policy.Spec.MaxReplicas < policy.Spec.MinReplicas:
		return fmt.Errorf("spec.maxReplicas must be greater than or equal to spec.minReplicas")
	case policy.Spec.ScaleUp.Step <= 0:
		return fmt.Errorf("spec.scaleUp.step must be greater than zero")
	case policy.Spec.ScaleDown.Step <= 0:
		return fmt.Errorf("spec.scaleDown.step must be greater than zero")
	}

	for i, schedule := range policy.Spec.Schedules {
		if err := validateSchedule(i, schedule, policy.Spec.MinReplicas, policy.Spec.MaxReplicas); err != nil {
			return err
		}
	}

	for i := 0; i < len(policy.Spec.Schedules); i++ {
		for j := i + 1; j < len(policy.Spec.Schedules); j++ {
			if schedulesOverlap(policy.Spec.Schedules[i], policy.Spec.Schedules[j]) {
				return fmt.Errorf("spec.schedules[%d] overlaps with spec.schedules[%d]", i, j)
			}
		}
	}

	return nil
}

func validateSchedule(index int, schedule kubehalov1.ScheduleSpec, defaultMinReplicas, defaultMaxReplicas int32) error {
	if len(schedule.Days) == 0 {
		return fmt.Errorf("spec.schedules[%d].days must not be empty", index)
	}

	start, err := parseClock(schedule.StartTime)
	if err != nil {
		return fmt.Errorf("spec.schedules[%d].startTime must be in HH:MM format", index)
	}
	end, err := parseClock(schedule.EndTime)
	if err != nil {
		return fmt.Errorf("spec.schedules[%d].endTime must be in HH:MM format", index)
	}
	if start >= end {
		return fmt.Errorf("spec.schedules[%d].startTime must be before endTime", index)
	}

	seenDays := make(map[string]struct{}, len(schedule.Days))
	for _, day := range schedule.Days {
		if !slices.Contains(validScheduleDays, day) {
			return fmt.Errorf("spec.schedules[%d].days contains invalid day %q", index, day)
		}
		if _, exists := seenDays[day]; exists {
			return fmt.Errorf("spec.schedules[%d].days contains duplicate day %q", index, day)
		}
		seenDays[day] = struct{}{}
	}

	minReplicas := schedule.MinReplicas
	maxReplicas := schedule.MaxReplicas
	if minReplicas == 0 {
		minReplicas = defaultMinReplicas
	}
	if maxReplicas == 0 {
		maxReplicas = defaultMaxReplicas
	}
	if minReplicas <= 0 {
		return fmt.Errorf("spec.schedules[%d].minReplicas must be greater than zero", index)
	}
	if maxReplicas < minReplicas {
		return fmt.Errorf("spec.schedules[%d].maxReplicas must be greater than or equal to minReplicas", index)
	}

	return nil
}

func schedulesOverlap(left, right kubehalov1.ScheduleSpec) bool {
	if !daysIntersect(left.Days, right.Days) {
		return false
	}

	leftStart, err := parseClock(left.StartTime)
	if err != nil {
		return false
	}
	leftEnd, err := parseClock(left.EndTime)
	if err != nil {
		return false
	}
	rightStart, err := parseClock(right.StartTime)
	if err != nil {
		return false
	}
	rightEnd, err := parseClock(right.EndTime)
	if err != nil {
		return false
	}

	return leftStart < rightEnd && rightStart < leftEnd
}

func daysIntersect(left, right []string) bool {
	rightSet := make(map[string]struct{}, len(right))
	for _, day := range right {
		rightSet[day] = struct{}{}
	}
	for _, day := range left {
		if _, ok := rightSet[day]; ok {
			return true
		}
	}
	return false
}

func parseClock(value string) (int, error) {
	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 0, err
	}
	return parsed.Hour()*60 + parsed.Minute(), nil
}
