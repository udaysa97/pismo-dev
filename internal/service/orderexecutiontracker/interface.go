package orderexecutiontracker

type OrderExecutionTrackerInterface interface {
	ProcessJobEvent(message []byte) bool
	SetRequiredRepos(repos RequiredRepos)
	InitJobTrackerConsumer()
}
