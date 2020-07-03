package log

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestFileLogRegex(t *testing.T) {
	r, err := regexp.Compile(fileLogRegexString)
	if err != nil {
		t.Error(err)
	}

	log := `[2019/09/27 15:45:50.651 +08:00] [INFO] [patchouli] [knowledge] [mahou] [ga] [umaretahi]`
	submatch := r.FindStringSubmatch(log)

	assert.Equal(t, len(submatch), 4, "there should be 4 sub matches")
	assert.Equal(t, submatch[1], "2019/09/27 15:45:50.651 +08:00", "match log datetime")
	assert.Equal(t, submatch[2], "[INFO]", "match log level")
	assert.Equal(t, submatch[3], "[patchouli] [knowledge] [mahou] [ga] [umaretahi]", "match log content")

	log = `[2019/09/29 12:00:30.468 +00:00] [info] [util.go:58] [pd] [git-hash=e8f4f0084407b53ac382d4998e30759b0e4b9c8e]`
	submatch = r.FindStringSubmatch(log)

	assert.Equal(t, len(submatch), 4, "there should be 4 sub matches")
	assert.Equal(t, submatch[1], "2019/09/29 12:00:30.468 +00:00", "match log datetime")
	assert.Equal(t, submatch[2], "[info]", "match log level")
	assert.Equal(t, submatch[3], "[util.go:58] [pd] [git-hash=e8f4f0084407b53ac382d4998e30759b0e4b9c8e]", "match log content")

	nomatch := r.FindStringSubmatch(`2019/09/29 11:45:14 executor.go:34: [info] job 2 finished`)
	assert.Equal(t, len(nomatch), 0, "should not match")
}

func TestCommonLogRegex(t *testing.T) {
	r, err := regexp.Compile(commonLogRegexString)
	if err != nil {
		t.Error(err)
	}

	log := `2019/09/29 11:45:14 executor.go:34: [info] job 2 finished`
	submatch := r.FindStringSubmatch(log)

	assert.Equal(t, len(submatch), 3, "there should be 3 sub matches")
	assert.Equal(t, submatch[1], "2019/09/29 11:45:14", "match log datetime")
	assert.Equal(t, submatch[2], "[info]", "match log level")
	assert.Equal(t, submatch[0], "2019/09/29 11:45:14 executor.go:34: [info] job 2 finished")

	nomatch := r.FindStringSubmatch( `[2019/09/27 15:45:50.651 +08:00] [INFO] [patchouli] [knowledge] [mahou] [ga] [umaretahi]`)
	assert.Equal(t, len(nomatch), 0, "should not match")
}
