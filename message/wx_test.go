package message

import "testing"
import "github.com/lvfeiyang/xiaobai/common/flog"

func TestGenerateNoncestr(t *testing.T) {
	flog.Init()
	noncestr := generateNoncestr(16)
	flog.LogFile.Println("nonce string:", noncestr)
	if 16 != len(noncestr) {
		t.Error("generate noncestr len err")
	}
}
