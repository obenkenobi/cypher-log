package topics

import "fmt"

func AppendDeadLetter(topic string) string {
	return topic + "-dead-letter"
}

func AppendRetry(topic string, index int) string {
	return fmt.Sprintf("%v-retry-%v", topic, index)
}
