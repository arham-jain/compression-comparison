package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"github.com/arham-jain/compression-comparison/proto"
	protobuf "github.com/golang/protobuf/proto"
	"log"
	"time"
)

type encoding func(entity Entity) []byte

type decoding func([]byte) Entity

type Entity struct {
	Name    string   `json:"name"`
	Details string   `json:"details"`
	Tags    []string `json:"tags"`
}

func main() {
	entity := Entity{
		Name:    "arham",
		Details: "jain",
		Tags:    []string{"97", "male"},
	}
	timeExecution(jsonEncoding, jsonDecoding, entity, "json")
	timeExecution(gobEncoding, gobDecoding, entity, "gob")
	timeExecution(gobGzipEncoding, gobGzipDecoding, entity, "gob+gzip")
	timeExecution(protoEncoding, protoDecoding, entity, "proto")
}

func timeExecution(encoder encoding, decoder decoding, entity Entity, typ string) {
	startTime1 := time.Now()
	out := encoder(entity)
	endTime1 := time.Now()
	startTime2 := time.Now()
	decoder(out)
	endTime2 := time.Now()
	log.Printf("Type: %s\t\tEncoding time: %v\t Decoding time: %v\t Encoded length: %v", typ,
		endTime1.Sub(startTime1), endTime2.Sub(startTime2), len(out))
}

func verify(entity Entity) {
	jsonBytes := jsonEncoding(entity)
	gobBytes := gobEncoding(entity)
	gzipGobBytes := gobGzipEncoding(entity)
	protoBytes := protoEncoding(entity)
	jsonDecoding(jsonBytes)
	gobDecoding(gobBytes)
	gobGzipDecoding(gzipGobBytes)
	protoDecoding(protoBytes)
}

func jsonEncoding(entity Entity) []byte {
	byteList, err := json.Marshal(entity)
	if err == nil {
		return byteList
	}
	return nil
}

func jsonDecoding(byteList []byte) (entity Entity) {
	_ = json.Unmarshal(byteList, &entity)
	//log.Printf("Json Decoding: %v\n", entity)
	return
}

func gobEncoding(entity Entity) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(entity)
	if err == nil {
		return buf.Bytes()
	}
	return nil
}

func gobDecoding(byteList []byte) (entity Entity) {
	var buf bytes.Buffer
	buf.Write(byteList)
	dec := gob.NewDecoder(&buf)
	_ = dec.Decode(&entity)
	return
}

func gobGzipEncoding(entity Entity) []byte {
	var buf bytes.Buffer
	out := gzip.NewWriter(&buf)
	enc := gob.NewEncoder(out)
	err := enc.Encode(entity)
	if err == nil {
		return buf.Bytes()
	}
	return nil
}

func gobGzipDecoding(byteList []byte) (entity Entity) {
	var buf bytes.Buffer
	buf.Write(byteList)
	decompressed, _ := gzip.NewReader(&buf)
	dec := gob.NewDecoder(decompressed)
	_ = dec.Decode(&entity)
	return
}

func protoEncoding(entity Entity) []byte {
	protoEntity := proto.Entity{
		Name:    entity.Name,
		Details: entity.Details,
		Tags:    entity.Tags,
	}
	out, err := protobuf.Marshal(&protoEntity)
	if err == nil {
		return out
	}
	return nil
}

func protoDecoding(byteList []byte) (entity Entity) {
	var protoEntity proto.Entity
	_ = protobuf.Unmarshal(byteList, &protoEntity)
	entity = Entity{
		Name:    protoEntity.Name,
		Details: protoEntity.Details,
		Tags:    protoEntity.Tags,
	}
	return
}
