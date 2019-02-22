package rsapush

import (
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmbrick"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmrouter"
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmsingleton/bmpkg"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/auth"
	"github.com/alfredyang1986/blackmirror/bmsecurity"
	"github.com/alfredyang1986/blackmirror/bmerror"
	"github.com/alfredyang1986/blackmirror/jsonapi"
	"io"
	"net/http"
)

type BmRsaKeyGenerateBrick struct {
	bk *bmbrick.BMBrick
}

/*------------------------------------------------
 * brick interface
 *------------------------------------------------*/

func (b *BmRsaKeyGenerateBrick) Exec() error {
	//TODO:Chech Auth and check company&date exist?
	rsaKey, err := bmsecurity.GetRsaKey(512)
	if err != nil {
		return err
	}
	tmp := b.bk.Pr.(auth.BmRsaKey)
	tmp.PublicKey = rsaKey.PublicKey
	tmp.PrivateKey = rsaKey.PrivateKey
	tmp.InsertBMObject()
	tmp.PrivateKey = ""
	b.BrickInstance().Pr = tmp
	return nil
}

func (b *BmRsaKeyGenerateBrick) Prepare(pr interface{}) error {
	req := pr.(auth.BmRsaKey)
	b.BrickInstance().Pr = req
	return nil
}

func (b *BmRsaKeyGenerateBrick) Done(pkg string, idx int64, e error) error {
	tmp, _ := bmpkg.GetPkgLen(pkg)
	if int(idx) < tmp-1 {
		bmrouter.NextBrickRemote(pkg, idx+1, b)
	}
	return nil
}

func (b *BmRsaKeyGenerateBrick) BrickInstance() *bmbrick.BMBrick {
	if b.bk == nil {
		b.bk = &bmbrick.BMBrick{}
	}
	return b.bk
}

func (b *BmRsaKeyGenerateBrick) ResultTo(w io.Writer) error {
	pr := b.BrickInstance().Pr
	tmp := pr.(auth.BmRsaKey)
	err := jsonapi.ToJsonAPI(&tmp, w)
	return err
}

func (b *BmRsaKeyGenerateBrick) Return(w http.ResponseWriter) {
	ec := b.BrickInstance().Err
	if ec != 0 {
		bmerror.ErrInstance().ErrorReval(ec, w)
	} else {
		var reval auth.BmRsaKey = b.BrickInstance().Pr.(auth.BmRsaKey)
		jsonapi.ToJsonAPI(&reval, w)
	}
}
