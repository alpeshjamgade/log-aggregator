package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
)

func compress(jsonBytes []byte) []byte {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, _ = gzipWriter.Write(jsonBytes)
	_ = gzipWriter.Close()
	compressed := buf.Bytes()

	return compressed
}

func decompress(compressed []byte) []byte {
	encodedStr := base64.StdEncoding.EncodeToString(compressed)

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		panic(err)
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(decodedBytes))
	if err != nil {
		panic(err)
	}
	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		panic(err)
	}
	_ = gzipReader.Close()

	return decompressed
}

func main() {
	// Your log data
	//data := []map[string]interface{}{
	//	{
	//		"date": 1752659073.838018,
	//		"log":  `2025-07-16T09:44:33.837928722Z stdout F {"request_id":"abcd","level":"\u001b[33mWARN\u001b[0m","ts":"2025-07-16T15:14:33.837+0530","caller":"cashPositions/controller_v2.go:495","msg":"no balance for cash positions -> <nil>"}`,
	//		"kubernetes": map[string]interface{}{
	//			"pod_name":        "oms-7568cdc985-mpwk2",
	//			"namespace_name":  "default",
	//			"pod_id":          "1359d67e-b915-43dc-9eb0-59f07538c64c",
	//			"host":            "gke-tradelab-uat-gke-a-main-node-pool-d458a698-sahj",
	//			"container_name":  "oms",
	//			"docker_id":       "99e43629af1d56781d00aa92b490100225a26f94ee154f74d5ea9a12773b6c06",
	//			"container_hash":  "registry.tradelab.in/orms@sha256:55144c5ba7d9e3ee01347c31a6c6f87c7d29be9b5c70b518bd27fac827d2459a",
	//			"container_image": "registry.tradelab.in/orms:tradelab_uat-a60a7bc5",
	//			"labels": map[string]string{
	//				"app":               "oms",
	//				"pod-template-hash": "7568cdc985",
	//			},
	//			"annotations": map[string]string{
	//				"cni.projectcalico.org/containerID": "d8d24cf258d4d2b40b237594642a2e47efbfa9f2034ca3cd56988b759a25bab6",
	//				"cni.projectcalico.org/podIP":       "10.96.1.206/32",
	//				"cni.projectcalico.org/podIPs":      "10.96.1.206/32",
	//				"kubectl.kubernetes.io/restartedAt": "2025-07-15T13:41:41+05:30",
	//			},
	//		},
	//	},
	//}
	//
	//jsonBytes, err := json.Marshal(data)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Original JSON size: %v Bytes\n", len(jsonBytes))
	//
	//compressed := compress(jsonBytes)
	//fmt.Printf("Compressed size: %v Bytes\n", len(compressed))
	//
	//fmt.Println(string(compressed))
	//
	//decompressed := decompress(compressed)
	//fmt.Println("\nDecompressed JSON string:\n", string(decompressed))
}
