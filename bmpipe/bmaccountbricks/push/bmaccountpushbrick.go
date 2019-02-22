package accountpush

import (
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmbrick"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmconfig"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmrouter"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmsingleton/bmpkg"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/account"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/auth"
	"github.com/alfredyang1986/blackmirror/bmerror"
	"github.com/alfredyang1986/blackmirror/jsonapi"
	"io"
	"net/http"
)

type BmAccountPushBrick struct {
	bk *bmbrick.BMBrick
}

/*------------------------------------------------
 * brick interface
 *------------------------------------------------*/

func (b *BmAccountPushBrick) Exec() error {
	var err error
	var tmp account.BmAccount = b.bk.Pr.(account.BmAccount)

	if tmp.Id != "" && tmp.Id_.Valid() {
		if tmp.Valid() && tmp.IsAccountRegisted() {
			//TODO: error处理
			b.bk.Err = -8
		} else {
			//TODO: 参数 配置文件化【以年为周期更新rsa】
			var bmRsa bmconfig.BmRsaConfig
			bmRsa.GenerateConfig()
			err = tmp.DecodeByCompanyDate(bmRsa.Company, bmRsa.Date)
			if err != nil {
				return err
			}
			tmp.Secret2MD5()

			tmp.InsertBMObject()
			tmp.SecretWord = ""
			b.bk.Pr = tmp
		}
	}

	return err
}

func (b *BmAccountPushBrick) Prepare(pr interface{}) error {

	req := pr.(account.BmAccount)
	b.BrickInstance().Pr = req
	return nil
}

func (b *BmAccountPushBrick) Done(pkg string, idx int64, e error) error {
	tmp, _ := bmpkg.GetPkgLen(pkg)
	ec := b.BrickInstance().Err
	if int(idx) < tmp-1 && ec == 0 {
		bmrouter.NextBrickRemote(pkg, idx+1, b)
	}
	return nil
}

func (b *BmAccountPushBrick) BrickInstance() *bmbrick.BMBrick {
	if b.bk == nil {
		b.bk = &bmbrick.BMBrick{}
	}
	return b.bk
}

func (b *BmAccountPushBrick) ResultTo(w io.Writer) error {
	pr := b.BrickInstance().Pr
	tmp := pr.(account.BmAccount)
	err := jsonapi.ToJsonAPI(&tmp, w)
	return err
}

func (b *BmAccountPushBrick) Return(w http.ResponseWriter) {
	ec := b.BrickInstance().Err
	if ec != 0 {
		bmerror.ErrInstance().ErrorReval(ec, w)
	} else {
		reval := b.BrickInstance().Pr.(auth.BmLoginSucceed)
		jsonapi.ToJsonAPI(&reval, w)
	}
}
