/*
MIT License

Copyright (c) 2017 Master.G

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net"
	"os"
	"runtime/debug"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/maruel/panicparse/stack"
)

// PrintPanicStack prints panic stack info
func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		in := bytes.NewBufferString(string(debug.Stack()[:]))
		goroutines, err := stack.ParseDump(in, os.Stdout)
		if err != nil {
			return
		}

		// Optional: Check for GOTRACEBACK being set, in particular if there is only
		// one goroutine returned.

		// Use a color palette based on ANSI code.
		p := &stack.Palette{}
		buckets := stack.SortBuckets(stack.Bucketize(goroutines, stack.AnyValue))
		srcLen, pkgLen := stack.CalcLengths(buckets, false)
		for _, bucket := range buckets {
			log.Error(p.BucketHeader(&bucket, false, len(buckets) > 1))
			log.Error("\n" + p.StackLines(&bucket.Signature, srcLen, pkgLen, false))
		}

		for k := range extras {
			log.Errorf("EXTRAS#%v DATA:%v\n", k, spew.Sdump(extras[k]))
		}
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// GetStringMD5Hash will return text's md5 digest hex encoded
func GetStringMD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// Timestamp returns current unix timestamp in uint64 format
func Timestamp() uint64 {
	return uint64(time.Now().Unix())
}
