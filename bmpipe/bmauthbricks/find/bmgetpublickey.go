package rsafind

import (
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmbrick"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmrouter"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmsingleton/bmpkg"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/auth"
	"github.com/alfredyang1986/blackmirror/bmmodel/request"
	"github.com/alfredyang1986/blackmirror/bmerror"
	"github.com/alfredyang1986/blackmirror/jsonapi"
	"net/http"
	"io"
)

type BmGetPublicKeyBrick struct {
	bk *bmbrick.BMBrick
}

/*------------------------------------------------
 * brick interface
 *------------------------------------------------*/

func (b *BmGetPublicKeyBrick) Exec() error {
	var tmp auth.BmRsaKey
	err := tmp.FindOne(*b.bk.Req)
	tmp.PrivateKey = ""
	b.bk.Pr = tmp
	return err
}

func (b *BmGetPublicKeyBrick) Prepare(pr interface{}) error {
	req := pr.(request.Request)
	b.BrickInstance().Req = &req
	return nil
}

func (b *BmGetPublicKeyBrick) Done(pkg string, idx int64, e error) error {
	tmp, _ := bmpkg.GetPkgLen(pkg)
	if int(idx) < tmp-1 {
		bmrouter.NextBrickRemote(pkg, idx+1, b)
	}
	return nil
}

func (b *BmGetPublicKeyBrick) BrickInstance() *bmbrick.BMBrick {
	if b.bk == nil {
		b.bk = &bmbrick.BMBrick{}
	}
	return b.bk
}

func (b *BmGetPublicKeyBrick) ResultTo(w io.Writer) error {
	pr := b.BrickInstance().Pr
	tmp := pr.(auth.BmRsaKey)
	err := jsonapi.ToJsonAPI(&tmp, w)
	return err
}

func (b *BmGetPublicKeyBrick) Return(w http.ResponseWriter) {
	ec := b.BrickInstance().Err
	if ec != 0 {
		bmerror.ErrInstance().ErrorReval(ec, w)
	} else {
		var reval auth.BmRsaKey = b.BrickInstance().Pr.(auth.BmRsaKey)
		jsonapi.ToJsonAPI(&reval, w)
	}
}

