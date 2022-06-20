package test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/pelletier/go-toml"
)

func TestName(t *testing.T) {
	log.Println("run test")

	toml_conf_path := "/Users/bruce/workspace/go/project/peer-node/configs/default.toml"

	log.Println(toml_conf_path)

	doc, err := ioutil.ReadFile(toml_conf_path)
	if err != nil {
		log.Println(err)
	}

	tree, err := toml.Load(string(doc))
	log.Println(tree)

	pp := tree.Get("storage222.password")
	log.Println(pp)

	p := tree.Get("storage.password")
	s := tree.Get("cache.size")

	log.Println(p)
	log.Println(s)

	//v := tree.Values()
	//log.Println(v)
	//
	//m := tree.ToMap()
	//log.Println(m)
	//
	//m["https_port"] = 555
	//
	//b, _ := toml.Marshal(m)
	//
	//f, err := os.OpenFile("/Users/bruce/workspace/go/project/peer-node/configs/default.test.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	//if err != nil {
	//	//return err
	//}
	//
	//defer f.Close()
	//
	//err = f.Truncate(0)
	//if err != nil {
	//	//return err
	//}
	//_, err = f.Seek(0, 0)
	//if err != nil {
	//	//return err
	//}
	//
	//_, err = f.Write(b)
	//if err != nil {
	//	//return err
	//}
}
