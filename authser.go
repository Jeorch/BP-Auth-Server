package main

import (
	"github.com/Jeorch/BP-Auth-Server/bmcommon/bmrouter"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/account"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/auth"
	"github.com/Jeorch/BP-Auth-Server/bmpipe/bmaccountbricks/find"
	"github.com/Jeorch/BP-Auth-Server/bmpipe/bmaccountbricks/other"
	"github.com/Jeorch/BP-Auth-Server/bmpipe/bmaccountbricks/push"
	"github.com/Jeorch/BP-Auth-Server/bmpipe/bmauthbricks/find"
	"github.com/Jeorch/BP-Auth-Server/bmpipe/bmauthbricks/push"
	"github.com/alfredyang1986/blackmirror/bmcommon/bmsingleton"
	"github.com/alfredyang1986/blackmirror/bmconfighandle"
	"github.com/alfredyang1986/blackmirror/bmerror"
	"github.com/alfredyang1986/blackmirror/bmmodel/request"
	"net/http"
	"sync"
)

func main() {

	fac := bmsingleton.GetFactoryInstance()

	/*------------------------------------------------
	 * model object
	 *------------------------------------------------*/
	fac.RegisterModel("Request", &request.Request{})
	fac.RegisterModel("Eqcond", &request.Eqcond{})
	fac.RegisterModel("Necond", &request.Necond{})
	fac.RegisterModel("Gtcond", &request.Gtcond{})
	fac.RegisterModel("Gtecond", &request.Gtecond{})
	fac.RegisterModel("Ltcond", &request.Ltcond{})
	fac.RegisterModel("Ltecond", &request.Ltecond{})
	fac.RegisterModel("Incond", &request.Incond{})
	fac.RegisterModel("Nincond", &request.Nincond{})
	fac.RegisterModel("Upcond", &request.Upcond{})
	fac.RegisterModel("Fmcond", &request.Fmcond{})
	fac.RegisterModel("BmErrorNode", &bmerror.BmErrorNode{})

	fac.RegisterModel("BmRsaKey", &auth.BmRsaKey{})
	fac.RegisterModel("BmAccount", &account.BmAccount{})
	fac.RegisterModel("BmLoginSucceed", &auth.BmLoginSucceed{})

	/*------------------------------------------------
	 * rsa bricks
	 *------------------------------------------------*/
	fac.RegisterModel("BmGetPublicKeyBrick", &rsafind.BmGetPublicKeyBrick{})
	fac.RegisterModel("BmRsaKeyGenerateBrick", &rsapush.BmRsaKeyGenerateBrick{})

	/*------------------------------------------------
	 * account bricks
	 *------------------------------------------------*/
	fac.RegisterModel("BmAccountPushBrick", &accountpush.BmAccountPushBrick{})
	fac.RegisterModel("BmAccountFindBrick", &accountfind.BmAccountFindBrick{})
	fac.RegisterModel("BmAccountGenerateToken", &accountother.BmAccountGenerateToken{})

	r := bmrouter.BindRouter()

	var once sync.Once
	var bmRouter bmconfig.BMRouterConfig
	once.Do(bmRouter.GenerateConfig)

	http.ListenAndServe(":"+bmRouter.Port, r)
}
