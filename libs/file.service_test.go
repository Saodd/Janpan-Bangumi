package libs

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileManager(t *testing.T) {
	log.SetFlags(0)
	fm := NewFileManager()
	fm.SetDirPath(filepath.Join("..", "bangumi"))
	yearReqs, err := fm.ReadAndParse()
	assert.NoError(t, err)

	fm.CheckDuplicated(yearReqs)
	t.Logf("解析完成，年份数量: %d\n", len(yearReqs))
	for _, yr := range yearReqs {
		t.Logf("年份: %s, 番剧数量: %d\n", yr.Year, len(yr.Bangumis))
		// for _, bg := range yr.Bangumis {
		// 	t.Logf("  - %s %s\n", bg.FileYear, bg.Title)
		// }
	}
}
