package RunningCommand

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestStartAndDoWithFeedback(t *testing.T) {
	Convey("Testing StartAndDoWithFeedback", t, func() {
		Convey("Start application to run echo with no errors", func() {
			for i := 0; i < 20; i++ {
				rc := New(nil, nil, "cmd", "/c", "echo", "hello")

				outLines := []string{}
				errLines := []string{}
				rc.StartAndDoWithFeedback(0, func(isErr bool, fb string) {
					if isErr {
						errLines = append(errLines, fb)
					} else {
						outLines = append(outLines, fb)
					}
				})

				So(len(outLines), ShouldEqual, 2)
				So(strings.TrimSpace(outLines[0]), ShouldEqual, "Command has started.")
				So(strings.TrimSpace(outLines[1]), ShouldEqual, "hello")

				So(len(errLines), ShouldEqual, 0)
			}
		})

		Convey("Start application WITH errors", func() {
			for i := 0; i < 20; i++ {
				rc := New(nil, nil, "cmd", "/c", "echo", "hello", "&", "exit", "/b", "13")

				outLines := []string{}
				errLines := []string{}
				rc.StartAndDoWithFeedback(0, func(isErr bool, fb string) {
					if isErr {
						errLines = append(errLines, fb)
					} else {
						outLines = append(outLines, fb)
					}
				})

				So(len(outLines), ShouldEqual, 2)
				So(strings.TrimSpace(outLines[0]), ShouldEqual, "Command has started.")
				So(strings.TrimSpace(outLines[1]), ShouldEqual, "hello")

				So(len(errLines), ShouldEqual, 1)
				So(strings.TrimSpace(errLines[0]), ShouldEqual, "Unable to finish waiting for command: exit status 13")
			}
		})
	})
}
