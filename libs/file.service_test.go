package libs

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileManager(t *testing.T) {
	log.SetFlags(0)
	fm := NewFileManager()
	yearReqs, err := fm.ReadAndParse()
	assert.NoError(t, err)

	fm.CheckDuplicated(yearReqs)
	t.Logf("解析完成，年份数量: %d\n", len(yearReqs))
	for _, yr := range yearReqs {
		t.Logf("年份: %s, 番剧数量: %d\n", yr.Year, len(yr.Bangumis))
	}
}
