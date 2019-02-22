package auth

import (
	"fmt"
	"github.com/alfredyang1986/blackmirror/bmmodel"
	"github.com/alfredyang1986/blackmirror/bmmodel/request"
	"gopkg.in/mgo.v2/bson"
)

type BmRsaKey struct {
	Id  string        `json:"id"`
	Id_ bson.ObjectId `bson:"_id"`

	Company    string `json:"company" bson:"company"`
	Date       string `json:"date" bson:"date"`
	PublicKey  string `json:"publicKey" bson:"publicKey"`
	PrivateKey string `json:"privateKey" bson:"privateKey"`
}

/*------------------------------------------------
 * bm object interface
 *------------------------------------------------*/

func (bd *BmRsaKey) ResetIdWithId_() {
	bmmodel.ResetIdWithId_(bd)
}

func (bd *BmRsaKey) ResetId_WithID() {
	bmmodel.ResetId_WithID(bd)
}

/*------------------------------------------------
 * bmobject interface
 *------------------------------------------------*/

func (bd *BmRsaKey) QueryObjectId() bson.ObjectId {
	return bd.Id_
}

func (bd *BmRsaKey) QueryId() string {
	return bd.Id
}

func (bd *BmRsaKey) SetObjectId(id_ bson.ObjectId) {
	bd.Id_ = id_
}

func (bd *BmRsaKey) SetId(id string) {
	bd.Id = id
}

/*------------------------------------------------
 * relationships interface
 *------------------------------------------------*/
func (bd BmRsaKey) SetConnect(tag string, v interface{}) interface{} {
	return bd
}

func (bd BmRsaKey) QueryConnect(tag string) interface{} {
	return bd
}

/*------------------------------------------------
 * mongo interface
 *------------------------------------------------*/

func (bd *BmRsaKey) InsertBMObject() error {
	return bmmodel.InsertBMObject(bd)
}

func (bd *BmRsaKey) FindOne(req request.Request) error {
	return bmmodel.FindOne(req, bd)
}

func (bd *BmRsaKey) UpdateBMObject(req request.Request) error {
	return bmmodel.UpdateOne(req, bd)
}

func (bd *BmRsaKey) GetPrivateKey() (string, error) {

	eq1 := request.Eqcond{}
	eq1.Ky = "company"
	eq1.Vy = bd.Company
	eq2 := request.Eqcond{}
	eq2.Ky = "date"
	eq2.Vy = bd.Date
	req := request.Request{}
	req.Res = "BmRsaKey"
	var condi []interface{}
	condi = append(condi, eq1, eq2)
	c := req.SetConnect("conditions", condi)
	fmt.Println(c)

	err := bd.FindOne(c.(request.Request))
	return bd.PrivateKey, err
}
