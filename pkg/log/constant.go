package log

const (
	BATCH_SIZE = 128
	fileLogRegexString = "^\\[(\\d{4}\\/[0-1][0-9]\\/[0-3][0-9]\\s[0-2][0-9]:[0-5][0-9]:[0-5][0-9].*?)\\]\\s(\\[[a-zA-Z0-9]+?\\])\\s(.*)$"
	commonLogRegexString = "^(\\d{4}\\/[0-1][0-9]\\/[0-3][0-9]\\s[0-2][0-9]:[0-5][0-9]:[0-5][0-9])\\s\\S+\\.go:\\d+:\\s(\\[[a-zA-Z0-9]+?\\])\\s.*$"
)