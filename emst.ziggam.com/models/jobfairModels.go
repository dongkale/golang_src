package models

type JobFairEntpRegister struct {
	RslData        DefaultResult
	RslJobFairList []JobfairInfo
}

type JobFairEntpDelete struct {
	RslData        DefaultResult
	RslJobFairList []JobfairInfo
}
