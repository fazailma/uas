package models

// AchievementStatistics represents comprehensive statistics for achievements
type AchievementStatistics struct {
	TotalByType        []TypeStatistic   `json:"total_by_type"`
	TotalByPeriod      []PeriodStatistic `json:"total_by_period"`
	TopStudents        []StudentRanking  `json:"top_students"`
	CompetitionDistrib []CompetitionStat `json:"competition_distribution"`
	Summary            StatisticsSummary `json:"summary"`
}

// TypeStatistic represents count of achievements per type/category
type TypeStatistic struct {
	Type  string  `json:"type"`
	Count int64   `json:"count"`
	Pct   float64 `json:"percentage"`
}

// PeriodStatistic represents count of achievements per period (month/year)
type PeriodStatistic struct {
	Period string `json:"period"` // Format: YYYY-MM
	Count  int64  `json:"count"`
}

// StudentRanking represents top performing students
type StudentRanking struct {
	StudentID  string `json:"student_id"`
	FullName   string `json:"full_name"`
	TotalCount int64  `json:"total_achievements"`
	Verified   int64  `json:"verified_count"`
	Pending    int64  `json:"pending_count"`
}

// CompetitionStat represents distribution of achievement competition levels
type CompetitionStat struct {
	Level string `json:"level"` // Lokal, Regional, Nasional, Internasional
	Count int64  `json:"count"`
}

// StatisticsSummary provides overall achievement statistics
type StatisticsSummary struct {
	TotalAchievements     int64   `json:"total_achievements"`
	VerifiedCount         int64   `json:"verified_count"`
	PendingCount          int64   `json:"pending_count"`
	RejectedCount         int64   `json:"rejected_count"`
	VerificationRate      float64 `json:"verification_rate"` // Percentage of verified
	TotalStudentsInvolved int64   `json:"total_students_involved"`
}
