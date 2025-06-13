package tasks

type TaskType string

const (
	TaskSummarize   TaskType = "summarize"
	TaskKeyPoints   TaskType = "key-points"
	TaskRisks       TaskType = "risks"
	TaskActionItems TaskType = "action-items"
)

var TaskPrompts = map[TaskType]string{
	TaskSummarize: `Provide a comprehensive summary of this document in the following format:

OVERVIEW:
[2-3 sentence high-level summary]

MAIN SECTIONS:
[Bullet points of major sections/topics]

KEY FINDINGS:
[Most important discoveries or statements]

CONCLUSION:
[Brief wrap-up of document's purpose and outcome]`,

	TaskKeyPoints: `Extract and list the most important points from this document:

1. [First key point with brief explanation]
2. [Second key point with brief explanation]
... (continue for all major points)

Focus on actionable information, critical deadlines, important numbers, and binding commitments.`,

	TaskRisks: `Analyze this document for potential risks, concerns, or red flags:

LEGAL RISKS:
- [Any legal vulnerabilities or unclear terms]

FINANCIAL RISKS:
- [Financial obligations or exposures]

OPERATIONAL RISKS:
- [Process or execution challenges]

COMPLIANCE RISKS:
- [Regulatory or policy concerns]

RECOMMENDATIONS:
- [Suggested mitigations for identified risks]`,

	TaskActionItems: `Extract all action items and next steps from this document:

IMMEDIATE ACTIONS (Within 7 days):
□ [Action item with responsible party if mentioned]
□ [Action item with deadline if specified]

SHORT-TERM ACTIONS (Within 30 days):
□ [Action items]

LONG-TERM ACTIONS (30+ days):
□ [Action items]

DEPENDENCIES:
- [Items that require other actions to complete first]`,
}

var TaskLabels = map[TaskType]string{
	TaskSummarize:   "📝 Summary",
	TaskKeyPoints:   "🎯 Key Points",
	TaskRisks:       "⚠️ Risk Analysis",
	TaskActionItems: "✅ Action Items",
}

func GetTaskPrompt(taskType TaskType) string {
	if prompt, exists := TaskPrompts[taskType]; exists {
		return prompt
	}
	return "Analyze this document and provide insights."
}

func GetTaskLabel(taskType TaskType) string {
	if label, exists := TaskLabels[taskType]; exists {
		return label
	}
	return "📄 Analysis"
}

func GetAllTaskTypes() []TaskType {
	return []TaskType{
		TaskSummarize,
		TaskKeyPoints,
		TaskRisks,
		TaskActionItems,
	}
}

func IsValidTaskType(taskType string) bool {
	_, exists := TaskPrompts[TaskType(taskType)]
	return exists
}