package accountfind

import (
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmbrick"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmconfig"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmrouter"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmsingleton/bmpkg"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/account"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/auth"
	"github.com/alfredyang1986/blackmirror/bmerror"
	"github.com/alfredyang1986/blackmirror/bmmodel/request"
	"github.com/alfredyang1986/blackmirror/jsonapi"
	"io"
	"net/http"
)

type BmAccountFindBrick struct {
	bk *bmbrick.BMBrick
}

/*------------------------------------------------
 * brick interface
 *------------------------------------------------*/

func (b *BmAccountFindBrick) Exec() error {
	var tmp account.BmAccount
	err := tmp.FindOne(*b.bk.Req)

	if err == nil && tmp.Account != "" {
		b.bk.Pr = tmp
		return nil
	}
	b.bk.Err = -11

	return err
}

func (b *BmAccountFindBrick) Prepare(pr interface{}) error {
	req := pr.(request.Request)

	var eqCondArr []request.Eqcond
	for _, e := range req.Eqcond {
		if e.Ky == "secretword" {
			tmpAccount := account.BmAccount{
				SecretWord: e.Vy.(string),
			}
			var bmRsa bmconfig.BmRsaConfig
			bmRsa.GenerateConfig()
			tmpAccount.DecodeByCompanyDate(bmRsa.Company, bmRsa.Date)
			tmpAccount.Secret2MD5()
			e.Vy = tmpAccount.SecretWord
		}
		eqCondArr = append(eqCondArr, e)
	}
	req.Eqcond = eqCondArr

	b.BrickInstance().Req = &req
	return nil
}

func (b *BmAccountFindBrick) Done(pkg string, idx int64, e error) error {
	ec := b.BrickInstance().Err
	if ec != 0 {
		return e
	}
	tmp, _ := bmpkg.GetPkgLen(pkg)
	if int(idx) < tmp-1 {
		bmrouter.NextBrickRemote(pkg, idx+1, b)
	}
	return nil
}

func (b *BmAccountFindBrick) BrickInstance() *bmbrick.BMBrick {
	if b.bk == nil {
		b.bk = &bmbrick.BMBrick{}
	}
	return b.bk
}

func (b *BmAccountFindBrick) ResultTo(w io.Writer) error {
	pr := b.BrickInstance().Pr
	tmp := pr.(account.BmAccount)
	err := jsonapi.ToJsonAPI(&tmp, w)
	return err
}

func (b *BmAccountFindBrick) Return(w http.ResponseWriter) {
	ec := b.BrickInstance().Err
	if ec != 0 {
		bmerror.ErrInstance().ErrorReval(ec, w)
	} else {
		reval := b.BrickInstance().Pr.(auth.BmLoginSucceed)
		jsonapi.ToJsonAPI(&reval, w)
	}
}
