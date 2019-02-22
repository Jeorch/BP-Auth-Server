package accountother

import (
	"crypto/md5"
	"fmt"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmbrick"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmoauth"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmrouter"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/account"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/auth"
	"github.com/alfredyang1986/blackmirror/bmcommon/bmsingleton/bmpkg"
	"github.com/alfredyang1986/blackmirror/bmerror"
	"github.com/alfredyang1986/blackmirror/jsonapi"
	"io"
	"net/http"
)

type BmAccountGenerateToken struct {
	bk *bmbrick.BMBrick
}

/*------------------------------------------------
 * brick interface
 *------------------------------------------------*/

func (b *BmAccountGenerateToken) Exec() error {
	var err error
	tmp := b.bk.Pr.(account.BmAccount)

	h := md5.New()
	io.WriteString(h, tmp.Id)
	token := fmt.Sprintf("%x", h.Sum(nil))
	err = bmoauth.PushToken(token)

	bmls := auth.BmLoginSucceed{
		Id:    tmp.Id,
		Id_:   tmp.Id_,
		Token: token,
	}

	b.BrickInstance().Pr = bmls
	return err
}

func (b *BmAccountGenerateToken) Prepare(pr interface{}) error {
	req := pr.(account.BmAccount)
	b.BrickInstance().Pr = req
	return nil
}

func (b *BmAccountGenerateToken) Done(pkg string, idx int64, e error) error {
	tmp, _ := bmpkg.GetPkgLen(pkg)
	ec := b.BrickInstance().Err
	if int(idx) < tmp-1 && ec == 0 {
		bmrouter.NextBrickRemote(pkg, idx+1, b)
	}
	return nil
}

func (b *BmAccountGenerateToken) BrickInstance() *bmbrick.BMBrick {
	if b.bk == nil {
		b.bk = &bmbrick.BMBrick{}
	}
	return b.bk
}

func (b *BmAccountGenerateToken) ResultTo(w io.Writer) error {
	pr := b.BrickInstance().Pr
	tmp := pr.(account.BmAccount)
	err := jsonapi.ToJsonAPI(&tmp, w)
	return err
}

func (b *BmAccountGenerateToken) Return(w http.ResponseWriter) {
	ec := b.BrickInstance().Err
	if ec != 0 {
		bmerror.ErrInstance().ErrorReval(ec, w)
	} else {
		reval := b.BrickInstance().Pr.(auth.BmLoginSucceed)
		jsonapi.ToJsonAPI(&reval, w)
	}
}
