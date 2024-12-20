package main

import (
	aw "github.com/deanishe/awgo"
	"tools/util/decode"
	"tools/util/ip"

	"strconv"
	"time"

	"tools/util/consts"
	"tools/util/datetimes"
	"tools/util/encode"
	"tools/util/hash"
)

func run() {
	var err error
	args := consts.Workflow.Args()
	if len(args) == 0 {
		return
	}

	defer func() {
		if err == nil {
			consts.Workflow.SendFeedback()
			return
		}
	}()

	if args[0] == "time" {
		if len(args) < 2 {
			return
		}
		if args[1] == "now" {
			datetimes.ProcessNow()
			return
		}
		if datetimes.RegexpTimestamp.MatchString(args[1]) {
			v, err := strconv.ParseInt(args[1], 10, 64) // 转为 int64 支持更大范围
			if err == nil {
				switch length := len(args[1]); {
				case length <= 10: // 秒级时间戳
					datetimes.ProcessTimestamp(time.Unix(v, 0)) // 秒部分直接使用
				case length <= 13: // 毫秒级时间戳
					seconds := v / 1000                   // 提取秒部分
					nanoseconds := (v % 1000) * 1_000_000 // 毫秒部分转为纳秒
					datetimes.ProcessTimestamp(time.Unix(seconds, nanoseconds))
				case length <= 16: // 微秒级时间戳
					seconds := v / 1_000_000               // 提取秒部分
					nanoseconds := (v % 1_000_000) * 1_000 // 微秒部分转为纳秒
					datetimes.ProcessTimestamp(time.Unix(seconds, nanoseconds))
				case length > 16: // 纳秒级时间戳
					seconds := v / 1_000_000_000     // 提取秒部分
					nanoseconds := v % 1_000_000_000 // 纳秒部分直接使用
					datetimes.ProcessTimestamp(time.Unix(seconds, nanoseconds))
				}
				return
			}
			return
		}
		// 处理时间字符串
		err = datetimes.ProcessTimeStr(args[1])
	} else if args[0] == "hash" {
		if len(args) < 2 {
			return
		}
		hash.HashMain(args[1])
		return
	} else if args[0] == "encode" {
		if len(args) < 2 {
			return
		}
		encode.Main(args[1])
		return
	} else if args[0] == "decode" {
		if len(args) < 2 {
			return
		}
		decode.Main(args[1])
		return
	} else if args[0] == "ip" {
		if len(args) < 2 {
			return
		}
		ip.Main(args[1])
		return
	}

}

func main() {
	consts.Workflow = aw.New()
	consts.Workflow.Run(run)
}
