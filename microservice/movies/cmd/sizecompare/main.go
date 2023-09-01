package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/ahmad-khatib0/go/microservice/movies/gen"
	"github.com/ahmad-khatib0/go/microservice/movies/metadata/pkg/model"
	"github.com/golang/protobuf/proto"
)

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

var genMetadata = &gen.Metadata{
	Id:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

func main() {

	jsonBytes, err := serializeToJson(metadata)
	if err != nil {
		panic(err)
	}
	xmlBytes, err := serializeToXML(metadata)
	if err != nil {
		panic(err)
	}
	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size:\t%dB\n", len(jsonBytes))
	fmt.Printf("XML size:\t%dB\n", len(xmlBytes))
	fmt.Printf("Proto size:\t%dB\n", len(protoBytes))
}

func serializeToJson(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}

func serializeToXML(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func serializeToProto(m *gen.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}
