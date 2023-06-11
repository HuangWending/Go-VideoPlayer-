# Go-VideoPlayer-
Go视频播放器
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nsf/termbox-go"
)

导入所需的包和库。


const (
	playbackSpeed = 1.0 // 默认播放速度
)

定义常量，设置默认播放速度。

var (
	videoFile      string // 视频文件路径
	isPlaying      bool   // 是否正在播放
	isPaused       bool   // 是否已暂停
	currentSpeed   float64
	volumeLevel    int // 音量级别
	audioSupported bool
)

声明变量，用于存储视频文件路径、播放状态、播放速度、音量级别以及音频支持情况。


func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// 初始化播放器状态
	isPlaying = false
	isPaused = false
	currentSpeed = playbackSpeed
	volumeLevel = 50

	// 获取视频文件路径
	videoFile = getVideoFilePath()

	// 检查是否支持音频
	audioSupported = checkAudioSupport()

	// 显示播放器界面
	render()

	// 监听键盘事件
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return
			}

			switch ev.Ch {
			case 'p', 'P': // 播放/暂停
				togglePlayPause()
			case 's', 'S': // 减速
				decreaseSpeed()
			case 'f', 'F': // 加速
				increaseSpeed()
			case 'u', 'U': // 音量增加
				increaseVolume()
			case 'd', 'D': // 音量减少
				decreaseVolume()
			}

			render()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

主函数，程序入口。首先初始化终端界面，然后初始化播放器状态和相关变量。接下来获取用户输入的视频文件路径，并检查是否支持音频。之后显示播放器界面，并开始监听键盘事件。根据用户的输入执行相应的操作，并不断重新渲染界面。

func getVideoFilePath() string {
	fmt.Print("请输入视频文件的路径：")
	var filePath string
	fmt.Scanln(&filePath)

	return filePath
}

获取用户输入的视频文件路径。

func checkAudioSupport() bool {
	cmd := exec.Command("ffprobe", "-show_streams", "-i", videoFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("检查音频支持时出错:", err)
		return false
	}

	streams := strings.Split(string(output), "[STREAM]")
	for _, stream := range streams {
		if strings.Contains(stream, "codec_type=audio") {
			return true
		}
	}

	return false
}

检查视频文件是否支持音频。使用外部工具ffprobe执行命令来获取视频文件的信息，然后检查输出中是否包含音频流的信息。


func togglePlayPause() {
	if !

isPlaying && !isPaused {
		playVideo()
	} else if isPlaying && !isPaused {
		pauseVideo()
	} else if isPlaying && isPaused {
		resumeVideo()
	}
}

切换播放/暂停状态。根据当前播放器的状态，调用相应的函数来执行播放、暂停或恢复播放的操作。

func playVideo() {
	cmd := exec.Command("ffplay", "-i", videoFile, "-vf", fmt.Sprintf("setpts=%.2f*PTS", currentSpeed))
	err := cmd.Start()
	if err != nil {
		fmt.Println("播放视频时出错:", err)
		return
	}

	isPlaying = true
	isPaused = false
}
播放视频。使用外部工具ffplay执行命令来播放视频，并根据当前的播放速度调整视频播放速度。

func pauseVideo() {
	cmd := exec.Command("pkill", "-STOP", "-f", "ffplay")
	err := cmd.Run()
	if err != nil {
		fmt.Println("暂停视频时出错:", err)
		return
	}

	isPaused = true
}

暂停视频播放。使用外部工具pkill执行命令来发送停止信号给ffplay进程，从而暂停视频的播放。

func resumeVideo() {
	cmd := exec.Command("pkill", "-CONT", "-f", "ffplay")
	err := cmd.Run()
	if err != nil {
		fmt.Println("恢复视频播放时出错:", err)
		return
	}

	isPaused = false
}

恢复视频播放。使用外部工具pkill执行命令来发送继续信号给ffplay进程，从而恢复视频的播放。

func increaseSpeed() {
	currentSpeed += 0.1
	if currentSpeed > 2.0 {
		currentSpeed = 2.0
	}
}

增加播放速度。将当前播放速度增加0.1，并限制最大速度为2.0。


func decreaseSpeed() {
	currentSpeed -= 0.1
	if currentSpeed < 0.1 {
		currentSpeed = 0.1
	}
}

减少播放速度。将当前播放速度减少0.1，并限制最小速度为0.1。


func increaseVolume() {
	volumeLevel += 10
	if volumeLevel > 100 {
		volumeLevel = 100
	}
}

增加音量。将当前音量级别增加10，并限制最大音量级别为100。


func decreaseVolume() {
	volumeLevel -= 10
	if volumeLevel < 0 {
		volumeLevel = 0
	}
}
```
减少音量。将当前音量级别减少10，并限制最小音量级别为0。

func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// 绘制视频播放器框架
	drawPlayerFrame()

	// 绘制功能按钮
	drawControls()

	termbox.Flush()
}

绘制界面。清空终端窗口并绘制视频播放器的框架和

功能按钮。

func drawPlayerFrame() {
	width, height := termbox.Size()

	for y := 0; y < height; y++ {
		termbox.SetCell(0, y, '│', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(width-1, y, '│', termbox.ColorDefault, termbox.ColorDefault)
	}

	for x := 1; x < width-1; x++ {
		termbox.SetCell(x, 0, '─', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x, height-1, '─', termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.SetCell(0, 0, '┌', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(width-1, 0, '┐', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(0, height-1, '└', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(width-1, height-1, '┘', termbox.ColorDefault, termbox.ColorDefault)
}

绘制视频播放器框架。使用termbox库的SetCell函数来设置字符的显示位置和样式，绘制出视频播放器的边框。

func drawControls() {
	width, _ := termbox.Size()
	y := 2

	drawButton(2, y, "[P] 播放/暂停", !isPlaying || isPaused)
	y += 2

	drawButton(2, y, "[S] 减速", isPlaying && !isPaused && currentSpeed <= 0.1)
	drawButton(16, y, "[F] 加速", isPlaying && !isPaused && currentSpeed >= 2.0)
	y += 2

	drawButton(2, y, "[U] 音量增加", !isPlaying || isPaused || volumeLevel >= 100)
	drawButton(18, y, "[D] 音量减少", !isPlaying || isPaused || volumeLevel <= 0)
	y += 2

	if !audioSupported {
		drawMessage(2, y, "警告: 视频不包含音频轨道")
	}
}

绘制功能按钮。根据播放器的状态和音频支持情况，调用drawButton函数和drawMessage函数来绘制不同样式的功能按钮。

func drawButton(x, y int, label string, disabled bool) {
	fgColor := termbox.ColorDefault
	bgColor := termbox.ColorGreen

	if disabled {
		fgColor = termbox.ColorBlack
		bgColor = termbox.ColorRed
	}

	for _, ch := range label {
		termbox.SetCell(x, y, ch, fgColor, bgColor)
		x++
	}
}

绘制按钮。根据按钮的位置、标签、禁用状态，使用termbox库的SetCell函数来设置字符的显示位置和样式，绘制出按钮。


func drawMessage(x, y int, message string) {
	fgColor := termbox.ColorRed
	bgColor := termbox.ColorDefault

	for _, ch := range message {
		termbox.SetCell(x, y, ch, fgColor, bgColor)
		x++
	}
}

绘制消息。根据消息的位置

、内容，使用termbox库的SetCell函数来设置字符的显示位置和样式，绘制出消息。
使用了termbox-go库来实现终端界面，结合外部工具ffplay和ffprobe来处理视频播放和音频相关的操作。
