package bmpkg

import (
	"errors"
	"fmt"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmbrick"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmsingleton/bmconf"
	"sync"
)

var t = make(map[string][]string)
var k []string
var oc sync.Once

func initEPipeline() {
	//TODO：ddsaas 待提出
	t["generatersakey"] = []string{"BmRsaKeyGenerateBrick"}
	t["getpublickey"] = []string{"BmGetPublicKeyBrick"}

	t["insertaccount"] = []string{"BmAccountPushBrick", "BmAccountGenerateToken"}
	t["accountlogin"] = []string{"BmAccountFindBrick", "BmAccountGenerateToken"}

	k = []string{
		"getpublickey", "generatersakey", "insertaccount", "accountlogin",
	}
}

func GetPkgLen(pkg string) (int, error) {
	oc.Do(initEPipeline)

	tmp := t[pkg]
	var err error
	if tmp == nil {
		err = errors.New("query resource router error")
	}

	return len(tmp), err
}

func GetCurBrick(pkg string, idx int64) (bmbrick.BMBrickFace, error) {

	oc.Do(initEPipeline)

	tmp := t[pkg]
	var err error
	if tmp == nil {
		err = errors.New("query resource router error")
	}

	reval := tmp[idx]
	fmt.Println(reval)
	if reval == "" {
		err = errors.New("query resource router error")
	}

	face, err := bmconf.GetBMBrick(reval)
	return face, err
}

func IsNeedAuth(pkg string, cur int64) bool {
	oc.Do(initEPipeline)
	for _, itm := range k {
		if itm == pkg {
			return false
		}
	}
	return cur == 0
}
