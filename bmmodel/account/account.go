package account

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/alfredyang1986/blackmirror/bmconfighandle"
	"github.com/alfredyang1986/blackmirror/bmmodel"
	"github.com/alfredyang1986/blackmirror/bmmodel/request"
	"github.com/alfredyang1986/blackmirror/bmsecurity"
	"github.com/Jeorch/BP-Auth-Server/bmmodel/auth"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"sync"
)

type BmAccount struct {
	Id  string        `json:"id"`
	Id_ bson.ObjectId `bson:"_id"`

	Account    string `json:"account" bson:"account"`
	SecretWord string `json:"secretword" bson:"secretword"`
	BrandId string `json:"brandId" bson:"brandId"`
}

/*------------------------------------------------
 * bm object interface
 *------------------------------------------------*/

func (bd *BmAccount) ResetIdWithId_() {
	bmmodel.ResetIdWithId_(bd)
}

func (bd *BmAccount) ResetId_WithID() {
	bmmodel.ResetId_WithID(bd)
}

/*------------------------------------------------
 * bmobject interface
 *------------------------------------------------*/

func (bd *BmAccount) QueryObjectId() bson.ObjectId {
	return bd.Id_
}

func (bd *BmAccount) QueryId() string {
	return bd.Id
}

func (bd *BmAccount) SetObjectId(id_ bson.ObjectId) {
	bd.Id_ = id_
}

func (bd *BmAccount) SetId(id string) {
	bd.Id = id
}

/*------------------------------------------------
 * relationships interface
 *------------------------------------------------*/
func (bd BmAccount) SetConnect(tag string, v interface{}) interface{} {
	return bd
}

func (bd BmAccount) QueryConnect(tag string) interface{} {
	return bd
}

/*------------------------------------------------
 * mongo interface
 *------------------------------------------------*/

func (bd *BmAccount) InsertBMObject() error {
	return bmmodel.InsertBMObject(bd)
}

func (bd *BmAccount) FindOne(req request.Request) error {
	return bmmodel.FindOne(req, bd)
}

func (bd *BmAccount) UpdateBMObject(req request.Request) error {
	return bmmodel.UpdateOne(req, bd)
}

func (bd *BmAccount) DecodeByCompanyDate(company string, date string) error {

	bmRsaKey := auth.BmRsaKey{
		Company: company,
		Date:    date,
	}

	privateKey, err := bmRsaKey.GetPrivateKey()
	if err != nil {
		return err
	}

	secretWord := bd.SecretWord
	secretByte, err := base64.StdEncoding.DecodeString(secretWord)
	if err != nil {
		return err
	}

	originByte, err := bmsecurity.PhRsaDecrypt(privateKey, secretByte)
	if err != nil {
		return err
	}

	bd.SecretWord = string(originByte)

	return nil
}

func (bd *BmAccount) Secret2MD5() {

	secretWord := bd.SecretWord

	h := md5.New()
	io.WriteString(h, secretWord)

	secretWordMd5 := fmt.Sprintf("%x", h.Sum(nil))
	bd.SecretWord = secretWordMd5

}

var once sync.Once
var bmMongoConfig bmconfig.BMMongoConfig

func (bd BmAccount) IsAccountRegisted() bool {
	once.Do(bmMongoConfig.GenerateConfig)
	session, err := mgo.Dial(bmMongoConfig.Host + ":" + bmMongoConfig.Port)
	if err != nil {
		panic("dial db error")
	}
	defer session.Close()

	c := session.DB(bmMongoConfig.Database).C("BmAccount")
	n, err := c.Find(bson.M{"account": bd.Account}).Count()
	if err != nil {
		panic(err)
	}

	return n > 0
}

func (bd BmAccount) Valid() bool {
	return bd.Account != ""
}
